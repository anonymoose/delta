package main

import (
	_ "bufio"
	_ "encoding/csv"
	"log"
	"math"
	_ "os"
)

func ProcessOptionsFile(config *Config, equityQuotes map[string]*Quote) {
	// path := config.ohlcFilePath + "/OPRA/OPRA_" + config.date + ".txt"
	// log.Printf("Loading %v...", path)
	// f, err := os.Open(path)
	// chkerr(err)
	// defer f.Close()

	// r := csv.NewReader(bufio.NewReader(f))

	// for {
	// 	record, err := r.Read()
	// 	if err != nil {
	// 		break
	// 	}
	// 	// q := parseOptionQuote(record)
	// 	// processOptionQuote(config, q, equityQuotes)
	// }
}

// Given an option quote from an EOD file, scan for options that match our criteria.
//    done: volume > config.volumeMin
//    done: expiration date == config.expDate
//    done: config.deltaMin < delta > config.deltaMax
//    done: open interest > config.oiMin
//    done: config.priceMin < contract price > config.priceMax
//    todo: IV Rank > config.ivMin
//    todo: true delta
//    todo: min distance from strike
//    todo: talib indicators on underlying (MACD, CCI, RSI)
func processOptionQuote(config *Config, optionQuote *Quote, equityQuotes map[string]*Quote) {
	sym, yr, mo, day, right, strikeDollars, strikeCents := parseOptionSymbol(optionQuote.symbol)

	if equityQuotes[sym] == nil || sym == "" || optionQuote.vol < config.volumeMin {
		// can't find it in symbols.  skip it.
		return
	}

	S0 := equityQuotes[sym].close              // underlying current
	K := strikeDollars + (strikeCents * 0.001) // contract strike
	exp_date := "20" + yr + mo + day           // expiration date
	eval_date := config.date                   // "today"
	r := 0.001                                 // risk free rate

	// if we care about expiration date, then filter out non-matches here.
	if config.expDate != "" && exp_date != config.expDate {
		return
	}

	opt := NewOption(right, S0, K, eval_date, exp_date, r, -1, optionQuote.close)

	aDelta := math.Abs(opt.delta) * 100

	if aDelta > config.deltaMin && aDelta < config.deltaMax && optionQuote.oi > config.oiMin && optionQuote.close > config.priceMin && optionQuote.close < config.priceMax {
		if math.IsNaN(opt.sigma) {
			log.Printf("nan")
		}
		log.Printf("%v\t%v\t%v\t%v\t%v\t%v underlying: %f strike: %f price: %f delta: %f vol: %f symbol: %v\n",
			optionQuote.vol, optionQuote.oi, sym, right, K, exp_date, S0, K, optionQuote.close, opt.delta, opt.sigma, optionQuote.symbol)
	}
	// fmt.Println("CALL")
	// fmt.Printf("Price: %v\n", opt.price)
	// fmt.Printf("Delta: %v\n", opt.delta)
	// fmt.Printf("Theta: %v\n", opt.theta)
	// fmt.Printf("Gamma: %v\n", opt.gamma)

	// opt = NewOption("C", S0, K, eval_date, exp_date, r, -1, 0.20)
	// fmt.Printf("Volatility: %v\n", opt.sigma)

	// log.Printf("%s %f", q.symbol, q.close)
}
