package params

// Value mimics amod.fieldValue but without tokens.
type Value struct {
	ID     *string
	Str    *string
	Number *float64
}

type Param struct {
	Key   string
	Value Value
}

type ParamError = int

const (
	NoError ParamError = iota
	NumberRequired
	NumberMustBePositive

	UnrecognizedParam
)
