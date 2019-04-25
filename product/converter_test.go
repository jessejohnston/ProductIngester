package product

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type converterTestSuite struct {
	suite.Suite
	convert *Converter
}

func Test_Converter(t *testing.T) {
	s := new(converterTestSuite)
	suite.Run(t, s)
}

func (s *converterTestSuite) SetupSuite() {
	s.convert, _ = NewConverter(8, 8, 9)
}

func (s *converterTestSuite) Test_NewConverter_ReturnsConverter() {
	c, err := NewConverter(1, 1, 1)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), c)
}

func (s *converterTestSuite) Test_NewConverter_BadNumberFieldLength_ReturnsError() {
	_, err := NewConverter(0, 1, 1)
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadParameter, errors.Cause(err))
}

func (s *converterTestSuite) Test_NewConverter_BadCurrencyFieldLength_ReturnsError() {
	_, err := NewConverter(1, 0, 1)
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadParameter, errors.Cause(err))
}

func (s *converterTestSuite) Test_NewConverter_BadFlagsFieldLength_ReturnsError() {
	_, err := NewConverter(1, 1, 0)
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadParameter, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToNumber_ShortField_ReturnsError() {
	_, err := s.convert.ToNumber([]byte("001"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFieldLength, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToNumber_LongField_ReturnsError() {
	_, err := s.convert.ToNumber([]byte("000100000"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFieldLength, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToNumber_ZeroLeftPadded_ReturnsNumber() {
	num, err := s.convert.ToNumber([]byte("00000001"))
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, num)
}

func (s *converterTestSuite) Test_ToNumber_SpaceLeftPadded_ReturnsError() {
	_, err := s.convert.ToNumber([]byte("       1"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFormat, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToNumber_ZeroRightadded_ReturnsNumber() {
	num, err := s.convert.ToNumber([]byte("10000000"))
	require.NoError(s.T(), err)
	require.Equal(s.T(), 10000000, num)
}

func (s *converterTestSuite) Test_ToNumber_ReturnsNumber() {
	num, err := s.convert.ToNumber([]byte("12345678"))
	require.NoError(s.T(), err)
	require.Equal(s.T(), 12345678, num)
}

func (s *converterTestSuite) Test_ToNumber_Negative_ReturnsNumber() {
	num, err := s.convert.ToNumber([]byte("-1234567"))
	require.NoError(s.T(), err)
	require.Equal(s.T(), -1234567, num)
}

func (s *converterTestSuite) Test_ToNumber_Decimal_ReturnsError() {
	_, err := s.convert.ToNumber([]byte("12345.67"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFormat, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToNumber_NonNumeric_ReturnsError() {
	_, err := s.convert.ToNumber([]byte("12345ABC"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFormat, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToString_Empty_ReturnsEmpty() {
	str := s.convert.ToString([]byte{})
	require.Empty(s.T(), str)
}

func (s *converterTestSuite) Test_ToString_ReturnsString() {
	str := s.convert.ToString([]byte("delicious pickles"))
	require.Equal(s.T(), "delicious pickles", str)
}

func (s *converterTestSuite) Test_ToString_WithLeadingTrailingWhitespace_ReturnsTrimmedString() {
	str := s.convert.ToString([]byte("  delicious pickles  "))
	require.Equal(s.T(), "delicious pickles", str)
}

func (s *converterTestSuite) Test_ToCurrency_ShortField_ReturnsError() {
	_, err := s.convert.ToCurrency([]byte("001"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFieldLength, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToCurrency_LongField_ReturnsError() {
	_, err := s.convert.ToCurrency([]byte("000100000"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFieldLength, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToCurrency_ZeroLeftPadded_ReturnsValue() {
	cur, err := s.convert.ToCurrency([]byte("00000001"))
	require.NoError(s.T(), err)

	expected, _ := decimal.NewFromString("0.01")
	require.True(s.T(), cur.Equal(expected))
}

func (s *converterTestSuite) Test_ToCurrency_SpaceLeftPadded_ReturnsError() {
	_, err := s.convert.ToCurrency([]byte("       1"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFormat, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToCurrency_WholeDollars_ReturnsValue() {
	cur, err := s.convert.ToCurrency([]byte("00009900"))
	require.NoError(s.T(), err)

	expected, _ := decimal.NewFromString("99.00")
	require.True(s.T(), cur.Equal(expected))
}

func (s *converterTestSuite) Test_ToCurrency_DollarsAndCents_ReturnsValue() {
	cur, err := s.convert.ToCurrency([]byte("00001999"))
	require.NoError(s.T(), err)

	expected, _ := decimal.NewFromString("19.99")
	require.True(s.T(), cur.Equal(expected))
}

func (s *converterTestSuite) Test_ToCurrency_Negative_ReturnsValue() {
	cur, err := s.convert.ToCurrency([]byte("-0001999"))
	require.NoError(s.T(), err)

	expected, _ := decimal.NewFromString("-19.99")
	require.True(s.T(), cur.Equal(expected))
}

func (s *converterTestSuite) Test_ToFlags_ShortField_ReturnsError() {
	_, err := s.convert.ToFlags([]byte("YYN"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFieldLength, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToFlags_LongField_ReturnsError() {
	_, err := s.convert.ToCurrency([]byte("YYYYNNNNYN"))
	require.Error(s.T(), err)
	require.Equal(s.T(), ErrBadFieldLength, errors.Cause(err))
}

func (s *converterTestSuite) Test_ToFlags_AllNo_ReturnsFlagNone() {
	flags, err := s.convert.ToFlags([]byte("NNNNNNNNN"))
	require.NoError(s.T(), err)
	require.Equal(s.T(), FlagNone, flags)
}

func (s *converterTestSuite) Test_ToFlags_PerWeight_ReturnsFlagPerWeight() {
	flags, err := s.convert.ToFlags([]byte("NNYNNNNNN"))
	require.NoError(s.T(), err)
	require.Equal(s.T(), FlagPerWeight, flags)
}

func (s *converterTestSuite) Test_ToFlags_Taxable_ReturnsFlagTaxable() {
	flags, err := s.convert.ToFlags([]byte("NNNNYNNNN"))
	require.NoError(s.T(), err)
	require.Equal(s.T(), FlagTaxable, flags)
}

func (s *converterTestSuite) Test_ToFlags_PerWeightTaxable_ReturnsBoth() {
	flags, err := s.convert.ToFlags([]byte("NNYNYNNNN"))
	require.NoError(s.T(), err)
	require.Equal(s.T(), FlagPerWeight|FlagTaxable, flags)
}
