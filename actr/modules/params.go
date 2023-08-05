package modules

import (
	"fmt"
	"strings"

	"github.com/asmaloney/gactar/util/numbers"
)

type ErrUnrecognizedOption struct {
	Option string
}

func (e ErrUnrecognizedOption) Error() string {
	return fmt.Sprintf("unrecognized option %q", e.Option)
}

type ErrValueOutOfRange struct {
	Min *float64
	Max *float64
}

func (e ErrValueOutOfRange) Error() string {
	if e.Min != nil && e.Max == nil {
		return fmt.Sprintf("is out of range (minimum %s)", numbers.Float64Str(*e.Min))
	}

	if e.Min == nil && e.Max != nil {
		return fmt.Sprintf("is out of range (maximum %s)", numbers.Float64Str(*e.Max))
	}

	return fmt.Sprintf("is out of range (%s-%s)", numbers.Float64Str(*e.Min), numbers.Float64Str(*e.Max))
}

type ErrInvalidValue struct {
	Expected []string
}

func (e ErrInvalidValue) Error() string {
	return fmt.Sprintf("must be must be one of %q", strings.Join(e.Expected, ", "))
}

// Ptr simply returns a pointer to a literal. e.g. Ptr(0.5)
// This is useful when passing literals to functions which require pointers to basic types.
func Ptr[T any](v T) *T {
	return &v
}

// ParamInfo is the basic info about a parameter
type ParamInfo struct {
	Name        string
	Description string
}

// ParamInt is an int parameter with optional min and max constraints
type ParamInt struct {
	ParamInfo

	Min *int
	Max *int
}

// ParamFloat is a float parameter with optional min and max constraints
type ParamFloat struct {
	ParamInfo

	Min *float64
	Max *float64
}

// ParamInterface provides an interface to a parameter
type ParamInterface interface {
	GetName() string
	GetDescription() string

	GetMin() *float64
	GetMax() *float64
}

// ParamInfoMap maps a name to the parameter's info
type ParamInfoMap map[string]ParamInterface

func (p ParamInfo) GetName() string {
	return p.Name
}

func (p ParamInfo) GetDescription() string {
	return p.Description
}

func (p ParamInt) GetMin() *float64 {
	if p.Min != nil {
		temp := float64(*p.Min)
		return &temp
	}
	return nil
}

func (p ParamInt) GetMax() *float64 {
	if p.Max != nil {
		temp := float64(*p.Max)
		return &temp
	}
	return nil
}

func (p ParamFloat) GetMin() *float64 { return p.Min }
func (p ParamFloat) GetMax() *float64 { return p.Max }

// NewParamInt creates a new int param with optional min/max constraints
func NewParamInt(name, description string, min, max *int) ParamInt {
	return ParamInt{
		ParamInfo{name, description},
		min, max,
	}
}

// NewParamFloat creates a new float param with optional min/max constraints
func NewParamFloat(name, description string, min, max *float64) ParamFloat {
	return ParamFloat{
		ParamInfo{name, description},
		min, max,
	}
}
