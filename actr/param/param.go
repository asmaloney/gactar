package param

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/asmaloney/gactar/util/keyvalue"
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

// Info is the basic info about a parameter
type Info struct {
	Name        string
	Description string
}

// Int is an int parameter with optional min and max constraints
type Int struct {
	Info

	Min *int
	Max *int
}

// Float is a float parameter with optional min and max constraints
type Float struct {
	Info

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

// InfoMap maps a name to the parameter's info
type InfoMap map[string]ParamInterface

// List is a slice of ParamInterface
type List []ParamInterface

type parameters struct {
	params InfoMap
}

// Add a ParametersInterface to a struct to store and validate parameters
type ParametersInterface interface {
	ParameterList() List

	ValidateParam(param *keyvalue.KeyValue) error
}

func (p Info) GetName() string {
	return p.Name
}

func (p Info) GetDescription() string {
	return p.Description
}

func (p Int) GetMin() *float64 {
	if p.Min != nil {
		temp := float64(*p.Min)
		return &temp
	}
	return nil
}

func (p Int) GetMax() *float64 {
	if p.Max != nil {
		temp := float64(*p.Max)
		return &temp
	}
	return nil
}

func (p Float) GetMin() *float64 { return p.Min }
func (p Float) GetMax() *float64 { return p.Max }

func NewParameters(paramMap InfoMap) ParametersInterface {
	return parameters{params: paramMap}
}

// Parameters returns a slice of parameters sorted by name
func (p parameters) ParameterList() List {
	params := maps.Values(p.params)

	slices.SortFunc[ParamInterface](params, func(a, b ParamInterface) bool { return a.GetName() < b.GetName() })

	return params
}

// ValidateParam given an actr param will validate it against our parameter info
func (p parameters) ValidateParam(param *keyvalue.KeyValue) (err error) {
	paramInfo := p.parameterInfo(param.Key)
	if paramInfo == nil {
		return ErrUnrecognizedOption{Option: param.Key}
	}

	min := paramInfo.GetMin()
	max := paramInfo.GetMax()

	value := param.Value

	// we currently only have numbers
	if value.Number == nil {
		return keyvalue.ErrInvalidType{ExpectedType: keyvalue.Number}
	}

	if (min != nil) && (max != nil) &&
		((*value.Number < *min) || (*value.Number > *max)) {
		return ErrValueOutOfRange{
			Min: min,
			Max: max,
		}
	}

	if min != nil && (*value.Number < *min) {
		return ErrValueOutOfRange{
			Min: min,
		}
	}

	if max != nil && (*value.Number > *max) {
		return ErrValueOutOfRange{
			Max: max,
		}
	}

	return
}

// parameterInfo returns detailed info about a specific parameter given by "name"
func (p parameters) parameterInfo(name string) ParamInterface {
	info, ok := p.params[name]
	if ok {
		return info
	}

	return nil
}

// NewInt creates a new int param with optional min/max constraints
func NewInt(name, description string, min, max *int) Int {
	return Int{
		Info{name, description},
		min, max,
	}
}

// NewFloat creates a new float param with optional min/max constraints
func NewFloat(name, description string, min, max *float64) Float {
	return Float{
		Info{name, description},
		min, max,
	}
}
