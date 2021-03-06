package product

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Converter provides format conversions for fixed-length fields.
type Converter struct {
	numberLength   int
	currencyLength int
	flagsLength    int
}

var (
	// ErrBadFieldLength is the error returned when a conversion method receives text of an unexpected length.
	ErrBadFieldLength = errors.New("Unexpected field length")

	// ErrBadParameter is the error returned when invalid input is provided.
	ErrBadParameter = errors.New("Invalid parameter")

	// ErrBadFormat is the error returned when a data field includes unexpected content
	ErrBadFormat = errors.New("Bad format")
)

// NewConverter returns a new converter of fixed length text fields.
func NewConverter(numFieldLength, currencyFieldLength, flagFieldLength int) (*Converter, error) {
	if numFieldLength < 1 || currencyFieldLength < 1 || flagFieldLength < 1 {
		return nil, errors.WithStack(ErrBadParameter)
	}

	return &Converter{
		numberLength:   numFieldLength,
		currencyLength: currencyFieldLength,
		flagsLength:    flagFieldLength,
	}, nil
}

// ToNumber converts text to an integer.
func (c *Converter) ToNumber(text []byte) (int, error) {
	if len(text) != c.numberLength {
		return 0, errors.WithStack(ErrBadFieldLength)
	}
	num, err := strconv.Atoi(string(text))
	if err != nil {
		return 0, errors.WithStack(ErrBadFormat)
	}

	return num, nil
}

// ToString converts text to a string.
func (c *Converter) ToString(text []byte) string {
	return strings.TrimSpace(string(text))
}

// ToCurrency converts text to a decimal value.
func (c *Converter) ToCurrency(text []byte) (decimal.Decimal, error) {
	if len(text) != c.currencyLength {
		return decimal.Decimal{}, errors.WithStack(ErrBadFieldLength)
	}
	unscaled, err := decimal.NewFromString(string(text))
	if err != nil {
		return decimal.Decimal{}, errors.WithStack(ErrBadFormat)
	}
	return unscaled.Shift(-2), nil
}

// ToFlags converts text to a set of flags.
func (c *Converter) ToFlags(text []byte) (Flags, error) {
	if len(text) != c.flagsLength {
		return FlagNone, errors.WithStack(ErrBadFieldLength)
	}

	var flags Flags

	for i, b := range text {
		if b != 'Y' && b != 'N' {
			return 0, errors.WithStack(ErrBadFormat)
		}
		if b == 'Y' {
			switch i {
			case 2:
				flags |= FlagPerWeight
			case 4:
				flags |= FlagTaxable
			}
		}
	}

	return flags, nil
}
