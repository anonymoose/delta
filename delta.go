package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

func LoadOHLCFiles(config *Config) map[string]*Quote {
	prices := make(map[string]*Quote)
	exchanges := [...]string{"NYSE", "NASDAQ", "AMEX"}

	for _, exchange := range exchanges {
		path := config.ohlcFilePath + "/" + exchange + "/" + exchange + "_" + config.date + ".txt"
		log.Printf("Loading %v...", path)
		f, _ := os.Open(path)
		//chkerr(err)
		defer f.Close()

		r := csv.NewReader(bufio.NewReader(f))

		for {
			record, _ := r.Read()
			//chkerr(err)
			q := parseQuote(record)
			prices[q.symbol] = q
		}
	}
	return prices
}

func ProcessOptionsFile(config *Config) {
	path := config.ohlcFilePath + "/OPRA/OPRA_" + config.date + ".txt"
	log.Printf("Loading %v...", path)
	f, _ := os.Open(path)
	//chkerr(err)
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))

	for {
		record, _ := r.Read()
		//chkerr(err)
		q := parseQuote(record)
		processOptionQuote(q)
	}

}

type Quote struct {
	symbol string
	dt     string
	open   float64
	high   float64
	low    float64
	close  float64
	vol    int64
	oi     int64
}

// Take an option record and shove it into the database.
//  Sym              Exp Dt   Open  High   Low   Close  Vol    OpenInterest
// [A151016C00032500 20151015 3.15  3.15   3.15  3.15   0      37]
func parseQuote(quote []string) *Quote {
	log.Println(quote)
	q := &Quote{
		symbol: quote[0],
		dt:     quote[1],
		open:   parseFloat(quote[2]),
		high:   parseFloat(quote[3]),
		low:    parseFloat(quote[4]),
		close:  parseFloat(quote[5]),
		vol:    parseInt(quote[6]),
	}

	if len(quote) == 8 {
		q.oi = parseInt(quote[7])
	}
	return q
}

func processOptionQuote(q *Quote) {
	// S0 := 15.45             // Underlying price
	// K := 17.0               // contract strike
	// exp_date := "20160115"  // expiration date
	// eval_date := "20151228" // "today"
	// r := 0.001              // risk free rate
	// vol := 0.525            // implied volatility

	// opt := NewOption("C", S0, K, eval_date, exp_date, r, vol, -1)

	// fmt.Println("CALL")
	// fmt.Printf("Price: %v\n", opt.price)
	// fmt.Printf("Delta: %v\n", opt.delta)
	// fmt.Printf("Theta: %v\n", opt.theta)
	// fmt.Printf("Gamma: %v\n", opt.gamma)

	// opt = NewOption("C", S0, K, eval_date, exp_date, r, -1, 0.20)
	// fmt.Printf("Volatility: %v\n", opt.sigma)

	// log.Printf("%s %f", q.symbol, q.close)
}
