package database

import (
	"fmt"
	"github.com/jessejohnston/ProductIngester/product"
)

// Db is a database implementation
type Db struct {
}

// New returns a new database instance.
func New() *Db {
	return &Db{}
}

// InsertProductRecord adds a product to the database, returning an error if the insertion fails.
func (db *Db) InsertProductRecord(r *product.Record) error {
	fmt.Println(r)
	return nil
}
