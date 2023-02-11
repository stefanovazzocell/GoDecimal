package decimal_test

import (
	"testing"

	"github.com/stefanovazzocell/GoDecimal/decimal"
)

func TestLogic(t *testing.T) {
	testCases := map[struct{ x, y decimal.Decimal }]struct {
		equals bool
	}{
		{decimal.Decimal{}, decimal.Decimal{}}:                                  {equals: true},
		{decimal.Decimal{Value: 1}, decimal.Decimal{}}:                          {equals: false},
		{decimal.Decimal{PowerOfTen: 1}, decimal.Decimal{}}:                     {equals: true},
		{decimal.Decimal{}, decimal.Decimal{Sign: true}}:                        {equals: true},
		{decimal.Decimal{Value: 100}, decimal.Decimal{Value: 1, PowerOfTen: 2}}: {equals: true},
	}

	for test, expected := range testCases {
		equals := test.x.Equals(test.y)
		if equals != expected.equals {
			t.Errorf("%v.Equals(%v) reported %v, but expected %v", test.x, test.y, equals, expected.equals)
		}
	}
}
