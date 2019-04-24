package parser

import (
	"bufio"
	"io"
	"log"

	"github.com/jessejohnston/ProductIngester/product"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// TaxRate is the product tax rate.
const TaxRate = 0.07775

var (
	// ErrBadParameter is the error returned when invalid input is provided.
	ErrBadParameter = errors.New("Invalid parameter")
)

// Converter is the behavior of a type that converts fixed-length text values to other types.
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
	// "go" runs p.execute() asynchronously
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
		record, err := p.toRecord(row, data)
		if err != nil {
			log.Println(errors.WithStack(err))
			p.errors <- err
		} else {
			p.records <- record
		}
	}

	p.done <- true
}

func (p *Parser) toRecord(row int, text []byte) (*product.Record, error) {
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
	record.Price, err = p.convert.ToCurrency(fragment)
	if err != nil {
		return nil, NewParserError(row, 69, fragment, "Error parsing price", err)
	}

	// If price is zero, read the split prices and use those instead.
	if record.Price.Equal(decimal.Zero) {
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
		record.Price = splitPrice.Div(decimal.New(int64(forX), 0))

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
	
			record.PromoPrice = splitPromoPrice.Div(decimal.New(int64(promoForX), 0))
		}
	} else {
		fragment = text[78:86]
		record.PromoPrice, err = p.convert.ToCurrency(fragment)
		if err != nil {
			return nil, NewParserError(row, 78, fragment, "Error parsing promotional price", err)
		}
	}

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
	}

	record.Size = p.convert.ToString(text[133:142])

	return record, nil
}
