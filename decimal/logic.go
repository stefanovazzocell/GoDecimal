package decimal

// Returns true if the two decimals are equal
func (d Decimal) Equals(x Decimal) bool {
	d.Compress()
	x.Compress()
	return (d.Value == 0 && x.Value == 0) ||
		(d.Sign == x.Sign && d.Value == x.Value && d.PowerOfTen == x.PowerOfTen)
}
