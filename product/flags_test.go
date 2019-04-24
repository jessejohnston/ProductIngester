package product

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_DefaultValue_FlagNone(t *testing.T) {
	var f Flags
	require.Equal(t, FlagNone, f)
}

func Test_DefaultValue_NotPerWeight(t *testing.T) {
	var f Flags
	require.False(t, f.PerWeight())
}

func Test_DefaultValue_NotTaxable(t *testing.T) {
	var f Flags
	require.False(t, f.PerWeight())
}

func Test_FlagPerWeight_IsPerWeight(t *testing.T) {
	f := FlagPerWeight
	require.True(t, f.PerWeight())
}

func Test_FlagPerWeightAndTaxable_IsPerWeight(t *testing.T) {
	f := FlagPerWeight | FlagTaxable
	require.True(t, f.PerWeight())
}

func Test_FlagTaxable_IsTaxable(t *testing.T) {
	f := FlagTaxable
	require.True(t, f.Taxable())
}

func Test_FlagPerWeightAndTaxable_IsTaxable(t *testing.T) {
	f := FlagPerWeight | FlagTaxable
	require.True(t, f.Taxable())
}
