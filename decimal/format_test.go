package decimal_test

import (
	"testing"

	"github.com/stefanovazzocell/GoDecimal/decimal"
)

func TestFormat(t *testing.T) {
	testCases := map[struct {
		number        string
		asDecimal     bool
		asPercentage  bool
		accuracyLimit int
	}]string{
		{"", true, true, 0}:            "0e2%",
		{"", true, false, 0}:           "0",
		{"", false, false, 0}:          "0",
		{"", true, true, 1}:            "0e2%",
		{"", true, false, 1}:           "0",
		{"", false, false, 1}:          "0",
		{"", true, true, 2}:            "0e2%",
		{"", true, false, 2}:           "0",
		{"", false, false, 2}:          "0",
		{"0.123", true, true, 2}:       "12%",
		{"0.123", true, false, 2}:      "12e-2",
		{"0.123", false, false, 2}:     "0.12",
		{"0.123", false, false, 0}:     "0.123",
		{"1.23", false, false, 2}:      "1.2",
		{"-12300000", false, false, 2}: "-12000000",
		{"0.000123", false, false, 0}:  "0.000123",
		{"0.0123%", false, true, 2}:    "0.012%",
		{"1234.5", false, false, 2}:    "1200",
	}

	for test, expected := range testCases {
		number, err := decimal.ParseString(test.number)
		if err != nil {
			t.Fatalf("Failed to setup test: %v", err)
		}
		actual := number.Format(test.asDecimal, test.asPercentage, test.accuracyLimit)
		if actual != expected {
			t.Fatalf("Expected %q for %v (%v), instead got %q", expected, test, number, actual)
		}
	}
}
