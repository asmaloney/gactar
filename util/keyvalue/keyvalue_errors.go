package keyvalue

import (
	"fmt"
)

// ValueType is used when outputting "invalid type" errors
type ValueType int

const (
	Boolean ValueType = iota
	Number
)

func (v ValueType) String() string {
	switch v {
	case Boolean:
		return "boolean"
	case Number:
		return "number"
	}

	return "unknown"
}

type ErrInvalidType struct {
	ExpectedType ValueType
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
