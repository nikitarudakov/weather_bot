package geoutils

import "strconv"

// FormatCoordinateToString formats coords (lan/lon) to string
func FormatCoordinateToString(c float32) string {
	return strconv.FormatFloat(float64(c), 'f', -1, 32)
}
