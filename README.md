# ProductIngester
[![Build Status](https://travis-ci.org/jessejohnston/ProductIngester.svg?branch=master)](https://travis-ci.org/jessejohnston/ProductIngester)  
[Latest Build](https://travis-ci.org/jessejohnston/ProductIngester/branches)

## Store product catalog ingestion package

This project includes the `parser` package which parses store catalog files.
The `product` package includes definitions of product data types and the `Converter` type that is used to convert flat-file text fields to
product data types.

The sample program `cmd/ingester` exercises the parser against a specific input file.

## To Build

Install Go: https://golang.org/doc/install

Install the mocking library:
```
$ go get github.com/vektra/mockery/.../
```
 
From the `cmd/ingester` folder:
```
$ go build
$ ./ingester input-sample.txt
```

The parsed product records are logged to the console.

## Catalog Parser Usage

A parser can be created from an input data source (io.Reader) and a field converter.

The `Parser.Parse()` method returns three output channels that produce product records, errors, and a done (EOF) signal.
```
// Get a reader to the input source.
file, err := os.Open(filename)
if err != nil {
	log.Fatalf("Error opening file %s: %v", filename, err)
}
defer file.Close()

// Create a field converter based on the known field lengths for each data type.
converter, err := product.NewConverter(parser.NumberFieldLength, parser.CurrencyFieldLength, parser.FlagsFieldLength)
if err != nil {
	log.Fatalf("Error creating field converter: %v", err)
}

// Create the parser.
parser, err := parser.New(file, converter)
if err != nil {
	log.Fatalf("Error creating parser: %v", err)
}

// Start parsing, receiving a stream of records and parsing errors.
records, errors, done := parser.Parse()

// As each record (or error) is generated, add the record to the results array or log the error.
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

// Do something with the parsed records.
...
```