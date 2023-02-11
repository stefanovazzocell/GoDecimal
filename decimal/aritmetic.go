package decimal

import (
	"math"
)

// Perform the addition x + y and store the result in this decimal.
// If the decimal is nil, this operation will be a noop.
// If the operation overflows/underflows, will return false
func (d *Decimal) Add(x, y Decimal) (ok bool) {
	// Special case: nil
	if d == nil {
		return false // NOOP
	}
	// Check if this is a subtraction
	sameSign := true
	if x.Sign != y.Sign {
		sameSign = false
	}
	// Reset n
	d.Sign = true
	d.Value = 0
	d.PowerOfTen = 0
	// Get x, y on the same power of 10
	x.Compress()
	y.Compress()
	if x.PowerOfTen != y.PowerOfTen {
		powerHigh := &x
		powerLow := &y
		if x.PowerOfTen < y.PowerOfTen {
			powerHigh = &y
			powerLow = &x
		}
		// Increase the value of the higher to decrease its power of 10
		for x.PowerOfTen != y.PowerOfTen && !overflow_multiplication(powerHigh.Value, 10) {
			powerHigh.Value *= 10
			powerHigh.PowerOfTen--
		}
		// [Loss of precision!] Decrease the value of the lower to increase its power of 10
		for x.PowerOfTen != y.PowerOfTen && powerLow.Value >= 10 {
			powerLow.Value /= 10
			powerLow.PowerOfTen++
		}
		// Check if lower went to zero
		if x.PowerOfTen != y.PowerOfTen {
			// The lower number is so small that it's insignificant
			d.Sign = powerHigh.Sign
			d.Value = powerHigh.Value
			d.PowerOfTen = powerHigh.PowerOfTen
			return true
		}
	}
	// Perform operation, being mindful of overflow
	if sameSign {
		if overflow_sum(x.Value, y.Value) {
			if x.PowerOfTen == math.MaxInt64 || y.PowerOfTen == math.MaxInt64 {
				// Overflow
				d.Sign = x.Sign
				d.Value = math.MaxUint64
				d.PowerOfTen = math.MaxInt64
				return false
			}
			// If we'll overflow, knock the value down, if we can
			if x.Value < 10 {
				d.Sign = y.Sign
				d.Value = y.Value/10 + 1 // Add back remainder
				d.PowerOfTen = y.PowerOfTen + 1
				return true
			}
			if y.Value < 10 {
				d.Sign = x.Sign
				d.Value = x.Value/10 + 1
				d.PowerOfTen = x.PowerOfTen + 1
				return true
			}
			x.PowerOfTen++
			y.PowerOfTen++
			x.Value /= 10
			y.Value /= 10
		}
		d.Sign = x.Sign
		d.Value = x.Value + y.Value
		d.PowerOfTen = x.PowerOfTen
	} else if !x.Sign {
		// y + -x
		if overflow_subtraction(y.Value, x.Value) {
			// Invert
			d.Sign = x.Sign
			d.Value = x.Value - y.Value
			d.PowerOfTen = y.PowerOfTen
			return true
		}
		d.Sign = y.Sign
		d.Value = y.Value - x.Value
		d.PowerOfTen = y.PowerOfTen
	} else {
		// x + -y
		if overflow_subtraction(x.Value, y.Value) {
			// Invert
			d.Sign = y.Sign
			d.Value = y.Value - x.Value
			d.PowerOfTen = x.PowerOfTen
			return true
		}
		d.Sign = x.Sign
		d.Value = x.Value - y.Value
		d.PowerOfTen = x.PowerOfTen
	}
	return true
}

// Perform the subtraction x - y and store the result in this decimal.
// If the decimal is nil, this operation will be a noop.
// If the operation overflows/underflows, will return false
func (d *Decimal) Sub(x, y Decimal) (ok bool) {
	y.Sign = !y.Sign
	return d.Add(x, y)
}

// Perform the multiplication x * y and store the result in this decimal.
// If the decimal is nil, this operation will be a noop.
// If the operation overflows/underflows, will return false
func (d *Decimal) Mult(x, y Decimal) (ok bool) {
	// Special case: nil
	if d == nil {
		return false // NOOP
	}
	// Reset n
	d.Sign = x.Sign == y.Sign
	d.Value = 0
	d.PowerOfTen = 0
	// Compress x, y, then adjust the two values (losing precision) until multiplication doesn't overflow
	x.Compress()
	y.Compress()
	for overflow_multiplication(x.Value, y.Value) {
		larger := &x
		smaller := &y
		if x.Value < y.Value {
			larger = &y
			smaller = &x
		}
		// Let's try to decrease the larger first (if we can)
		if larger.PowerOfTen < math.MaxInt64 {
			// NOTE: The larger between x and y must be much larger than 10
			larger.Value /= 10
			larger.PowerOfTen++
			continue // Done for this iteration
		}
		// Okay, then let's try to decrease the smaller value (if we can)
		if smaller.Value >= 10 && smaller.PowerOfTen < math.MaxInt64 {
			smaller.Value /= 10
			smaller.PowerOfTen++
			continue // Done for this iteration
		}
		// If we got here we can't further reduce x and y to process the multiplication (overflow)
		d.Value = math.MaxUint64
		d.PowerOfTen = math.MaxInt64
		return false // Overflow
	}
	// Check power-of-tens overflow
	if overflow_int64(x.PowerOfTen, y.PowerOfTen) {
		if x.PowerOfTen < 0 {
			// Overflow to -inf
			d.Value = 0
			d.PowerOfTen = math.MinInt64
		} else {
			// Positive to +inf
			d.Value = math.MaxUint64
			d.PowerOfTen = math.MaxInt64
		}
		return false
	}
	// Perform multiplication
	d.Value = x.Value * y.Value
	d.PowerOfTen = x.PowerOfTen + y.PowerOfTen
	return true
}

// Perform the division x / y and store the result in this decimal.
// If the decimal is nil, this operation will be a noop.
// If the operation overflows/underflows, will return false
// If a divide-by-zero error is encountered, will return false
func (d *Decimal) Div(x, y Decimal) (ok bool) {
	// Special case: nil
	if d == nil {
		return false // NOOP
	}
	// Reset n
	d.Sign = x.Sign == y.Sign
	d.Value = 0
	d.PowerOfTen = 0
	// Special case: divide by zero
	if y.IsZero() {
		// 0 / 0 = ?
		if !x.IsZero() {
			// 1 / 0 = infinity
			d.Sign = true
			d.Value = math.MaxUint64
			d.PowerOfTen = math.MaxInt64
		}
		return false
	}
	// Special case: divide zero
	if x.IsZero() {
		// 0 / 1 = 0
		return true
	}
	// Prepare for division
	// We want to do some tweaks here in order to be as precise result possible
	x.Expand()   // Expand nominator
	y.Compress() // Compress denominator
	// We'd want x to be at least as big as y if at all possible
	for x.Value < y.Value && math.MinInt64 < y.PowerOfTen {
		y.PowerOfTen++
		y.Value /= 10
	}
	// Check boundaries
	for overflow_int64(x.PowerOfTen, -y.PowerOfTen) {
		// Check if we can attempt to resolve this
		if math.MinInt64 == y.PowerOfTen || y.Value < 10 {
			// The result will be too small to process
			d.PowerOfTen = math.MinInt64
			return false
		}
		// Attempt to resolve this
		y.Value /= 10
		y.PowerOfTen--
	}
	// Perform operation
	d.PowerOfTen = x.PowerOfTen - y.PowerOfTen
	d.Value = x.Value / y.Value
	return true
}
