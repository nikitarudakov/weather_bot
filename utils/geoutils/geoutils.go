package geoutils

import "strconv"

func FormatCoordinateToString(c float32) string {
	return strconv.FormatFloat(float64(c), 'f', -1, 32)
}
