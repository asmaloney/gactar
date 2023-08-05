// Package keyvalue implements parsed key/value information.
package keyvalue

import "golang.org/x/exp/slices"

var boolean = []string{
	"true",
	"false",
}

// Value stores an ID, string, number, or another KeyValue (recursive).
type Value struct {
	ID     *string
	Str    *string
	Number *float64
	Field  *KeyValue
}

// KeyValue is the key/value of a parameter from the parsed amod code.
type KeyValue struct {
	Key   string
	Value Value
}

func (v Value) AsBool() (bool, error) {
	if (v.ID == nil) || !slices.Contains(boolean, *v.ID) {
		return false, ErrInvalidType{ExpectedType: Boolean}
	}

	return *v.ID == "true", nil
}
