# GoDecimal

A decimal number representation in Go where numbers are stored using scientific notation.

Useful to do operations involving decimals that require high precision and accuracy such as those involving money.

# Testing

This library goal is to have 100% code coverage and to Fuzz all the major functions.

Please open issues on this repository if anything was missed.

To run a comprehensive system test run `make full-test`.

# Limitations

A `Decimal` is represented by an `uint64` representing it's absolute value,
a `int64` value representing a power of 10 to multiply the number to,
and finally a `bool` representing the sign.

| Type                       | Value                                |
| :------------------------- | :----------------------------------: |
| Largest Decimal            | `math.MaxUint64`*10^`math.MaxInt64`  |
| Smallest Decimal           | -`math.MaxUint64`*10^`math.MaxInt64` |
| Smallest Positive Decimal  | `1`*10^`math.MinInt64`               |
| Largest Negative Decimal   | -`1`*10^`math.MinInt64`              |

## Precision

`18446744073709551615` (`math.MaxUint64`) is the largest number we can represent in our decimal with precision.

Due to the way in which we do certain operations we have between 18 and 20 significant figures represented.