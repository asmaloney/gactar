package params

import "github.com/asmaloney/gactar/util/container"

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

type ParamError = int

const (
	NoError ParamError = iota

	BoolRequired

	NumberRequired
	NumberMustBePositive
	NumberOutOfRange

	InvalidOption

	UnrecognizedParam
)

func (v Value) AsBool() (bool, ParamError) {
	if (v.ID == nil) || !container.Contains(*v.ID, boolean) {
		return false, BoolRequired
	}

	return *v.ID == "true", NoError
}
