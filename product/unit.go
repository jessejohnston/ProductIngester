package product

type UnitOfMeasure string

const (
	// UnitEach expresses a per-item pricing unit of measure
	UnitEach UnitOfMeasure = "Each"

	// UnitPound expresses a weight-based pricing unit of measure
	UnitPound UnitOfMeasure = "Pound"
)
