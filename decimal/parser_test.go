package decimal_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/stefanovazzocell/GoDecimal/decimal"
)

func TestParseString(t *testing.T) {
	testCases := map[string]struct {
		decimal  decimal.Decimal
		hasError bool
	}{
		"":                             {decimal.Decimal{Sign: true}, false},
		"0":                            {decimal.Decimal{Sign: true}, false},
		".e":                           {decimal.Decimal{Sign: true}, false},
		"0000000000000000000001":       {decimal.Decimal{Sign: true, Value: 1}, false},
		"1e0000000000000000000001":     {decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 1}, false},
		"123":                          {decimal.Decimal{Sign: true, Value: 123}, false},
		"1e10":                         {decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 10}, false},
		".1e10":                        {decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 9}, false},
		"1e+10":                        {decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 10}, false},
		"+1":                           {decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 0}, false},
		"-1":                           {decimal.Decimal{Sign: false, Value: 1, PowerOfTen: 0}, false},
		"-1!23.45e-23":                 {decimal.Decimal{Sign: false, Value: 12345, PowerOfTen: -25}, false},
		"23.e":                         {decimal.Decimal{Sign: true, Value: 23, PowerOfTen: 0}, false},
		"90000000000000000000":         {decimal.Decimal{Sign: true, Value: 9000000000000000000, PowerOfTen: 1}, false},
		"-123456789.0000123456789E123": {decimal.Decimal{Sign: false, Value: 12345678900001234567, PowerOfTen: 112}, false},
		"+987654321.0000987654321e123": {decimal.Decimal{Sign: true, Value: 9876543210000987654, PowerOfTen: 113}, false},
		"1e+" + strconv.FormatInt(math.MaxInt64, 10) + "0":                 {decimal.Decimal{}, true},
		"1e-" + strconv.FormatInt(math.MinInt64, 10) + "0":                 {decimal.Decimal{}, true},
		"1000000000000000000000e" + strconv.FormatInt(math.MaxInt64-1, 10): {decimal.Decimal{}, true},
		"0.01e" + strconv.FormatInt(math.MinInt64+1, 10):                   {decimal.Decimal{}, true},
	}
	for numberStr, testExpected := range testCases {
		// Prepare a readable string for error messages
		readable := numberStr
		if len(readable) > 20 {
			readable = readable[:20]
		}
		// Parse string
		actualNumber, err := decimal.ParseString(numberStr)
		// Check error
		if testExpected.hasError && err == nil {
			t.Errorf("Expected %q to produce an error, instead got nothing", readable)
		}
		if err != nil && !testExpected.hasError {
			t.Errorf("%q returned an unexpected error: %v", readable, err)
		}
		// Check number
		if actualNumber.Sign != testExpected.decimal.Sign ||
			actualNumber.Value != testExpected.decimal.Value ||
			actualNumber.PowerOfTen != testExpected.decimal.PowerOfTen {
			t.Errorf("%q expected:\n%v\ninstead got:\n%v", readable, testExpected.decimal, actualNumber)
		}
	}
}

func FuzzParseString(f *testing.F) {
	seeds := []string{
		"",
		"0",
		".e",
		"0000000000000000000001",
		"1e0000000000000000000001",
		"123",
		"1e10",
		".1e10",
		"1e+10",
		"+1",
		"-1",
		"-1!23.45e-23",
		"23.e",
		"90000000000000000000",
		"-123456789.0000123456789E123",
		"+987654321.0000987654321e123",
		"1e+" + strconv.FormatInt(math.MaxInt64, 10) + "0",
		"1e-" + strconv.FormatInt(math.MinInt64, 10) + "0",
		"1000000000000000000000e" + strconv.FormatInt(math.MaxInt64-1, 10),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, numberStr string) {
		decimal.ParseString(numberStr)
	})
}
