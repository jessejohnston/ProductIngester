package product

import (
	"github.com/shopspring/decimal"
)

// Record is the parsed Product
type Record struct {
	ID                int
	Description       string
	DisplayPrice      string
	Price             decimal.Decimal
	PromoDisplayPrice string
	PromoPrice        decimal.Decimal
	Unit              UnitOfMeasure
	Size              string
	TaxRate           decimal.Decimal
}
