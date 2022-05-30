// Package numbers provides utility functions for numbers
package numbers

import "strconv"

// Float64Str takes a float and returns a string of the minimal representation.
// e.g. 2.5000 becomes "2.5"
func Float64Str(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
