package decimal

// A representation of a decimal number in scientific notation
type Decimal struct {
	Sign       bool
	Value      uint64
	PowerOfTen int64
}

// Returns a decimal based on a given int64
func DecimalFromInt(i int64) Decimal {
	positive := i >= 0
	if !positive {
		i *= -1
	}
	return Decimal{
		Sign:       positive,
		Value:      uint64(i),
		PowerOfTen: 0,
	}
}

// Returns a decimal based on a given uint64
func DecimalFromUint(sign bool, u uint64) Decimal {
	return Decimal{
		Sign:       sign,
		Value:      u,
		PowerOfTen: 0,
	}
}
