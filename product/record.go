package product

import (
	"fmt"

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

func (r Record) String() string {
	return fmt.Sprintf("%d %60s %10s %10s %7v %s %8s", r.ID, r.Description, r.DisplayPrice, r.PromoDisplayPrice, r.Unit, r.Size, r.TaxRate.StringFixed(4))
}
