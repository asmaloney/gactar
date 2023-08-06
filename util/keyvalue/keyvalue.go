// Package keyvalue implements parsed key/value information.
package keyvalue

import (
	"fmt"

	"golang.org/x/exp/slices"
)

var BooleanValues = []string{
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

func (v Value) IsSet() bool {
	return v.ID != nil || v.Str != nil || v.Number != nil || v.Field != nil
}

func (v Value) String() string {
	switch {
	case v.ID != nil:
		return *v.ID

	case v.Str != nil:
		return *v.Str

	case v.Number != nil:
		return fmt.Sprintf("%f", *v.Number)

	case v.Field != nil:
		return fmt.Sprintf("{ %s }", *v.Field)
	}

	return ""
}

func (v Value) Type() string {
	switch {
	case v.ID != nil:
		return "id"

	case v.Str != nil:
		return "string"

	case v.Number != nil:
		return "number"

	case v.Field != nil:
		return "field"
	}

	return "<none>"
}

func (v Value) AsBool() (bool, error) {
	if (v.ID == nil) || !slices.Contains(BooleanValues, *v.ID) {
		return false, ErrInvalidType{ExpectedType: Boolean}
	}

	return *v.ID == "true", nil
}
