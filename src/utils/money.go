package utils

import (
	"fmt"
	"math"
)

// ToKobo converts a naira amount (float64) to kobo (int64), rounding to the nearest kobo.
func ToKobo(naira float64) int64 {
	return int64(math.Round(naira * 100))
}

// FromKobo converts kobo (int64) to naira (float64).
func FromKobo(kobo float64) float64 {
	return float64(kobo) / 100.0
}

// FormatNaira returns a string with 2 decimal places, e.g., 20.00
func FormatNaira(naira float64) string {
	return fmt.Sprintf("%.2f", naira)
}
