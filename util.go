package main

import (
	"log"
	"strconv"
)

func chkerr(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
		panic(e)
	}
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	chkerr(err)
	return f
}

func parseInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	chkerr(err)
	return i
}
