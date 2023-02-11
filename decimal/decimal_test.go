package decimal_test

import (
	"math"
	"testing"

	"github.com/stefanovazzocell/GoDecimal/decimal"
)

func TestDecimalFrom(t *testing.T) {
	// Int
	testCasesInt := map[int64]struct {
		positive bool
		value    uint64
	}{
		0:             {true, 0},
		-1:            {false, 1},
		1:             {true, 1},
		math.MaxInt64: {true, 9223372036854775807},
		math.MinInt64: {false, 9223372036854775808},
	}
	for i, expect := range testCasesInt {
		decimalInt := decimal.DecimalFromInt(i)
		if decimalInt.PowerOfTen != 0 {
			t.Fatalf("Failed to setup %d Decimal correctly: %v", i, decimalInt)
		}
		if decimalInt.Sign != expect.positive || decimalInt.Value != expect.value {
			t.Fatalf("[%d] Expected (%v, %d), instead got (%v, %d)",
				i,
				expect.positive, expect.value,
				decimalInt.Sign, decimalInt.Value)
		}
	}
	// Uint
	testCasesUint := []uint64{
		0, 1, math.MaxUint64,
	}
	for _, i := range testCasesUint {
		t.Logf("Testing positive uint64 %d", i)
		decimalPositive := decimal.DecimalFromUint(true, i)
		if !decimalPositive.Sign || decimalPositive.PowerOfTen != 0 {
			t.Fatalf("Failed to setup Decimal %d correctly: %v", i, decimalPositive)
		}
		if decimalPositive.Value != i {
			t.Fatalf("Expected %d, instead got %d", i, decimalPositive.Value)
		}
		t.Logf("Testing negative uint64 %d", i)
		decimalNegative := decimal.DecimalFromUint(false, i)
		if decimalNegative.Sign || decimalNegative.PowerOfTen != 0 {
			t.Fatalf("Failed to setup Decimal %d correctly: %v", i, decimalNegative)
		}
		if decimalNegative.Value != i {
			t.Fatalf("Expected %d, instead got %d", i, decimalNegative.Value)
		}
	}
}
