package config

import (
	"encoding/json"
	"os"
)

// ServerConfigFile represents the JSON configuration file structure for the server.
type ServerConfigFile struct {
	Address       string `json:"address"`
	Restore       *bool  `json:"restore"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DatabaseDSN   string `json:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
}

func parseConfigFile(path string) (*ServerConfigFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg ServerConfigFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func applyConfigFile(options *Options, cfg *ServerConfigFile) error {
	if cfg.Address != "" {
		if err := options.Address.Set(cfg.Address); err != nil {
			return err
		}
	}
	if cfg.Restore != nil {
		val := "false"
		if *cfg.Restore {
			val = "true"
		}
		if err := options.Restore.Set(val); err != nil {
			return err
		}
	}
	if cfg.StoreInterval != "" {
		if err := options.StoreInterval.Set(cfg.StoreInterval); err != nil {
			return err
		}
	}
	if cfg.StoreFile != "" {
		if err := options.FileStoragePath.Set(cfg.StoreFile); err != nil {
			return err
		}
	}
	if cfg.DatabaseDSN != "" {
		if err := options.DatabaseDSN.Set(cfg.DatabaseDSN); err != nil {
			return err
		}
	}
	if cfg.CryptoKey != "" {
		if err := options.CryptoKey.Set(cfg.CryptoKey); err != nil {
			return err
		}
	}
	return nil
}
