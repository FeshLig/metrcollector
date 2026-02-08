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

	cfg := config.GetOptions()

	ctx, cancel := newStartapContext()
	defer cancel()

	db, err := newDB(ctx, cfg)
	if err != nil {
		return err
	}

	memStorage := repository.NewMemStorage()

	persister, err := newPersister(cfg, memStorage)
	if err != nil {
		return err
	}
	defer persister.Stop()

	service := NewService(cfg, memStorage, persister)
	handlers := handler.NewHandlers(service, db)
	router := router.NewRouter(handlers)

	router.Run(cfg.Address.String())

	return nil

}

func newStartapContext() (context.Context, context.CancelFunc) {

	const t = 5 * time.Second

	return context.WithTimeout(context.Background(), t)

}

func newDB(ctx context.Context, cfg config.Options) (*pgxpool.Pool, error) {

	// dsn := "postgres://metrics:Fjytotbytn4rjtvju@localhost:5432/metrics"
	dsn := cfg.DatabaseDSN.String()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("db init failed: %w", err)
	}

	return db, nil

}

func newPersister(cfg config.Options, storage persister.MetricsStorage) (*persister.FilePersister, error) {

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
