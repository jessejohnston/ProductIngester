package parser

import (
	"strings"
	"testing"

	"github.com/jessejohnston/ProductIngester/product"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type parserTestSuite struct {
	suite.Suite
	converter Converter
}

func Test_Parser(t *testing.T) {
	s := new(parserTestSuite)
	suite.Run(t, s)
}

func (s *parserTestSuite) SetupSuite() {
	s.converter, _ = product.NewConverter(NumberFieldLength, CurrencyFieldLength, FlagsFieldLength)
}

func (s *parserTestSuite) Test_NewParser_NoInput_ReturnsError() {
	t := s.T()

	_, err := New(nil, s.converter)
	require.Error(t, err)
	require.Equal(t, ErrBadParameter, errors.Cause(err))
}

func (s *parserTestSuite) Test_NewParser_NoConverter_ReturnsError() {
	t := s.T()

	reader := strings.NewReader("the record")
	_, err := New(reader, nil)
	require.Error(t, err)
	require.Equal(t, ErrBadParameter, errors.Cause(err))
}

func (s *parserTestSuite) Test_ParseRecord_NilRecord_ReturnsError() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	_, err := p.ParseRecord(1, nil)
	require.Error(t, err)
	require.Equal(t, ErrBadParameter, errors.Cause(err))
}

func (s *parserTestSuite) Test_ParseRecord_ShortRecord_ReturnsError() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("80000001 Kimchi-flavored white rice                                  00000567 00000000 00000000 00000000 00000000 00000000 NNNN")
	_, err := p.ParseRecord(1, row)
	require.Error(t, err)
	require.Equal(t, ErrBadParameter, errors.Cause(err))
}

func (s *parserTestSuite) Test_ParseRecord_LongRecord_ReturnsError() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("80000001 Kimchi-flavored white rice                                  00000567 00000000 00000000 00000000 00000000 00000000 NNNNNNNNN      18ozzzzzz")
	_, err := p.ParseRecord(1, row)
	require.Error(t, err)
	require.Equal(t, ErrBadParameter, errors.Cause(err))
}

func (s *parserTestSuite) Test_ParseRecord_SingularPrice_NoPromoPrice_PriceIsSingularPrice() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("80000001 Kimchi-flavored white rice                                  00000567 00000000 00000000 00000000 00000000 00000000 NNNNNNNNN      18oz")
	r, err := p.ParseRecord(1, row)
	require.NoError(t, err)

	expectedPrice, _ := decimal.NewFromString("5.67")
	require.True(t, r.Price.Equals(expectedPrice))

	require.True(t, r.PromoPrice.Equals(decimal.Zero))
}

func (s *parserTestSuite) Test_ParseRecord_SingularPromoPrice_SplitPrice_PriceIsSplitPrice_PromoPriceIsSingularPrice() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("14963801 Generic Soda 12-pack                                        00000000 00000549 00001300 00000000 00000002 00000000 NNNNYNNNN   12x12oz")
	r, err := p.ParseRecord(1, row)
	require.NoError(t, err)

	expectedPrice, _ := decimal.NewFromString("6.50")
	require.True(t, r.Price.Equals(expectedPrice))

	expectedPromoPrice, _ := decimal.NewFromString("5.49")
	require.True(t, r.PromoPrice.Equals(expectedPromoPrice))
}

func (s *parserTestSuite) Test_ParseRecord_SingularPrice_SplitPromoPrice_PriceIsSingularPrice_PromoPriceIsSplitPrice() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("14963801 Generic Soda 12-pack                                        00000549 00000000 00000000 00001000 00000000 00000002 NNNNYNNNN   12x12oz")
	r, err := p.ParseRecord(1, row)
	require.NoError(t, err)

	expectedPrice, _ := decimal.NewFromString("5.49")
	require.True(t, r.Price.Equals(expectedPrice))

	expectedPromoPrice, _ := decimal.NewFromString("5.00")
	require.True(t, r.PromoPrice.Equals(expectedPromoPrice))
}

func (s *parserTestSuite) Test_ParseRecord_SingularPrice_SingularPromoPrice_PriceIsSingularPrice() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("40123401 Marlboro Cigarettes                                         00001000 00000549 00000000 00000000 00000000 00000000 YNNNNNNNN          ")
	r, err := p.ParseRecord(1, row)
	require.NoError(t, err)

	expectedPrice, _ := decimal.NewFromString("10.00")
	require.True(t, r.Price.Equals(expectedPrice))

	expectedPromoPrice, _ := decimal.NewFromString("5.49")
	require.True(t, r.PromoPrice.Equals(expectedPromoPrice))
}

func (s *parserTestSuite) Test_ParseRecord_PerWeightFlagNotSet_HasEachUnitOfMeasure() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("40123401 Marlboro Cigarettes                                         00001000 00000549 00000000 00000000 00000000 00000000 YNNNNNNNN          ")
	r, err := p.ParseRecord(1, row)
	require.NoError(t, err)

	require.Equal(t, product.UnitEach, r.Unit)
}

func (s *parserTestSuite) Test_ParseRecord_PerWeightFlagSet_HasPoundUnitOfMeasure() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("50133333 Fuji Apples (Organic)                                       00000349 00000000 00000000 00000000 00000000 00000000 NNYNNNNNN        lb")
	r, err := p.ParseRecord(1, row)
	require.NoError(t, err)

	require.Equal(t, product.UnitPound, r.Unit)
}

func (s *parserTestSuite) Test_ParseRecord_TaxableFlagNotSet_HasZeroTaxRate() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("80000001 Kimchi-flavored white rice                                  00000567 00000000 00000000 00000000 00000000 00000000 NNNNNNNNN      18oz")
	r, err := p.ParseRecord(1, row)
	require.NoError(t, err)

	require.Equal(t, decimal.Zero, r.TaxRate)
}

func (s *parserTestSuite) Test_ParseRecord_TaxableFlagSet_HasTaxRate() {
	t := s.T()

	reader := strings.NewReader("the file")
	p, _ := New(reader, s.converter)

	row := []byte("14963801 Generic Soda 12-pack                                        00000000 00000549 00001300 00000000 00000002 00000000 NNNNYNNNN   12x12oz")
	r, err := p.ParseRecord(1, row)
	require.NoError(t, err)

	expectedTaxRate := decimal.NewFromFloat32(TaxRate)
	require.Equal(t, expectedTaxRate, r.TaxRate)
}

func (s *parserTestSuite) Test_Parse_ReturnsAllRecords() {
	t := s.T()

	reader := strings.NewReader(
		"80000001 Kimchi-flavored white rice                                  00000567 00000000 00000000 00000000 00000000 00000000 NNNNNNNNN      18oz\n" +
			"14963801 Generic Soda 12-pack                                        00000000 00000549 00001300 00000000 00000002 00000000 NNNNYNNNN   12x12oz\n" +
			"40123401 Marlboro Cigarettes                                         00001000 00000549 00000000 00000000 00000000 00000000 YNNNNNNNN          \n" +
			"50133333 Fuji Apples (Organic)                                       00000349 00000000 00000000 00000000 00000000 00000000 NNYNNNNNN        lb")
	p, _ := New(reader, s.converter)

	records, errors, done := p.Parse()
	require.NotNil(t, records)
	require.NotNil(t, errors)
	require.NotNil(t, done)

	var results []*product.Record

	for {
		select {
		case <-errors:
			t.FailNow()
		case r := <-records:
			results = append(results, r)
		case <-done:
			goto finished
		}
	}

finished:
	require.Len(t, results, 4)
	require.Equal(t, 80000001, results[0].ID)
	require.Equal(t, 14963801, results[1].ID)
	require.Equal(t, 40123401, results[2].ID)
	require.Equal(t, 50133333, results[3].ID)
}

func (s *parserTestSuite) Test_Parse_BadRecord_ReturnsOtherRecords_AndError() {
	t := s.T()

	reader := strings.NewReader(
		"80000001 Kimchi-flavored white rice                                  00000567 00000000 00000000 00000000 00000000 00000000 NNNNNNNNN      18oz\n" +
			"14963801 Generic Soda 12-pack                                        00000000 00000549 00001300 00000000 00000002 00000000 NNNNYNNNN   12x12oz\n" +
			"40123401 Marlboro Cigare\n" +
			"50133333 Fuji Apples (Organic)                                       00000349 00000000 00000000 00000000 00000000 00000000 NNYNNNNNN        lb")
	p, _ := New(reader, s.converter)

	records, errors, done := p.Parse()
	require.NotNil(t, records)
	require.NotNil(t, errors)
	require.NotNil(t, done)

	var results []*product.Record
	var errs []error

	for {
		select {
		case e := <-errors:
			errs = append(errs, e)
		case r := <-records:
			results = append(results, r)
		case <-done:
			goto finished
		}
	}

finished:
	require.Len(t, results, 3)
	require.Equal(t, 80000001, results[0].ID)
	require.Equal(t, 14963801, results[1].ID)
	require.Equal(t, 50133333, results[2].ID)
	require.Len(t, errs, 1)
}
