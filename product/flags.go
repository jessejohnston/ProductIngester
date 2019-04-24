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

// PerWeight returns true if the flags include FlagPerWeight
func (f Flags) PerWeight() bool {
	return f&FlagPerWeight == FlagPerWeight
}

// Taxable returns true if the flags include FlagTaxable
func (f Flags) Taxable() bool {
	return f&FlagTaxable == FlagTaxable
}
