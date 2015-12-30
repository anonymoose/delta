package main

import (
	"flag"
)

type Config struct {
	ohlcFilePath string
	date         string
	deltaMin     float64
	deltaMax     float64
	priceMin     float64
	priceMax     float64
	volumeMin    int64
	oiMin        int64
}

func InitConfig() *Config {
	config := &Config{ohlcFilePath: ""}
	flag.StringVar(&config.ohlcFilePath, "ohlc", "", "Readable path to equity OHLC files")
	flag.StringVar(&config.date, "date", "", "Date stamp for 'today' (20151228)")

	flag.Float64Var(&config.deltaMin, "deltaMin", 10.0, "Minimum delta.")
	flag.Float64Var(&config.deltaMax, "deltaMax", 15.0, "Maximum delta.")
	flag.Float64Var(&config.priceMin, "priceMin", 0.15, "Minimum contract price.")
	flag.Float64Var(&config.priceMax, "priceMax", 0.75, "Maximum contract price.")
	flag.Int64Var(&config.volumeMin, "volumeMin", 25, "Minimum volume.")
	flag.Int64Var(&config.oiMin, "oiMin", 200, "Minimum open interest.")

	flag.Parse()
	return config
}
