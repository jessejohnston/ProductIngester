package product

// Flags specify boolean product characteristics.
type Flags uint

const (
	// FlagNone indicates no special product characteristics.
	FlagNone Flags = 0

	// FlagPerWeight indicates a product that is priced by weight.
	FlagPerWeight Flags = 1

	// FlagTaxable indicates a product that is taxable.
	FlagTaxable Flags = 2
)
