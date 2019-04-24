package main

import (
	"io"
	"log"
	"os"

	"github.com/jessejohnston/ProductIngester/database"
	"github.com/jessejohnston/ProductIngester/parser"
	"github.com/jessejohnston/ProductIngester/product"
	"github.com/pkg/errors"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println("usage: ingest <filename>")
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

	db := getDatabaseWriter()

	// Start parsing, receiving a stream of records and parsing errors.
	records, errors, done := p.Parse()

	// As each record (or error) is generated, add the record to the database or log the error.
	for {
		select {
		case err := <-errors:
			log.Println(err)
		case r := <-records:
			err := db.InsertProductRecord(r)
			if err != nil {
				log.Println(err)
			}
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
	return product.NewConverter(8, 8, 9)
}

func getDatabaseWriter() DatabaseWriter {
	return database.New()
}
