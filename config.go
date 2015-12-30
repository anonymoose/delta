package main

import (
	"flag"
)

type Config struct {
	ohlcFilePath string
	date         string
	deltaMin     float64
	deltaMax     float64
	volumeMin    int64
	oiMin        int64
}

func InitConfig() *Config {
	config := &Config{ohlcFilePath: ""}
	flag.StringVar(&config.ohlcFilePath, "ohlc", "", "Readable path to equity OHLC files")
	flag.StringVar(&config.date, "date", "", "Date stamp for 'today' (20151228)")

	flag.Float64Var(&config.deltaMin, "deltaMin", 10.0, "Minimum delta.")
	flag.Float64Var(&config.deltaMax, "deltaMax", 15.0, "Maximum delta.")
	flag.Int64Var(&config.volumeMin, "volumeMin", 100, "Minimum volume.")
	flag.Int64Var(&config.oiMin, "oiMin", 1000, "Minimum open interest.")

	flag.Parse()
	return config
}
