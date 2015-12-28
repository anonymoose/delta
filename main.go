package main

import (
	"fmt"
)

var config Config

func main() {
	//config := initConfig()

	//processFile(config)

	S0 := 15.45             // Underlying price
	K := 17.0               // contract strike
	exp_date := "20160115"  // expiration date
	eval_date := "20151228" // "today"
	r := 0.001              // risk free rate
	vol := 0.525            // implied volatility

	opt := NewOption("C", S0, K, eval_date, exp_date, r, vol, -1)

	fmt.Println("CALL")
	fmt.Printf("Price: %v\n", opt.price)
	fmt.Printf("Delta: %v\n", opt.delta)
	fmt.Printf("Theta: %v\n", opt.theta)
	fmt.Printf("Gamma: %v\n", opt.gamma)

	opt = NewOption("C", S0, K, eval_date, exp_date, r, -1, 0.20)
	fmt.Printf("Volatility: %v\n", opt.sigma)

	// right = "P" // 'C' = call, 'P' = put
	// opt = NewOption(right, S0, K, eval_date, exp_date, r, vol, -1)

	// fmt.Printf("Price %v: %v\n", opt.right, opt.price)
	// fmt.Printf("Delta %v: %v\n", opt.right, opt.delta)
	// fmt.Printf("Theta %v: %v\n", opt.right, opt.theta)
	// fmt.Printf("Gamma %v: %v\n", opt.right, opt.gamma)

	// opt = NewOption("P", S0, K, eval_date, exp_date, r, -1, 1.7622)
	// fmt.Printf("Volatility %v: %v\n", opt.right, opt.sigma)

}
