package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/FeshLig/metrcollector/internal/config"
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/persister"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/router"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	err := Run()
	if err != nil {
		log.Fatal(err)
	}

}

func Run() error {

	var storage repository.Storage
	var db *pgxpool.Pool
	var err error

	cfg := config.GetOptions()

	ctx, cancel := newStartupContext()
	defer cancel()

	if cfg.DatabaseDSN.String() != "" {

		db, err = newDB(ctx, cfg)
		if err != nil {
			return err
		}
		if db != nil {
			defer db.Close()
		}

		storage = repository.NewPostgresStorage(db)

	} else {

		storage = repository.NewMemStorage()
	}

	persister, err := newPersister(cfg, storage)
	if err != nil {
		return err
	}
	defer persister.Stop()

	service := NewService(cfg, storage, persister)
	handlers := handler.NewHandlers(service, db)
	router := router.NewRouter(handlers)

	if err := router.Run(cfg.Address.String()); err != nil {
		return err
	}

	return nil

}

func newStartupContext() (context.Context, context.CancelFunc) {

	const t = 5 * time.Second

	return context.WithTimeout(context.Background(), t)

}

func newDB(ctx context.Context, cfg config.Options) (*pgxpool.Pool, error) {

	dsn := cfg.DatabaseDSN.String()

	if err := repository.RunMigrations(dsn); err != nil {
		return nil, err
	}

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("db init failed: %w", err)
	}

	return db, nil

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
