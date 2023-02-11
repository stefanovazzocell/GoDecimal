package decimal_test

import (
	"testing"

	"github.com/stefanovazzocell/GoDecimal/decimal"
)

func TestUtils(t *testing.T) {
	// A decimal must be one of these types
	type numberType uint8
	var (
		nZero     numberType = 0
		nPositive numberType = 1
		nNegative numberType = 2
	)
	testCases := map[decimal.Decimal]numberType{
		{Sign: true, Value: 0}:                  nZero,
		{Sign: false, Value: 0}:                 nZero,
		{Sign: true, Value: 143}:                nPositive,
		{Sign: false, Value: 1, PowerOfTen: 43}: nNegative,
	}
	for testNumber, expected := range testCases {
		t.Logf("Testing decimal %v", testNumber)
		// Cloning
		cloned := testNumber.Clone()
		if testNumber.Sign != cloned.Sign || testNumber.Value != cloned.Value ||
			testNumber.PowerOfTen != cloned.PowerOfTen {
			t.Errorf("%v Clone() incorrectly returned %v", testNumber, cloned)
		}
		// Other utils
		if testNumber.IsZero() != (expected == nZero) {
			t.Errorf("IsZero() expected %v, got %v", (expected == nZero), testNumber.IsZero())
		}
		if testNumber.IsPositive() != (expected == nPositive) {
			t.Errorf("IsPositive() expected %v, got %v", (expected == nPositive), testNumber.IsPositive())
		}
		if testNumber.IsNegative() != (expected == nNegative) {
			t.Errorf("IsNegative() expected %v, got %v", (expected == nNegative), testNumber.IsNegative())
		}
		// If negative test Abs()
		testNumber.Abs()
		if testNumber.IsNegative() {
			t.Errorf("Failed to Abs() decimal, still registering as negative: got %v", testNumber)
		}
		if testNumber.IsZero() != (expected == nZero) {
			t.Errorf("After Abs(), IsZero() expected %v, got %v", (expected == nZero), testNumber.IsZero())
		}
		// Test Zero()
		testNumber.Zero()
		if !testNumber.IsZero() || testNumber.IsPositive() || testNumber.IsNegative() {
			t.Errorf("Failed to Zero() decimal: got %v", testNumber)
		}
	}
	// Can Abs() and Zero() handle nil without panic?
	var nilNumb *decimal.Decimal = nil
	nilNumb.Abs()
	nilNumb.Zero()
}

func TestExpandCompressReduce(t *testing.T) {
	var nilDecimal *decimal.Decimal = nil
	nilDecimal.Compress()
	if nilDecimal != nil {
		t.Fatal("Incorrect handling of nil [1/2]")
	}
	nilDecimal.Expand()
	if nilDecimal != nil {
		t.Fatal("Incorrect handling of nil [2/2]")
	}
	testCases := map[decimal.Decimal]struct {
		compressed decimal.Decimal
		expanded   decimal.Decimal
	}{ // Sign, Value, Denumerator, PowerOfTen
		{false, 0, 2}: {
			compressed: decimal.Decimal{false, 0, 0},
			expanded:   decimal.Decimal{false, 0, 0},
		},
		{true, 1, 5}: {
			compressed: decimal.Decimal{true, 1, 5},
			expanded:   decimal.Decimal{true, 10000000000000000000, -14},
		},
		{false, 90, -4}: {
			compressed: decimal.Decimal{false, 9, -3},
			expanded:   decimal.Decimal{false, 9000000000000000000, -21},
		},
		{true, 5000, 2}: {
			compressed: decimal.Decimal{true, 5, 5},
			expanded:   decimal.Decimal{true, 5000000000000000000, -13},
		},
		{false, 1480604557035703029, -3}: {
			compressed: decimal.Decimal{false, 1480604557035703029, -3},
			expanded:   decimal.Decimal{false, 14806045570357030290, -4},
		},
	}
	for test, expected := range testCases {
		compressed := test.Clone()
		compressed.Compress()
		expanded := test.Clone()
		expanded.Expand()
		if expected.compressed.Sign != compressed.Sign ||
			expected.compressed.Value != compressed.Value ||
			expected.compressed.PowerOfTen != compressed.PowerOfTen {
			t.Errorf("%v Compress() incorrectly returned %v", test, compressed)
		}
		if expected.expanded.Sign != expanded.Sign ||
			expected.expanded.Value != expanded.Value ||
			expected.expanded.PowerOfTen != expanded.PowerOfTen {
			t.Errorf("%v Expand() incorrectly returned %v", test, expanded)
		}
	}
}

func FuzzUtils(f *testing.F) {
	seeds := []decimal.Decimal{
		{},
		{true, 10, -2},
		{true, 654245, 234},
		{false, 2304042349828359000, 6},
	}
	for _, seed := range seeds {
		f.Add(seed.Sign, seed.Value, seed.PowerOfTen)
	}
	f.Fuzz(func(t *testing.T, sign bool, value uint64, powerOfTen int64) {
		n := decimal.Decimal{
			Sign:       sign,
			Value:      value,
			PowerOfTen: powerOfTen,
		}
		if n.IsZero() && (n.IsPositive() || n.IsNegative()) {
			t.Fatalf("If the decimal is zero, it can't be positive or negative.")
		}
		if n.IsPositive() && (n.IsNegative() || n.IsZero()) {
			t.Fatalf("If the decimal is finite and negative, it can't be positive or zero")
		}
		if n.IsNegative() && (n.IsPositive() || n.IsZero()) {
			t.Fatalf("If the decimal is finite and positive, it can't be negative or zero")
		}
		// Try to set absolute
		n.Abs()
		if n.IsNegative() {
			t.Fatalf("Failed to Abs() decimal, still registering as negative: got %v", n)
		}
		// Try to zero
		n.Zero()
		if !n.IsZero() || n.IsPositive() || n.IsNegative() {
			t.Errorf("Failed to Zero() decimal: got %v", n)
		}
	})
}
