package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	_ "github.com/lib/pq" // https://godoc.org/github.com/lib/pq
	"log"
	"os"
	"regexp"
	// "math"
	// "github.com/d4l3k/talib"
)

func DBOpen() *sql.DB {
	db, err := sql.Open("postgres", "user=pickr dbname=pickr password=pickr sslmode=disable")
	chkerr(err)
	return db
}

type Quote struct {
	symbol     string
	exchange   string
	dt         string
	open       float64
	high       float64
	low        float64
	close      float64
	vol        int64
	oi         int64 // only valid for options
	fullSymbol string
	right      string
	strike     float64
	expDate    string
}

type EquityName struct {
	symbol   string
	exchange string
	name     string
}

var nameInsertStmt *sql.Stmt
var SQL_NAME_INSERT = "insert into symbol_name (exchange, symbol, name) values ($1, $2, $3)"

func LoadNamesFiles(config *Config) map[string]*EquityName {
	namesCache := make(map[string]*EquityName)
	exchanges := [...]string{"NYSE", "NASDAQ", "AMEX", "INDEX"}

	db := DBOpen()
	defer db.Close()
	db.Exec("truncate table symbol_name")
	nameInsertStmt, err := db.Prepare(SQL_NAME_INSERT)
	chkerr(err)
	for _, exchange := range exchanges {
		path := config.dataPath + "/meta/src/names/" + exchange + ".txt"
		log.Printf("Loading %v...", path)
		f, err := os.Open(path)
		chkerr(err)
		defer f.Close()

		r := csv.NewReader(bufio.NewReader(f))
		r.Comma = '\t'

		for {
			record, err := r.Read()
			//chkerr(err)
			if err != nil {
				break
			}
			eqName := insertNameRecord(record, nameInsertStmt, exchange)
			namesCache[eqName.symbol] = eqName
		}
	}
	return namesCache
}

func insertNameRecord(record []string, insert *sql.Stmt, exchange string) *EquityName {
	_, err := insert.Exec(exchange, record[0], record[1])

	chkerr(err)
	return &EquityName{exchange: exchange, symbol: record[0], name: record[1]}
}

var equityInsertStmt *sql.Stmt
var SQL_EQUITY_INSERT_STMT = "insert into symbol_eod (exchange, symbol, dt, open, high, low, close, vol) values ($1, $2, $3, $4, $5, $6, $7, $8)"

var optionInsertStmt *sql.Stmt
var SQL_OPTION_INSERT_STMT = "insert into option_eod (exchange,symbol,symbol_full,option_type,dt,expiration_dt,strike,open,high,low,close,vol,oi) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)"

func LoadAllOHLCFiles(config *Config, namesCache map[string]*EquityName) {
	exchanges := [...]string{"NYSE", "NASDAQ", "AMEX", "INDEX", "OPRA"}

	db := DBOpen()
	defer db.Close()
	db.Exec("delete from symbol_eod where dt = $1", config.date)
	db.Exec("delete from option_eod where dt = $1", config.date)
	equityInsertStmt, err := db.Prepare(SQL_EQUITY_INSERT_STMT)
	chkerr(err)
	optionInsertStmt, err := db.Prepare(SQL_OPTION_INSERT_STMT)
	chkerr(err)
	for _, exchange := range exchanges {
		path := config.dataPath + "/quotes/db/" + exchange + "/" + exchange + "_" + config.date + ".txt"
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

			if exchange == "OPRA" {
				q := parseOptionQuote(record, exchange)
				insertOptionQuote(q, optionInsertStmt, namesCache)
			} else {
				q := parseEquityQuote(record, exchange)
				insertEquityQuote(q, equityInsertStmt)
			}
		}
	}
}

func insertEquityQuote(quote *Quote, insert *sql.Stmt) {
	_, err := insert.Exec(quote.exchange, quote.symbol, quote.dt, quote.open, quote.high, quote.low, quote.close, quote.vol)
	chkerr(err)
}

//  KB: [2015-12-31]: Some names don't map consistently for indexes.  Fixups are here.

var nameMappings = map[string]*EquityName{
	"NDX.X":  &EquityName{symbol: "NDX", exchange: "INDEX"},
	"RUT.X":  &EquityName{symbol: "RUTX", exchange: "INDEX"},
	"SPX.XO": &EquityName{symbol: "SPX", exchange: "INDEX"},
}

func insertOptionQuote(quote *Quote, insert *sql.Stmt, namesCache map[string]*EquityName) {
	name := namesCache[quote.symbol]
	if name == nil {
		name = nameMappings[quote.symbol]
	}

	if name != nil {
		_, err := insert.Exec(name.exchange, quote.symbol, quote.fullSymbol, quote.right, quote.dt, quote.expDate,
			quote.strike, quote.open, quote.high, quote.low, quote.close, quote.vol, quote.oi)
		if err != nil {
			log.Printf("%v / %v : %v\n", quote.fullSymbol, quote.dt, err)
		}
	}
}

// Take an option record and shove it into the database.
//  Sym              Exp Dt   Open  High   Low   Close  Vol
// [AAPL             20151015 3.15  3.15   3.15  3.15   100      ]
func parseEquityQuote(quote []string, exchange string) *Quote {
	//log.Println(quote)
	q := &Quote{
		symbol:   quote[0],
		exchange: exchange,
		dt:       quote[1],
		open:     parseFloat(quote[2]),
		high:     parseFloat(quote[3]),
		low:      parseFloat(quote[4]),
		close:    parseFloat(quote[5]),
		vol:      parseInt(quote[6]),
	}

	return q
}

// Take an option record and shove it into the database.
//  Sym              Exp Dt   Open  High   Low   Close  Vol    OpenInterest
// [A151016C00032500 20151015 3.15  3.15   3.15  3.15   0      37]
func parseOptionQuote(quote []string, exchange string) *Quote {
	sym, yr, mo, day, right, strikeDollars, strikeCents := parseOptionSymbol(quote[0])

	q := &Quote{
		symbol:     sym,
		exchange:   exchange,
		dt:         quote[1],
		open:       parseFloat(quote[2]),
		high:       parseFloat(quote[3]),
		low:        parseFloat(quote[4]),
		close:      parseFloat(quote[5]),
		vol:        parseInt(quote[6]),
		oi:         parseInt(quote[7]),
		fullSymbol: quote[0],
		right:      right,
		strike:     strikeDollars + (strikeCents * 0.001),
		expDate:    "20" + yr + mo + day,
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
	r, err := regexp.Compile("^([A-Z.-]+)([0-9]{2})([0-9]{2})([0-9]{2})([[:alpha:]]{1})([0-9]{5})([0-9]{3})")
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
