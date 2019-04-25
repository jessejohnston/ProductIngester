package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jessejohnston/ProductIngester/parser"
	"github.com/jessejohnston/ProductIngester/product"
	"github.com/pkg/errors"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		println("usage: ingest <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", filename, err)
	}
	defer file.Close()

	p, err := getParser(file)
	if err != nil {
		log.Fatalf("Error creating parser: %v", err)
	}

	// Start parsing, receiving a stream of records and parsing errors.
	records, errors, done := p.Parse()

	// As each record (or error) is generated, add the record to the database or log the error.
	var results []*product.Record

	for {
		select {
		case err := <-errors:
			log.Println(err)
		case r := <-records:
			fmt.Println(r)
			results = append(results, r)
		case <-done:
			log.Println("Done")
			os.Exit(0)
		}
	}
}

func getParser(input io.Reader) (Parser, error) {
	convert, err := getConverter()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return parser.New(input, convert)
}

func getConverter() (parser.Converter, error) {
	return product.NewConverter(parser.NumberFieldLength, parser.CurrencyFieldLength, parser.FlagsFieldLength)
}
