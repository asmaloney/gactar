package params

import (
	"github.com/asmaloney/gactar/util/container"
)

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

type Param struct {
	Key   string
	Value Value
}

func (v Value) AsBool() (bool, error) {
	if (v.ID == nil) || !container.Contains(*v.ID, boolean) {
		return false, ErrInvalidType{ExpectedType: Boolean}
	}

	return *v.ID == "true", nil
}
