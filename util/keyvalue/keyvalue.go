// Package keyvalue implements parsed key/value information.
package keyvalue

import (
	"fmt"
	"slices"
	"strings"
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
	Fields *[]KeyValue
}

// KeyValue is the key/value of a parameter from the parsed amod code.
type KeyValue struct {
	Key   string
	Value Value
}

func (v Value) IsSet() bool {
	return v.ID != nil || v.Str != nil || v.Number != nil || v.Fields != nil
}

func (v Value) String() string {
	switch {
	case v.ID != nil:
		return *v.ID

	case v.Str != nil:
		return *v.Str

	case v.Number != nil:
		return fmt.Sprintf("%f", *v.Number)

	case v.Fields != nil:
		fieldsStr := []string{}
		for _, field := range *v.Fields {
			fieldsStr = append(fieldsStr, field.Value.String())
		}

		return fmt.Sprintf("{ %s }", strings.Join(fieldsStr, ", "))
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

	case v.Fields != nil:
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
