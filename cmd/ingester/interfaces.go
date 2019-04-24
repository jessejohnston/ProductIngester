package main

import (
	"github.com/jessejohnston/ProductIngester/product"
)

// Parser defines the behavior of a product catalog parser.
type Parser interface {
	Parse() (<-chan *product.Record, <-chan error, <-chan bool)
}

// DatabaseWriter defines the behavior of inserting into a product catalog database.
type DatabaseWriter interface {
	InsertProductRecord(r *product.Record) error
}
