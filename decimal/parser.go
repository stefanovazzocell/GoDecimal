package decimal

import (
	"errors"
	"math"
)

const (
	// How many digit can we read before we (almost) fill a uint64
	digitsCutoff = uint64(19)
)

var (
	ErrorParsingOverflow = errors.New("the given string has a number that is too large to parse correctly")
	ErrorNilPointer      = errors.New("use of a nil pointer as Decimal in call to UnmarshalText")
)

// Parse a decimal number from a given string, ignoring any unknown characters.
// `numberStr` is the string you want to parse.
//
// Notes:
//  1. If any part of the whole number is too large to parse will first lose precision, then overflow.
//  2. Will ignore any unrecognized character.
//  3. Parses anything after the first 'e' as the exponential.
//  4. If the number ends with '%' the number will be parsed as a percentage.
//
// Examples:
//  1. "-1!23.45e-23" parses as -12345 * 10 ^ -25
//  2. "12%" parses as 12 * 10 ^ -2
//  3. "-1!23.45e-23%" parses as -12345 * 10 ^ -27
//  4. "23.e" parses as 23
//  5. A string composed of 100 digits of "9" parses as 9999999999999999999 * 10 ^ (100 - 19)
//  6. "9." followed by a string of 100 digits of "9" parses as 9999999999999999999 * 10 ^ -(100  - 18)
func ParseString(numberStr string) (Decimal, error) {
	// Setup defaults
	decimal := Decimal{
		Sign:       true,
		Value:      0,
		PowerOfTen: 0,
	}
	// Special case: 0
	if len(numberStr) == 0 {
		return decimal, nil
	}
	// Check sign
	if numberStr[0] == '+' {
		numberStr = numberStr[1:]
	} else if numberStr[0] == '-' {
		decimal.Sign = false
		numberStr = numberStr[1:]
	}
	// Check if percentage
	if len(numberStr) > 1 && numberStr[len(numberStr)-1] == '%' {
		// 1% == 0.01 == 1e-2
		decimal.PowerOfTen = -2
		numberStr = numberStr[:len(numberStr)-1]
	}
	// Value (whole + decimal part)
	digits := uint64(0)
	for i := 0; i < len(numberStr); i++ {
		if '0' <= numberStr[i] && numberStr[i] <= '9' {
			if digits == 0 {
				// It's the first digit, initialize the number
				if numberStr[i] != '0' {
					digits++
				}
				decimal.Value = uint64(numberStr[i] - '0')
			} else if digits > digitsCutoff {
				// We can't add any more number to the value (max precision)
				/*
					Note: number.PowerOfTen can never overflow.
					len(numberStr) is an integer and therefore cannot store enough digits to overflow
					number.PowerOfTen no matter how many digits are "soft-overflown" there.
				*/
				decimal.PowerOfTen++
			} else if digits == digitsCutoff {
				// Almost max precision, add digit if possible
				digits++
				if overflow_multiplication(decimal.Value, 10) {
					// Precision-loss overflow
					decimal.PowerOfTen++
					continue
				}
				decimal.Value = (decimal.Value * 10) + uint64(numberStr[i]-'0')
			} else {
				// Can add digit to value
				digits++
				decimal.Value = (decimal.Value * 10) + uint64(numberStr[i]-'0')
			}
		} else if numberStr[i] == '.' {
			// Moving to decimal part
			for j := i + 1; j < len(numberStr); j++ {
				if '0' <= numberStr[j] && numberStr[j] <= '9' {
					if digits == 0 {
						// It's the first digit, initialize the number
						digits++
						decimal.PowerOfTen--
						decimal.Value = uint64(numberStr[j] - '0')
					} else if digits > digitsCutoff {
						// We can't add any more number to the value (max precision)
						continue
					} else if digits == digitsCutoff {
						// Almost max precision, add digit if possible
						digits++
						if overflow_multiplication(decimal.Value, 10) {
							// Precision-loss overflow
							continue
						}
						decimal.PowerOfTen--
						decimal.Value = (decimal.Value * 10) + uint64(numberStr[j]-'0')
					} else {
						digits++
						decimal.PowerOfTen--
						decimal.Value = (decimal.Value * 10) + uint64(numberStr[j]-'0')
					}
				} else if numberStr[j] == '/' || numberStr[j] == 'e' || numberStr[j] == 'E' {
					numberStr = numberStr[j:]
					break
				}
			}
			break
		} else if numberStr[i] == '/' || numberStr[i] == 'e' || numberStr[i] == 'E' {
			numberStr = numberStr[i:]
			break
		}
	}
	// Power of 10 ('E'/'e')
	if len(numberStr) > 1 && (numberStr[0] == 'e' || numberStr[0] == 'E') && decimal.Value != 0 {
		// Parse positiveE
		numberStr = numberStr[1:]
		positiveE := true
		if len(numberStr) > 0 && numberStr[0] == '+' {
			numberStr = numberStr[1:]
		} else if len(numberStr) > 0 && numberStr[0] == '-' {
			positiveE = false
			numberStr = numberStr[1:]
		}
		// Read number
		powerOfTen := int64(0)
		first := true
		for i := 0; i < len(numberStr); i++ {
			if '0' <= numberStr[i] && numberStr[i] <= '9' {
				if first {
					if numberStr[i] == '0' {
						continue
					}
					first = false
					powerOfTen = int64(numberStr[i] - '0')
					continue
				}
				powerOfTen = powerOfTen*10 + int64(numberStr[i]-'0')
				if powerOfTen < 0 {
					// Overflow!
					return Decimal{}, ErrorParsingOverflow
				}
			}
		}
		// Merge powers of 10
		if !positiveE {
			powerOfTen *= -1
			if decimal.PowerOfTen < 0 && powerOfTen < math.MinInt64-decimal.PowerOfTen {
				// Overflow!
				return Decimal{}, ErrorParsingOverflow
			}
		} else if 0 < decimal.PowerOfTen && math.MaxInt64-decimal.PowerOfTen < powerOfTen {
			// Overflow!
			return Decimal{}, ErrorParsingOverflow
		}
		decimal.PowerOfTen += powerOfTen
	}
	return decimal, nil
}

// Implementation of the TextUnmarshaler interface using the ParseString function
func (d *Decimal) UnmarshalText(text []byte) (err error) {
	if d == nil {
		return ErrorNilPointer
	}
	*d, err = ParseString(string(text))
	return err
}
