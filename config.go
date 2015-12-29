package main

import (
	"flag"
)

type Config struct {
	ohlcFilePath string
	date         string
}

func InitConfig() *Config {
	config := &Config{ohlcFilePath: ""}
	flag.StringVar(&config.ohlcFilePath, "ohlc", "", "Readable path to equity OHLC files")
	flag.StringVar(&config.date, "date", "", "Date stamp for 'today' (20151228)")

	flag.Parse()
	return config
}
