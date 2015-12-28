package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

func processFile(config *Config) {
	f, err := os.Open(config.filePath)
	chkerr(err)
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))

	for {
		record, err := r.Read()
		chkerr(err)
		q := parseQuote(record)

		processQuote(q)
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
	q := &Quote{
		symbol: quote[0],
		dt:     quote[1],
		open:   parseFloat(quote[2]),
		high:   parseFloat(quote[3]),
		low:    parseFloat(quote[4]),
		close:  parseFloat(quote[5]),
		vol:    parseInt(quote[6]),
		oi:     parseInt(quote[7]),
	}

	return q
}

func processQuote(q *Quote) {
	log.Printf("%s %f", q.symbol, q.close)
}
