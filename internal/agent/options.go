package agent

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/FeshLig/metrcollector/internal/flags"
)

type Options struct {
	Address        flags.NetAddress
	ReportInterval flags.SecondsDuration
	PollInterval   flags.SecondsDuration
	Key            flags.Key
	RateLimit      flags.RateLimit
}

func GetOptions() Options {

	const defaultReportInterval = time.Duration(10) * time.Second
	const defaultPollInterval = time.Duration(2) * time.Second

	options := Options{
		Address: flags.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
		ReportInterval: flags.SecondsDuration{
			Duration: defaultReportInterval,
		},
		PollInterval: flags.SecondsDuration{
			Duration: defaultPollInterval,
		},
		Key:       "",
		RateLimit: flags.RateLimit(5),
	}

	parseFlags(&options)
	parseEnv(&options)

	return options

}

func parseFlags(options *Options) {

	flag.Var(&options.Address, "a", "net address host:port")
	flag.Var(&options.ReportInterval, "r", "frequency of sending metrics")
	flag.Var(&options.PollInterval, "p", "metrics polling frequency")
	flag.Var(&options.Key, "k", "hash key")
	flag.Var(&options.Key, "l", "rate limit")

	flag.Parse()

}

func parseEnv(options *Options) error {

	if addrStr, ok := os.LookupEnv("ADDRESS"); ok {
		err := options.Address.Set(addrStr)
		if err != nil {
			return fmt.Errorf("wrong value of ADDRESS: %w", err)
		}
	}
	if reportStr, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if err := options.ReportInterval.Set(reportStr); err != nil {
			return fmt.Errorf("wrong value of REPORT_INTERVAL: %w", err)
		}
	}
	if pollStr, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if err := options.PollInterval.Set(pollStr); err != nil {
			return fmt.Errorf("wrong value of POLL_INTERVAL: %w", err)
		}
	}
	if keyStr, ok := os.LookupEnv("KEY"); ok {
		err := options.Key.Set(keyStr)
		if err != nil {
			return fmt.Errorf("wrong value of KEY: %w", err)
		}
	}
	if rateStr, ok := os.LookupEnv("RATE_LIMIT"); ok {
		err := options.RateLimit.Set(rateStr)
		if err != nil {
			return fmt.Errorf("wrong value of RATE_LIMIT: %w", err)
		}
	}

	return nil

}
