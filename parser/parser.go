package parser

import (
	"bufio"
	"io"
	"log"

	"github.com/jessejohnston/ProductIngester/product"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	// RecordLength is the expected length of each flat-file record
	RecordLength = 142

	// TaxRate is the product tax rate.
	TaxRate = 0.07775

	// NumberFieldLength is the expected length of all number fields.
	NumberFieldLength   = 8

	// CurrencyFieldLength is the expected length of all currency fields.
	CurrencyFieldLength = 8

	// FlagsFieldLength is the expected length of all flag fields.
	FlagsFieldLength    = 9
)

var (
	// ErrBadParameter is the error returned when invalid input is provided.
	ErrBadParameter = errors.New("Invalid parameter")
)

// Converter is the behavior of a type that converts fixed-length text values to other types.
//go:generate mockery -name Converter
type Converter interface {
	ToNumber(text []byte) (int, error)
	ToString(text []byte) string
	ToCurrency(text []byte) (decimal.Decimal, error)
	ToFlags(text []byte) (product.Flags, error)
}

// Parser reads from an input source, producing parsed records in it's Output channel.
type Parser struct {
	src     io.Reader
	convert Converter
	records chan *product.Record
	errors  chan error
	done    chan bool
}

// New creates a new product parser.
func New(input io.Reader, c Converter) (*Parser, error) {
	if input == nil || c == nil {
		return nil, ErrBadParameter
	}

	return &Parser{
		src:     input,
		convert: c,
		records: make(chan *product.Record),
		errors:  make(chan error),
		done:    make(chan bool),
	}, nil
}

// Parse reads each line from the input and sends parsed records to the output channel.
func (p *Parser) Parse() (<-chan *product.Record, <-chan error, <-chan bool) {
	// "go" runs p.execute() asynchronously so that the caller can start reading
	// records and errors off the returned channels.
	go p.execute()

	return p.records, p.errors, p.done
}

func (p *Parser) execute() {
	defer func() {
		close(p.done)
		close(p.records)
		close(p.errors)
	}()

	scanner := bufio.NewScanner(p.src)

	for row := 0; scanner.Scan(); row++ {
		data := scanner.Bytes()
		record, err := p.ParseRecord(row, data)
		if err != nil {
			log.Println(errors.WithStack(err))
			p.errors <- err
		} else {
			p.records <- record
		}
	}

	p.done <- true
}

func (p *Parser) ParseRecord(row int, text []byte) (*product.Record, error) {
	if len(text) != RecordLength {
		return nil, errors.WithStack(ErrBadParameter)
	}

	record := &product.Record{}
	var err error

	fragment := text[0:8]
	record.ID, err = p.convert.ToNumber(fragment)
	if err != nil {
		return nil, NewParserError(row, 0, fragment, "Error parsing ID", err)
	}

	fragment = text[9:68]
	record.Description = p.convert.ToString(fragment)

	fragment = text[69:77]
	singularPrice, err := p.convert.ToCurrency(fragment)
	if err != nil {
		return nil, NewParserError(row, 69, fragment, "Error parsing singular price", err)
	}

	// If singular price is zero, read the split price and use it instead.
	if singularPrice.Equal(decimal.Zero) {
		fragment = text[87:95]
		splitPrice, err := p.convert.ToCurrency(fragment)
		if err != nil {
			return nil, NewParserError(row, 87, fragment, "Error parsing split price", err)
		}

		fragment = text[105:113]
		forX, err := p.convert.ToNumber(fragment)
		if err != nil {
			return nil, NewParserError(row, 105, fragment, "Error parsing for X", err)
		}
		if forX == 0 {
			return nil, NewParserError(row, 105, fragment, "Error calculating split price (zero for X)", err)
		}

		// Round to 4 decimal places, half down
		record.Price = splitPrice.Div(decimal.New(int64(forX), 0)).RoundBank(4)
	} else {
		record.Price = singularPrice
	}

	fragment = text[78:86]
	singularPromoPrice, err := p.convert.ToCurrency(fragment)
	if err != nil {
		return nil, NewParserError(row, 78, fragment, "Error parsing singular promotional price", err)
	}

	// If singular promo price is zero, read the split promo price and use it instead.
	if singularPromoPrice.Equal(decimal.Zero) {
		fragment = text[96:104]
		splitPromoPrice, err := p.convert.ToCurrency(fragment)
		if err != nil {
			return nil, NewParserError(row, 96, text[96:104], "Error parsing split promo price", err)
		}

		if splitPromoPrice.GreaterThan(decimal.Zero) {
			fragment = text[114:122]
			promoForX, err := p.convert.ToNumber(fragment)
			if err != nil {
				return nil, NewParserError(row, 114, fragment, "Error parsing promo for X", err)
			}
			if promoForX == 0 {
				return nil, NewParserError(row, 114, fragment, "Error calculating promo split price (zero for X)", err)
			}
	
			// Round to 4 decimal places, half down
			record.PromoPrice = splitPromoPrice.Div(decimal.New(int64(promoForX), 0)).RoundBank(4)
		}
	} else {
		record.PromoPrice = singularPromoPrice
	}

	record.DisplayPrice = "$" + record.Price.StringFixed(2)
	record.PromoDisplayPrice = "$" + record.PromoPrice.StringFixed(2)

	fragment = text[123:132]
	flags, err := p.convert.ToFlags(fragment)
	if err != nil {
		return nil, NewParserError(row, 123, fragment, "Error parsing flags", err)
	}

	if flags.PerWeight() {
		record.Unit = product.UnitPound
	} else {
		record.Unit = product.UnitEach
	}

	if flags.Taxable() {
		record.TaxRate = decimal.NewFromFloat32(TaxRate)
	} else {
		record.TaxRate = decimal.Zero
	}

	record.Size = p.convert.ToString(text[133:142])

	return record, nil
}
