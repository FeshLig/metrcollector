package main

import (
	"context"
	"log"
	"time"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/FeshLig/metrcollector/internal/config"
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/persister"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/router"
	"github.com/FeshLig/metrcollector/internal/service"
)

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

}

func run() error {

	var storage repository.Storage
	var persister *persister.FilePersister

	cfg := config.GetOptions()

	ctx, cancel := newStartupContext()
	defer cancel()

	if cfg.DatabaseDSN.String() != "" {

		s, db, err := newDB(ctx, cfg)
		if err != nil {
			return err
		}
		if db != nil {
			defer db.Close()
		}

		persister = nil
		storage = s

	} else {

		storage = repository.NewMemStorage()
		persister, err := newPersister(cfg, storage)
		if err != nil {
			return err
		}
		defer persister.Stop()
	}

	auditPublisher, err := newAudit(cfg)
	if err != nil {
		return err
	}

	service := NewService(cfg, storage, persister)
	handlers := handler.NewHandlers(service, auditPublisher)
	router := router.NewRouter(handlers, cfg)

	if err := router.Run(cfg.Address.String()); err != nil {
		return err
	}

	return nil

}

func newAudit(cfg config.Options) (*audit.Publisher, error) {
	publisher := audit.NewPublisher()
	if cfg.AuditFile.String() != "" {
		fileObserver, err := audit.NewFileObserver(string(cfg.AuditFile))
		if err != nil {
			return nil, err
		}

		// defer fileObserver.Close()

		publisher.Subscribe(fileObserver)
	}

	if cfg.AuditURL.String() != "" {
		httpObserver := audit.NewHTTPObserver(string(cfg.AuditURL))

		publisher.Subscribe(httpObserver)
	}

	return publisher, nil
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

func NewService(cfg config.Options, storage repository.Storage, persister *persister.FilePersister) service.MetricsService {

	service := service.NewMetricService(
		storage,
		persister,
		cfg.StoreInterval.Duration == 0,
	)

	return service

}
