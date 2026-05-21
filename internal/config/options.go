package config

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/FeshLig/metrcollector/internal/flags"
)

// Options contains application configuration values.
type Options struct {
	Address         flags.NetAddress
	StoreInterval   flags.SecondsStoreInterval
	FileStoragePath flags.FileStoragePath
	Restore         flags.Restore
	DatabaseDSN     flags.DatabaseDSN
	Key             flags.Key
	AuditFile       flags.AuditFile
	AuditURL        flags.AuditURL
}

// GetOptions parses application configuration from flags and environment variables.
func GetOptions() Options {

	const fileName = "metrics.txt"

	options := Options{
		Address: flags.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
		StoreInterval: flags.SecondsStoreInterval{
			Duration: time.Duration(300) * time.Second,
		},
		FileStoragePath: "/tmp/" + fileName,
		Restore:         false,
		DatabaseDSN:     "",
		Key:             "",
		AuditFile:       "",
		AuditURL:        "",
	}

	parseFlags(&options)
	parseEnv(&options)

	pathStr := options.FileStoragePath.String()

	info, err := os.Stat(pathStr)
	if err == nil && info.IsDir() {
		pathStr = filepath.Join(pathStr, fileName)
	} else if errors.Is(err, fs.ErrNotExist) {
		if filepath.Ext(pathStr) == "" {
			pathStr = filepath.Join(pathStr, fileName)
		}
	}
	_ = options.FileStoragePath.Set(pathStr)

	return options

}

// parseFlags parses command-line flags.
func parseFlags(options *Options) {

	flag.Var(&options.Address, "a", "net address host:port")
	flag.Var(&options.StoreInterval, "i", "store interval (seconds)")
	flag.Var(&options.FileStoragePath, "f", "file storage path")
	flag.Var(&options.Restore, "r", "restore file (true/false)")
	flag.Var(&options.DatabaseDSN, "d", "postgres dsn (format: postgres://user:password@host:port/dbname)")
	flag.Var(&options.Key, "k", "hash key")
	flag.Var(&options.AuditFile, "audit-file", "audit file path")
	flag.Var(&options.AuditURL, "audit-url", "audit url")

	flag.Parse()

}

// parseEnv parses environment variables.
func parseEnv(options *Options) error {

	if addrStr, ok := os.LookupEnv("ADDRESS"); ok {
		err := options.Address.Set(addrStr)
		if err != nil {
			return fmt.Errorf("wrong value of ADDRESS: %w", err)
		}
	}

	if intvlStr, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		err := options.StoreInterval.Set(intvlStr)
		if err != nil {
			return fmt.Errorf("wrong value of STORE_INTERVAL: %w", err)
		}
	}

	if pathStr, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		err := options.FileStoragePath.Set(pathStr)
		if err != nil {
			return fmt.Errorf("wrong value of FILE_STORAGE_PATH: %w", err)
		}
	}

	if rstrStr, ok := os.LookupEnv("RESTORE"); ok {
		err := options.Restore.Set(rstrStr)
		if err != nil {
			return fmt.Errorf("wrong value of RESTORE: %w", err)
		}
	}

	if databaseDSNStr, ok := os.LookupEnv("DATABASE_DSN"); ok {
		err := options.DatabaseDSN.Set(databaseDSNStr)
		if err != nil {
			return fmt.Errorf("wrong value of DATABASE_DSN: %w", err)
		}
	}

	if keyStr, ok := os.LookupEnv("KEY"); ok {
		err := options.Key.Set(keyStr)
		if err != nil {
			return fmt.Errorf("wrong value of KEY: %w", err)
		}
	}

	if auditFileStr, ok := os.LookupEnv("AUDIT_FILE"); ok {
		err := options.AuditFile.Set(auditFileStr)
		if err != nil {
			return fmt.Errorf("wrong value of AUDIT_FILE: %w", err)
		}
	}

	if auditURLStr, ok := os.LookupEnv("AUDIT_URL"); ok {
		err := options.AuditURL.Set(auditURLStr)
		if err != nil {
			return fmt.Errorf("wrong value of AUDIT_URL: %w", err)
		}
	}

	return nil

}
