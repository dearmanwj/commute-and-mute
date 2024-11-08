package serialization

func GetNilFromNegative1(number float64) *float64 {
	var nillable *float64
	if number != -1 {
		nillable = &number
	}
	return nillable
}
