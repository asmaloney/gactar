package params

import (
	"errors"
	"fmt"
	"strings"

	"github.com/asmaloney/gactar/util/numbers"
)

var (
	ErrUnrecognizedParam = errors.New("unrecognized option")
)

// ParamType is used when outputting "invalid type" errors
type ParamType int

const (
	Boolean ParamType = iota
	Number
)

func (p ParamType) String() string {
	switch p {
	case Boolean:
		return "boolean"
	case Number:
		return "number"
	}

	return "unknown"
}

type ErrInvalidType struct {
	ExpectedType ParamType
}

func (e ErrInvalidType) Error() string {
	expected := e.ExpectedType.String()

	if e.ExpectedType == Boolean {
		expected = "'true' or 'false'"
	} else {
		expected = fmt.Sprintf("a %s", expected)
	}

	return fmt.Sprintf("must be %s", expected)
}

type ErrInvalidOption struct {
	Expected []string
}

func (e ErrInvalidOption) Error() string {
	return fmt.Sprintf("must be must be one of %q", strings.Join(e.Expected, ", "))
}

type ErrOutOfRange struct {
	Min *float64
	Max *float64
}

func (e ErrOutOfRange) Error() string {
	if e.Min != nil && e.Max == nil {
		return fmt.Sprintf("is out of range (minimum %s)", numbers.Float64Str(*e.Min))
	}

	if e.Min == nil && e.Max != nil {
		return fmt.Sprintf("is out of range (maximum %s)", numbers.Float64Str(*e.Max))
	}

	return fmt.Sprintf("is out of range (%s-%s)", numbers.Float64Str(*e.Min), numbers.Float64Str(*e.Max))
}
