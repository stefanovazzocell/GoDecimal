package decimal

import "math"

// Returns true if the number is zero
func (d Decimal) IsZero() bool {
	return d.Value == 0
}

// Returns true if the number is positive
func (d Decimal) IsPositive() bool {
	return d.Sign && !d.IsZero()
}

// Returns true if the number is negative
func (d Decimal) IsNegative() bool {
	return !d.Sign && !d.IsZero()
}

// Returns a clone of a number
func (d Decimal) Clone() Decimal {
	return Decimal{
		Sign:       d.Sign,
		Value:      d.Value,
		PowerOfTen: d.PowerOfTen,
	}
}

// Makes the number value absolute
func (d *Decimal) Abs() {
	if d == nil {
		return
	}
	d.Sign = true
}

// Zeroes the number
func (d *Decimal) Zero() {
	if d == nil {
		return
	}
	d.Sign = true
	d.Value = 0
	d.PowerOfTen = 0
}

// Increase the Value of the number compensating by adjusting it's power of ten.
// Will not affect accuracy or precision.
func (d *Decimal) Expand() {
	if d == nil {
		return
	} else if d.Value == 0 {
		d.PowerOfTen = 0
	} else {
		for !overflow_multiplication(d.Value, 10) && math.MinInt64 < d.PowerOfTen {
			d.PowerOfTen--
			d.Value *= 10
		}
	}
}

// Compress the Value of the number compensating by adjusting it's power of ten.
// Will not affect accuracy or precision.
func (d *Decimal) Compress() {
	if d == nil {
		return
	} else if d.Value == 0 {
		d.PowerOfTen = 0
	} else {
		for d.Value%10 == 0 && d.PowerOfTen < math.MaxInt64 {
			d.PowerOfTen++
			d.Value /= 10
		}
	}
}
