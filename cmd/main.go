package main

import (
	"io"
	"log"
	"os"

	"github.com/jessejohnston/ProductIngester/database"
	"github.com/jessejohnston/ProductIngester/parser"
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

	db := getDatabase()

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
	return parser.New(input)
}

func getDatabase() Database {
	return database.New()
}
