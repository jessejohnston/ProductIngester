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

	record.ID, err = p.convert.ToNumber(text[0:8])
	if err != nil {
		return nil, NewError(row, 0, text[0:8], "Error parsing ID", err)
	}

	record.Description = p.convert.ToString(text[9:68])

	record.Price, err = p.convert.ToCurrency(text[69:77])
	if err != nil {
		return nil, NewError(row, 69, text[69:77], "Error parsing price", err)
	}

	// If price is zero, read the split prices and use those instead.
	if record.Price.Equal(decimal.Zero) {
		splitPrice, err := p.convert.ToCurrency(text[87:95])
		if err != nil {
			return nil, NewError(row, 87, text[87:95], "Error parsing split price", err)
		}
		forX, err := p.convert.ToNumber(text[105:113])
		if err != nil {
			return nil, NewError(row, 105, text[87:95], "Error parsing for X", err)
		}
		record.Price = splitPrice.Div(decimal.New(int64(forX), 0))

		splitPromoPrice, err := p.convert.ToCurrency(text[96:104])
		if err != nil {
			return nil, NewError(row, 96, text[96:104], "Error parsing split promo price", err)
		}

		if splitPromoPrice.GreaterThan(decimal.Zero) {
			promoForX, err := p.convert.ToNumber(text[114:122])
			if err != nil {
				return nil, NewError(row, 114, text[114:122], "Error parsing promo for X", err)
			}

			record.PromoPrice = splitPromoPrice.Div(decimal.New(int64(promoForX), 0))
		}
	} else {
		record.PromoPrice, err = p.convert.ToCurrency(text[78:86])
		if err != nil {
			return nil, NewError(row, 78, text[78:86], "Error parsing promotional price", err)
		}
	}

	flags, err := p.convert.ToFlags(text[123:132])
	if err != nil {
		return nil, NewError(row, 123, text[123:132], "Error parsing flags", err)
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
