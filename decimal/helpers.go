package decimal

import (
	"math"
)

// Calculate the Greatest Common Divisor (GCD) between two uint64
func gcd(x, y uint64) uint64 {
	for y != 0 {
		t := y
		y = x % y
		x = t
	}
	return x
}

// Calculate the Least Common Multiple (LCM) between two uint64
// Returns (0, false) if an overflow is detected
func lcm(x, y uint64) (uint64, bool) {
	if x == 0 || y == 0 {
		return 0, false
	}
	t := x / gcd(x, y)
	if overflow_multiplication(t, y) {
		return 0, true
	}
	return t * y, false
}

// Returns true if x + y (int64) overflows
func overflow_int64(x, y int64) bool {
	// Overflow up
	if x > 0 && y > 0 && x > (math.MaxInt64-y) {
		return true
	}
	// Overflow down
	if x < 0 && y < 0 && x < (math.MinInt64-y) {
		return true
	}
	// Not an overflow
	return false
}

// Returns true if x + y will overflow
func overflow_sum(x, y uint64) bool {
	return math.MaxUint64-x < y
}

// Returns true if x - y will overflow
func overflow_subtraction(x, y uint64) bool {
	return x < y
}

// Returns true if x * y will overflow
func overflow_multiplication(x, y uint64) bool {
	if x == 0 {
		return false
	}
	return math.MaxUint64/x < y
}
