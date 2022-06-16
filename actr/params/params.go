package params

var Boolean = []string{
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
	NumberRequired
	NumberMustBePositive

	InvalidOption

	UnrecognizedParam
)

func BooleanStrToBool(b string) bool {
	return b == "true"
}
