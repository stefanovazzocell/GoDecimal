package decimal

import (
	"strconv"
	"strings"
)

// Formats a number as a string, with the following parameters:
//   - if true, asDecimal uses the power-of-ten notation with "e"
//   - if true, asPercentage prints the number as a percentage and ends it with "%"
//   - accuracyLimit limts the maximum number of significant digits printed out; <=0 means unlimited
//
// Examples:
//   - {true, 123, -2}.Format(true, false, 2): "12e-1"
//   - {true, 123, -2}.Format(false, false, 0): "1.23"
//   - {true, 123, -2}.Format(false, true, 0): "123"
func (d Decimal) Format(asDecimal bool, asPercentage bool, accuracyLimit int) string {
	// Prepare number
	d.Compress()
	// Parepare core
	core := strconv.FormatUint(d.Value, 10)
	powerDelta := uint64(0)
	if accuracyLimit > 0 && len(core) > accuracyLimit {
		powerDelta = uint64(len(core) - accuracyLimit)
		core = core[:accuracyLimit]
	}
	// grab uint64 powerOfTen
	powerOfTenPositive := true
	if d.PowerOfTen < 0 {
		powerOfTenPositive = false
		d.PowerOfTen *= -1
	}
	powerOfTen := uint64(d.PowerOfTen)
	// Adjust power of 10 based on core changes
	if powerOfTenPositive {
		powerOfTen += powerDelta
	} else if powerOfTen > powerDelta {
		powerOfTen -= powerDelta
	} else {
		powerOfTenPositive = true
		powerOfTen = powerDelta - powerOfTen
	}
	// If we want a percentage, adjust accordingly
	if asPercentage {
		if powerOfTenPositive {
			powerOfTen += 2
		} else if powerOfTen > 2 {
			powerOfTen -= 2
		} else {
			powerOfTenPositive = true
			powerOfTen = 2 - powerOfTen
		}
	}
	// Prepare builder and add '-' if necessary
	resultBuilder := strings.Builder{}
	if !d.Sign {
		resultBuilder.WriteByte('-')
	}
	// Process core
	if asDecimal {
		// Write this as a decimal
		resultBuilder.WriteString(core)
		if powerOfTen != 0 {
			resultBuilder.WriteByte('e')
			if !powerOfTenPositive {
				resultBuilder.WriteByte('-')
			}
			resultBuilder.WriteString(strconv.FormatUint(powerOfTen, 10))
		}
	} else {
		// Core part
		if !powerOfTenPositive {
			// 0.00123 or 1.23
			if powerOfTen < uint64(len(core)) {
				resultBuilder.WriteString(core[0:powerOfTen])
				resultBuilder.WriteByte('.')
			} else {
				// Write first zeroes
				resultBuilder.WriteString("0.")
				for i := powerOfTen - uint64(len(core)); i > 0; i-- {
					resultBuilder.WriteByte('0')
				}
				powerOfTen = uint64(len(core))
			}
			if powerOfTen <= uint64(len(core)) {
				resultBuilder.WriteString(core[uint64(len(core))-powerOfTen:])
			}
		} else {
			// 123
			resultBuilder.WriteString(core)
		}
		// "000" end
		if powerOfTenPositive {
			// Add zeroes
			for i := uint64(0); i < powerOfTen; i++ {
				resultBuilder.WriteByte('0')
			}
		}
	}
	// "%" end
	if asPercentage {
		resultBuilder.WriteByte('%')
	}
	return resultBuilder.String()
}
