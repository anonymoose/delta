package main

import (
	"log"
)

var config Config

func main() {
	// config := InitConfig()
	// equityQuotes := LoadOHLCFiles(config)
	// log.Println(len(equityQuotes))
	// ProcessOptionsFile(config, equityQuotes)
	bstest()
}

func bstest() {

	S0 := 693.5

	K := 640.0
	right := "P"
	price := 2.4

	// K := 720.0
	// right := "C"
	// price := 4.5

	// vol := 0.2939
	// price := -1.0
	vol := -1.0

	r := 0.001 // risk free rate
	eval_date := "20151230"
	exp_date := "20160115"

	opt := NewOption(right, S0, K, eval_date, exp_date, r, vol, price)

	log.Println("CALL")
	log.Printf("Price: %v\n", opt.price)
	log.Printf("Delta: %v\n", opt.delta)
	log.Printf("Theta: %v\n", opt.theta)
	log.Printf("Gamma: %v\n", opt.gamma)
	log.Printf("Volatility: %v\n", opt.sigma)
}
