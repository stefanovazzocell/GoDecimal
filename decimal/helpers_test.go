package decimal

import (
	"math"
	"testing"
)

func TestHelpers(t *testing.T) {
	testCases := map[struct{ x, y uint64 }]struct {
		gcd                     uint64
		lcm                     uint64
		lcm_overflow            bool
		sum_overflow            bool
		subtraction_overflow    bool
		multiplication_overflow bool
	}{
		{0, 0}:                           {0, 0, false, false, false, false},
		{math.MaxUint64, 0}:              {math.MaxUint64, 0, false, false, false, false},
		{0, math.MaxUint64}:              {math.MaxUint64, 0, false, false, true, false},
		{1, math.MaxUint64}:              {1, math.MaxUint64, false, true, true, false},
		{math.MaxUint64, 1}:              {1, math.MaxUint64, false, true, false, false},
		{13, math.MaxUint64}:             {1, 0, true, true, true, true},
		{math.MaxUint64, 13}:             {1, 0, true, true, false, true},
		{math.MaxUint64, math.MaxUint64}: {math.MaxUint64, math.MaxUint64, false, true, false, true},
		{543534, 134325}:                 {3, 24336734850, false, false, false, false},
		{134325, 543534}:                 {3, 24336734850, false, false, true, false},
	}

	for test, expect := range testCases {
		t.Logf("Testing with %d, %d", test.x, test.y)
		if actual := gcd(test.x, test.y); actual != expect.gcd {
			t.Fatalf("Expected %d from gcm, instead got %d", expect.gcd, actual)
		}
		if actual, overflow := lcm(test.x, test.y); actual != expect.lcm || overflow != expect.lcm_overflow {
			t.Fatalf("Expected (%d, %v) from lcm, instead got (%d, %v)", expect.lcm, expect.lcm_overflow, actual, overflow)
		}
		// Overflows
		if actual := overflow_sum(test.x, test.y); actual != expect.sum_overflow {
			t.Fatalf("Expected %v from overflow_sum, instead got %v", expect.sum_overflow, actual)
		}
		if actual := overflow_subtraction(test.x, test.y); actual != expect.subtraction_overflow {
			t.Fatalf("Expected %v from overflow_subtraction, instead got %v", expect.subtraction_overflow, actual)
		}
		if actual := overflow_multiplication(test.x, test.y); actual != expect.multiplication_overflow {
			t.Fatalf("Expected %v from overflow_multiplication, instead got %v", expect.multiplication_overflow, actual)
		}
	}

	testCasesInt64Overflow := map[struct{ x, y int64 }]bool{
		{0, 0}:                         false,
		{1, 2}:                         false,
		{2, -1}:                        false,
		{math.MaxInt64, math.MinInt64}: false,
		{math.MaxInt64, math.MaxInt64}: true,
		{math.MinInt64 + 20, -21}:      true,
		{math.MinInt64 + 20, -20}:      false,
	}
	for test, expected := range testCasesInt64Overflow {
		if actual := overflow_int64(test.x, test.y); actual != expected {
			t.Fatalf("Expected %v from overflow_int64(%d, %d), but got %v", expected, test.x, test.y, actual)
		}
	}
}

func BenchmarkHelpers(b *testing.B) {
	b.Run("gcd", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = gcd(869505786, 640448604)
		}
	})
	b.Run("lcm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = lcm(869505786, 640448604)
		}
	})
	b.Run("overflow_int64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = overflow_int64(869505786, -640448604)
		}
	})
	b.Run("overflow_sum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = overflow_sum(869505786, 640448604)
		}
	})
	b.Run("overflow_subtraction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = overflow_subtraction(869505786, 640448604)
		}
	})
	b.Run("overflow_multiplication", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = overflow_multiplication(869505786, 640448604)
		}
	})
}

func FuzzHelpers(f *testing.F) {
	f.Fuzz(func(t *testing.T, x uint64, y uint64) {
		z := gcd(x, y)
		if z != 0 && (x/z)*z != x {
			t.Fatalf("GDC %d seems wrong", z)
		}
		// Overflows
		sum_overflow := (x+y)-x != y
		if actual := overflow_sum(x, y); actual != sum_overflow {
			t.Fatalf("Expected overflow_sum %v, instead got %v", sum_overflow, actual)
		}
		subtraction_overflow := x < y
		if actual := overflow_subtraction(x, y); actual != subtraction_overflow {
			t.Fatalf("Expected overflow_subtraction %v, instead got %v", subtraction_overflow, actual)
		}
		multiplication_overflow := (x > 0 && y > 0 && x*y < x && x*y < y)
		if actual := overflow_multiplication(x, y); actual != multiplication_overflow {
			t.Fatalf("Expected overflow_multiplication %v, instead got %v", multiplication_overflow, actual)
		}
	})
}
