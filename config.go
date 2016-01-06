package main

import (
	"flag"
)

type Config struct {
	action    string
	dataPath  string
	date      string
	deltaMin  float64
	deltaMax  float64
	priceMin  float64
	priceMax  float64
	volumeMin int64
	oiMin     int64
	expDate   string
}

func InitConfig() *Config {
	config := &Config{}
	flag.StringVar(&config.action, "action", "optionScan", "What do you want to do?")
	flag.StringVar(&config.dataPath, "dataPath", "./data", "Readable path to data file root")
	flag.StringVar(&config.date, "date", "", "Date stamp for 'today' (20151228)")

	flag.Float64Var(&config.deltaMin, "deltaMin", 10.0, "Minimum delta.")
	flag.Float64Var(&config.deltaMax, "deltaMax", 15.0, "Maximum delta.")
	flag.Float64Var(&config.priceMin, "priceMin", 0.15, "Minimum contract price.")
	flag.Float64Var(&config.priceMax, "priceMax", 0.75, "Maximum contract price.")
	flag.Int64Var(&config.volumeMin, "volumeMin", 25, "Minimum volume.")
	flag.Int64Var(&config.oiMin, "oiMin", 200, "Minimum open interest.")
	flag.StringVar(&config.expDate, "expDate", "", "Expiration date (optional)")

	flag.Parse()
	return config
}
