package util

import (
	"strconv"
)

//mod11 calculate Mod11 DV from string
func mod11(valueSequence string) int {
	digit := 0
	sum := 0
	weight := 2

	var values []int

	for _, r := range valueSequence {
		c := string(r)
		n, _ := strconv.Atoi(c)
		values = append(values, n)
	}
	for i := len(values) - 1; i >= 0; i-- {
		sum += values[i] * weight

		if weight < 9 {
			weight = weight + 1
		} else {
			weight = 2
		}
	}
	digit = 11 - (sum % 11)
	return digit
}

//OurNumberDv calculate DV from OurNumber
func OurNumberDv(valueSequence string) string {
	digit := mod11(valueSequence)
	if digit > 9 {
		digit = 0
	}
	return strconv.Itoa(digit)
}

//BarcodeDv calculate DV from barcode
func BarcodeDv(valueSequence string) string {
	digit := mod11(valueSequence)
	if digit <= 1 || digit > 9 {
		digit = 1
	}
	return strconv.Itoa(digit)
}
