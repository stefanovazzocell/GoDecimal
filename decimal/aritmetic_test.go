package decimal_test

import (
	"math"
	"testing"

	"github.com/stefanovazzocell/GoDecimal/decimal"
)

var sharedTestCases = map[struct{ x, y decimal.Decimal }]struct {
	add, sub, mult, div             decimal.Decimal
	add_ok, sub_ok, mult_ok, div_ok bool
}{
	{decimal.Decimal{}, decimal.Decimal{}}: {
		add:    decimal.Decimal{},
		sub:    decimal.Decimal{},
		mult:   decimal.Decimal{},
		div:    decimal.Decimal{},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: false,
	},
	{decimal.Decimal{Value: 1}, decimal.Decimal{}}: {
		add:    decimal.Decimal{Value: 1},
		sub:    decimal.Decimal{Value: 1},
		mult:   decimal.Decimal{},
		div:    decimal.Decimal{Sign: true, Value: math.MaxUint64, PowerOfTen: math.MaxInt64},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: false,
	},
	{decimal.Decimal{}, decimal.Decimal{Sign: true, Value: 1}}: {
		add:    decimal.Decimal{Sign: true, Value: 1},
		sub:    decimal.Decimal{Value: 1},
		mult:   decimal.Decimal{},
		div:    decimal.Decimal{},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: true,
	},
	{decimal.Decimal{PowerOfTen: 1}, decimal.Decimal{}}: {
		add:    decimal.Decimal{},
		sub:    decimal.Decimal{},
		mult:   decimal.Decimal{},
		div:    decimal.Decimal{Sign: true}, // 0/0
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: false,
	},
	{decimal.Decimal{}, decimal.Decimal{Sign: true}}: {
		add:    decimal.Decimal{},
		sub:    decimal.Decimal{},
		mult:   decimal.Decimal{},
		div:    decimal.Decimal{},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: false,
	},
	{decimal.Decimal{Value: 100}, decimal.Decimal{Value: 1, PowerOfTen: 2}}: {
		add:    decimal.Decimal{Value: 200},
		sub:    decimal.Decimal{},
		mult:   decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 4},
		div:    decimal.Decimal{Sign: true, Value: 1},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: true,
	},
	{decimal.Decimal{Value: math.MaxUint64}, decimal.Decimal{Sign: true, Value: 1}}: {
		add:    decimal.Decimal{Value: math.MaxUint64 - 1},
		sub:    decimal.Decimal{Sign: false, Value: (math.MaxUint64/10 + 1), PowerOfTen: 1},
		mult:   decimal.Decimal{Value: math.MaxUint64},
		div:    decimal.Decimal{Value: math.MaxUint64},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: true,
	},
	{decimal.Decimal{Value: math.MaxUint64}, decimal.Decimal{Sign: true, Value: 10}}: {
		add:    decimal.Decimal{Value: math.MaxUint64 - 10},
		sub:    decimal.Decimal{Sign: false, Value: (math.MaxUint64 / 10) + 1, PowerOfTen: 1},
		mult:   decimal.Decimal{Value: math.MaxUint64, PowerOfTen: 1},
		div:    decimal.Decimal{Value: math.MaxUint64, PowerOfTen: -1},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: true,
	},
	{decimal.Decimal{Sign: true, Value: 2}, decimal.Decimal{Sign: true, Value: math.MaxUint64}}: {
		add:    decimal.Decimal{Sign: true, Value: (math.MaxUint64/10 + 1), PowerOfTen: 1},
		sub:    decimal.Decimal{Sign: false, Value: math.MaxUint64 - 2},
		mult:   decimal.Decimal{Sign: true, Value: (math.MaxUint64 / 10 * 2), PowerOfTen: 1},
		div:    decimal.Decimal{Sign: true, Value: 1, PowerOfTen: -19},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: true,
	},
	{decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 0}, decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 100}}: {
		add:    decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 100},
		sub:    decimal.Decimal{Sign: false, Value: 10000000000000000000, PowerOfTen: 81},
		mult:   decimal.Decimal{Sign: true, Value: 1, PowerOfTen: 100},
		div:    decimal.Decimal{Sign: true, Value: 1, PowerOfTen: -100},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: true,
	},
	{decimal.Decimal{Sign: true, Value: 12, PowerOfTen: 0}, decimal.Decimal{Sign: true, Value: 12, PowerOfTen: 100}}: {
		add:    decimal.Decimal{Sign: true, Value: 12, PowerOfTen: 100},
		sub:    decimal.Decimal{Sign: false, Value: 12, PowerOfTen: 100},
		mult:   decimal.Decimal{Sign: true, Value: 144, PowerOfTen: 100},
		div:    decimal.Decimal{Sign: true, Value: 1, PowerOfTen: -100},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: true,
	},
	{decimal.Decimal{Sign: true, Value: math.MaxUint64 - 1, PowerOfTen: math.MaxInt64}, decimal.Decimal{Sign: true, Value: 2, PowerOfTen: math.MaxInt64}}: {
		add:    decimal.Decimal{Sign: true, Value: math.MaxUint64, PowerOfTen: math.MaxInt64},
		sub:    decimal.Decimal{Sign: true, Value: math.MaxUint64 - 3, PowerOfTen: math.MaxInt64},
		mult:   decimal.Decimal{Sign: true, Value: math.MaxUint64, PowerOfTen: math.MaxInt64},
		div:    decimal.Decimal{Sign: true, Value: (math.MaxUint64 - 1) / 2, PowerOfTen: 0},
		add_ok: false, sub_ok: true, mult_ok: false, div_ok: true,
	},
	{decimal.Decimal{Value: 1234567891234, PowerOfTen: math.MaxInt64}, decimal.Decimal{Value: 123456789, PowerOfTen: -1}}: {
		add:    decimal.Decimal{Value: 1234567891234, PowerOfTen: math.MaxInt64},
		sub:    decimal.Decimal{Value: 1234567891234, PowerOfTen: math.MaxInt64},
		mult:   decimal.Decimal{Sign: true, Value: (1234567891234 * 12345678), PowerOfTen: math.MaxInt64},
		div:    decimal.Decimal{Sign: true, Value: 100000000099, PowerOfTen: 9223372036854775801},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: true,
	},
	{decimal.Decimal{Value: 1, PowerOfTen: math.MaxInt64 - 1}, decimal.Decimal{Value: 1, PowerOfTen: 2}}: {
		add:    decimal.Decimal{Value: 1, PowerOfTen: math.MaxInt64 - 1},
		sub:    decimal.Decimal{Value: 1, PowerOfTen: math.MaxInt64 - 1},
		mult:   decimal.Decimal{Sign: true, Value: math.MaxUint64, PowerOfTen: math.MaxInt64},
		div:    decimal.Decimal{Sign: true, Value: 1, PowerOfTen: math.MaxInt64 - 3},
		add_ok: true, sub_ok: true, mult_ok: false, div_ok: true,
	},
	{decimal.Decimal{Value: 1, PowerOfTen: math.MinInt64 + 1}, decimal.Decimal{Value: 1, PowerOfTen: -2}}: {
		add:    decimal.Decimal{Value: 1, PowerOfTen: -2},
		sub:    decimal.Decimal{Sign: true, Value: 1, PowerOfTen: -2},
		mult:   decimal.Decimal{Sign: true, Value: 0, PowerOfTen: 0},
		div:    decimal.Decimal{Sign: true, Value: 1, PowerOfTen: math.MinInt64 + 3},
		add_ok: true, sub_ok: true, mult_ok: false, div_ok: true,
	},
	{decimal.Decimal{Value: math.MaxUint64, PowerOfTen: math.MinInt64 + 1}, decimal.Decimal{Value: math.MaxUint64, PowerOfTen: 2}}: {
		add:    decimal.Decimal{Value: math.MaxUint64, PowerOfTen: 2},
		sub:    decimal.Decimal{Sign: true, Value: math.MaxUint64, PowerOfTen: 2},
		mult:   decimal.Decimal{Sign: true, Value: 3402823667840801649, PowerOfTen: math.MinInt64 + 23},
		div:    decimal.Decimal{Sign: true, Value: 10, PowerOfTen: math.MinInt64},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: true,
	},
	{decimal.Decimal{Value: math.MaxUint64, PowerOfTen: math.MinInt64 + 1}, decimal.Decimal{Value: 2, PowerOfTen: 2}}: {
		add:    decimal.Decimal{Value: 2, PowerOfTen: 2},
		sub:    decimal.Decimal{Sign: true, Value: 2, PowerOfTen: 2},
		mult:   decimal.Decimal{Sign: true, Value: 3689348814741910322, PowerOfTen: math.MinInt64 + 4},
		div:    decimal.Decimal{Sign: true, Value: 0, PowerOfTen: math.MinInt64},
		add_ok: true, sub_ok: true, mult_ok: true, div_ok: false,
	},
}

func TestAritmetic(t *testing.T) {
	// Make sure we handle nil
	var nilDecimal *decimal.Decimal = nil
	if nilDecimal.Add(decimal.Decimal{}, decimal.Decimal{}) ||
		nilDecimal.Sub(decimal.Decimal{}, decimal.Decimal{}) ||
		nilDecimal.Mult(decimal.Decimal{}, decimal.Decimal{}) ||
		nilDecimal.Div(decimal.Decimal{}, decimal.Decimal{}) {
		t.Fatal("Expect aritmetic ops to return false if run on a nil decimal")
	}
	// run test
	for test, expected := range sharedTestCases {
		add := decimal.Decimal{}
		add_ok := add.Add(test.x, test.y)
		sub := decimal.Decimal{}
		sub_ok := sub.Sub(test.x, test.y)
		mult := decimal.Decimal{}
		mult_ok := mult.Mult(test.x, test.y)
		div := decimal.Decimal{}
		div_ok := div.Div(test.x, test.y)
		// Compare results
		if !add.Equals(expected.add) || add_ok != expected.add_ok {
			t.Fatalf("add with %v and %v expected (%v, %v) but got (%v, %v)",
				test.x, test.y,
				expected.add, expected.add_ok,
				add, add_ok,
			)
		}
		if !sub.Equals(expected.sub) || sub_ok != expected.sub_ok {
			t.Fatalf("sub with %v and %v expected (%v, %v) but got (%v, %v)",
				test.x, test.y,
				expected.sub, expected.sub_ok,
				sub, sub_ok,
			)
		}
		if !mult.Equals(expected.mult) || mult_ok != expected.mult_ok {
			t.Fatalf("mult with %v and %v expected (%v, %v) but got (%v, %v)",
				test.x, test.y,
				expected.mult, expected.mult_ok,
				mult, mult_ok,
			)
		}
		if !div.Equals(expected.div) || div_ok != expected.div_ok {
			t.Fatalf("div with %v and %v expected (%v, %v) but got (%v, %v)",
				test.x, test.y,
				expected.div, expected.div_ok,
				div, div_ok,
			)
		}
	}
}

func BenchmarkAritmetic(b *testing.B) {
	dec := decimal.Decimal{}
	cases := []struct{ x, y decimal.Decimal }{}
	for test := range sharedTestCases {
		cases = append(cases, test)
	}
	casesLen := len(cases)
	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec.Add(cases[i%casesLen].x, cases[i%casesLen].y)
		}
	})
	b.Run("Sub", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec.Sub(cases[i%casesLen].x, cases[i%casesLen].y)
		}
	})
	b.Run("Mult", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec.Mult(cases[i%casesLen].x, cases[i%casesLen].y)
		}
	})
	b.Run("Div", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec.Div(cases[i%casesLen].x, cases[i%casesLen].y)
		}
	})
}

// NOTE: Since subtraction uses the same base function as addition, this effectively tests both
func FuzzAdd(f *testing.F) {
	// Setup some helpers

	// Returns true if x + y will overflow
	overflow_sum := func(x, y uint64) bool {
		return math.MaxUint64-x < y
	}

	// Returns true if x - y will overflow
	overflow_subtraction := func(x, y uint64) bool {
		return x < y
	}

	// Add to seed corpus
	testCases := []struct {
		first_sign  bool
		first       uint64
		second_sign bool
		second      uint64
	}{
		{true, 0, true, 0},
		{false, 0, true, 0},
		{true, 0, false, 0},
		{false, 0, false, 0},
		{true, 1, true, 1},
		{false, 1, true, 1},
		{true, 1, false, 1},
		{false, 1, false, 1},
		{true, 123, true, 234},
		{false, 123, true, 234},
		{true, 123, false, 234},
		{false, 123, false, 234},
	}
	for _, seed := range testCases {
		f.Add(seed.first_sign, seed.first, seed.second_sign, seed.second)
	}

	// Fuzz
	f.Fuzz(func(t *testing.T, first_sign bool, first uint64, second_sign bool, second uint64) {
		// Prepare actual value
		first_decimal := decimal.Decimal{
			Sign:       first_sign,
			Value:      first,
			PowerOfTen: 0,
		}
		second_decimal := decimal.Decimal{
			Sign:       second_sign,
			Value:      second,
			PowerOfTen: 0,
		}
		actual := decimal.Decimal{}
		actual.Add(first_decimal, second_decimal)
		// Prepare expected value
		expected_sign := first_sign
		expected := first + second
		if first_sign != second_sign {
			// Yuck, nested IFs!
			expected_sign = true
			if first_sign {
				// first - second
				expected = first - second
				if overflow_subtraction(first, second) {
					expected_sign = false
					expected = second - first
				}
			} else {
				// second - first
				expected = second - first
				if overflow_subtraction(second, first) {
					expected_sign = false
					expected = first - second
				}
			}
		} else if overflow_sum(first, second) {
			// The logic in this fuzz is already pretty complex, so we'll actually skip this.
			// What we're testing here is that for trivial operations we can produce valid results
			// and for all operations we don't panic.
			t.SkipNow()
		}
		expected_decimal := decimal.Decimal{
			Sign:       expected_sign,
			Value:      expected,
			PowerOfTen: 0,
		}
		// Compare actual to expected
		if !expected_decimal.Equals(actual) {
			t.Fatalf("For %v + %v\ngot %v,\nbut expected %v", first, second, actual, expected)
		}
	})
}
