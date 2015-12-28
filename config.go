package main

import (
	"flag"
)

type Config struct {
	filePath string
}

func initConfig() *Config {
	config := &Config{filePath: ""}
	flag.StringVar(&config.filePath, "path", "", "Readable path to options EOD file")
	flag.Parse()
	return config
}
