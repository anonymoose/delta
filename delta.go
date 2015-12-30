package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"math"
	"os"
	"regexp"
)

func LoadOHLCFiles(config *Config) map[string]*Quote {
	equityQuotes := make(map[string]*Quote)
	exchanges := [...]string{"NYSE", "NASDAQ", "AMEX"}

	for _, exchange := range exchanges {
		path := config.ohlcFilePath + "/" + exchange + "/" + exchange + "_" + config.date + ".txt"
		log.Printf("Loading %v...", path)
		f, err := os.Open(path)
		chkerr(err)
		defer f.Close()

		r := csv.NewReader(bufio.NewReader(f))

		for {
			record, err := r.Read()
			//chkerr(err)
			if err != nil {
				break
			}
			q := parseQuote(record)
			equityQuotes[q.symbol] = q
		}
	}
	return equityQuotes
}

func ProcessOptionsFile(config *Config, equityQuotes map[string]*Quote) {
	path := config.ohlcFilePath + "/OPRA/OPRA_" + config.date + ".txt"
	log.Printf("Loading %v...", path)
	f, err := os.Open(path)
	chkerr(err)
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))

	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		q := parseQuote(record)
		processOptionQuote(config, q, equityQuotes)
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
	//log.Println(quote)
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

/*  AAPL131206C00400000
    + Symbol
    + Expiration Year(yy)
    + Expiration Month(mm)
    + Expiration Day(dd)
    + Call/Put Indicator (C or P)
    + Strike Price Dollars
    + Strike Price Fraction of Dollars (which can include decimals)

    'Understanding The 2010 Options Symbology'
     http://www.investopedia.com/articles/optioninvestor/10/options-symbol-rules.asp#ixzz3olRdqwvV
*/
func parseOptionSymbol(symbol string) (sym string, yr string, mo string, day string, right string, strikeDollars float64, strikeCents float64) {
	r, err := regexp.Compile("^([[:alpha:]]+)([0-9]{2})([0-9]{2})([0-9]{2})([[:alpha:]]{1})([0-9]{5})([0-9]{3})")
	if err != nil {
		panic(err)
	}
	// "AAPL131206C00400000" -> [AAPL131206C00400000 AAPL 13 12 06 C 00400 000]
	matches := r.FindStringSubmatch(symbol)
	if len(matches) == 8 {
		sym = matches[1]
		yr = matches[2]
		mo = matches[3]
		day = matches[4]
		right = matches[5]
		strikeDollars = parseFloat(matches[6])
		strikeCents = parseFloat(matches[7])
		return sym, yr, mo, day, right, strikeDollars, strikeCents
	}
	return "", "", "", "", "", -1.0, -1.0
}

func processOptionQuote(config *Config, optionQuote *Quote, equityQuotes map[string]*Quote) {
	sym, yr, mo, day, right, strikeDollars, strikeCents := parseOptionSymbol(optionQuote.symbol)

	if equityQuotes[sym] == nil || sym == "" {
		// can't find it in symbols.  skip it.
		return
	}

	S0 := equityQuotes[sym].close              // underlying current
	K := strikeDollars + (strikeCents * 0.001) // contract strike
	exp_date := "20" + yr + mo + day           // expiration date
	eval_date := config.date                   // "today"
	r := 0.001                                 // risk free rate

	opt := NewOption(right, S0, K, eval_date, exp_date, r, -1, optionQuote.close)

	//aDelta := math.Abs(opt.delta) * 100

	// if aDelta > config.deltaMin && aDelta < config.deltaMax && optionQuote.vol > config.volumeMin && optionQuote.oi > config.oiMin {
	//if ("P" == right && K > S0) || ("C" == right && K < S0) {
	//if "P" == right && S0 > K {
	if math.IsNaN(opt.sigma) {
		log.Printf("nan")
	}
	log.Printf("%v %v %v %v underlying: %f strike: %f price: %f delta: %f vol: %f symbol: %v\n",
		sym, right, K, exp_date, S0, K, optionQuote.close, opt.delta, opt.sigma, optionQuote.symbol)
	//}
	//}
	// fmt.Println("CALL")
	// fmt.Printf("Price: %v\n", opt.price)
	// fmt.Printf("Delta: %v\n", opt.delta)
	// fmt.Printf("Theta: %v\n", opt.theta)
	// fmt.Printf("Gamma: %v\n", opt.gamma)

	// opt = NewOption("C", S0, K, eval_date, exp_date, r, -1, 0.20)
	// fmt.Printf("Volatility: %v\n", opt.sigma)

	// log.Printf("%s %f", q.symbol, q.close)
}
