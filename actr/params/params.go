// Package params implements parsed parameter information.
package params

import "golang.org/x/exp/slices"

var boolean = []string{
	"true",
	"false",
}

// Value mimics amod.fieldValue but without tokens.
type Value struct {
	ID     *string
	Str    *string
	Number *float64
	Field  *Param
}

// Param is the key/value of a parameter from the parsed amod code.
type Param struct {
	Key   string
	Value Value
}

func (v Value) AsBool() (bool, error) {
	if (v.ID == nil) || !slices.Contains(boolean, *v.ID) {
		return false, ErrInvalidType{ExpectedType: Boolean}
	}

	return *v.ID == "true", nil
}
