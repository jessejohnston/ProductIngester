package main

import (
	"github.com/jessejohnston/ProductIngester/product"
)

// Parser defines the behavior of a product catalog parser.
type Parser interface {
	Parse() (<-chan *product.Record, <-chan error, <-chan bool)
}

// Database defines the behavior of a product catalog database.
type Database interface {
	InsertProductRecord(r *product.Record) error
}
