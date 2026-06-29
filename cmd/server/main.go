package main

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/FeshLig/metrcollector/internal/buildinfo"
	"github.com/FeshLig/metrcollector/internal/config"
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/persister"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/router"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/FeshLig/metrcollector/pkg/crypto"
	"go.uber.org/zap"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {

	buildinfo.Print(buildVersion, buildDate, buildCommit)
	err := run()
	if err != nil {
		log.Fatal(err)
	}

}

func run() error {

	var storage repository.Storage
	var persister *persister.FilePersister

	cfg := config.GetOptions()

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()

	ctx, cancel := newStartupContext()
	defer cancel()

	if cfg.DatabaseDSN.String() != "" {

		var postgresStorage *repository.PostgresStorage
		var db *repository.Postgres
		postgresStorage, db, err = newDB(ctx, cfg)
		if err != nil {
			return err
		}
		if db != nil {
			defer db.Close()
		}

		persister = nil
		storage = postgresStorage

	} else {

		storage = repository.NewMemStorage()
		persister, err = newPersister(cfg, storage)
		if err != nil {
			return err
		}
		defer persister.Stop()
	}

	auditPublisher, closeAudit, err := newAudit(cfg)
	if err != nil {
		return err
	}
	defer closeAudit()

	var privateKey *rsa.PrivateKey
	if keyPath := cfg.CryptoKey.String(); keyPath != "" {
		privateKey, err = crypto.LoadPrivateKey(keyPath)
		if err != nil {
			return fmt.Errorf("load private key: %w", err)
		}
	}

	if persister != nil {
		defer func() {
			if err := persister.Save(); err != nil {
				logger.Error("final save failed", zap.Error(err))
			} else {
				logger.Info("metrics flushed to disk")
			}
		}()
	}

	service := newService(cfg, storage, persister)
	handlers := handler.NewHandlers(service, auditPublisher)
	ginRouter := router.NewRouter(handlers, cfg, logger, privateKey)

	serverErr := make(chan error, 1)
	srv := &http.Server{
		Addr:    cfg.Address.String(),
		Handler: ginRouter,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case err := <-serverErr:
		return err
	case sig := <-sigChan:
		logger.Info("received signal, shutting down", zap.String("signal", sig.String()))
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return nil

}

func newAudit(cfg config.Options) (*audit.Publisher, func() error, error) {
	publisher := audit.NewPublisher()

	var closers []func() error

	if cfg.AuditFile.String() != "" {
		fileObserver, err := audit.NewFileObserver(string(cfg.AuditFile))
		if err != nil {
			return nil, nil, err
		}

		closers = append(closers, fileObserver.Close)

		publisher.Subscribe(fileObserver)
	}

	if cfg.AuditURL.String() != "" {
		httpObserver := audit.NewHTTPObserver(string(cfg.AuditURL))

		publisher.Subscribe(httpObserver)
	}

	closeFn := func() error {
		for _, closer := range closers {
			if err := closer(); err != nil {
				return err
			}
		}

		return nil
	}

	return publisher, closeFn, nil
}

func newStartupContext() (context.Context, context.CancelFunc) {

	const t = 5 * time.Second

	return context.WithTimeout(context.Background(), t)

}

func newDB(ctx context.Context, cfg config.Options) (*repository.PostgresStorage, *repository.Postgres, error) {

	dsn := cfg.DatabaseDSN.String()

	db, err := repository.NewPostgres(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}

	storage := repository.NewPostgresStorage(db)

	return storage, db, nil

}

func newPersister(cfg config.Options, storage repository.Storage) (*persister.FilePersister, error) {

	filePath := cfg.FileStoragePath.String()
	storeInterval := cfg.StoreInterval.Duration

	persister := persister.NewFilePersister(
		filePath,
		storage,
		storeInterval,
	)

	if cfg.Restore {
		err := persister.Load()
		if err != nil {
			return nil, err
		}
	}

	persister.Start()

	return persister, nil

}

func newService(cfg config.Options, storage repository.Storage, persister *persister.FilePersister) service.MetricsService {

	service := service.NewMetricService(
		storage,
		persister,
		cfg.StoreInterval.Duration == 0,
	)

	return service

}
