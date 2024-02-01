// Package param provides tools to declare and validate parameters.
package param

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"

	"github.com/asmaloney/gactar/util/keyvalue"
	"github.com/asmaloney/gactar/util/numbers"
)

type number interface {
	constraints.Integer | constraints.Float
}

type ErrUnrecognizedOption struct {
	Option string
}

func (e ErrUnrecognizedOption) Error() string {
	return fmt.Sprintf("unrecognized option %q", e.Option)
}

type ErrValueOutOfRange struct {
	Min *string
	Max *string
}

func (e ErrValueOutOfRange) Error() string {
	if e.Min != nil && e.Max == nil {
		return fmt.Sprintf("is out of range (minimum %s)", *e.Min)
	}

	if e.Min == nil && e.Max != nil {
		return fmt.Sprintf("is out of range (maximum %s)", *e.Max)
	}

	return fmt.Sprintf("is out of range (%s-%s)", *e.Min, *e.Max)
}

type ErrInvalidType struct {
	FoundType    string
	ExpectedType string
}

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("invalid type (found %s; expected %s)", e.FoundType, e.ExpectedType)
}

type ErrInvalidValue struct {
	ParameterName string
	Value         string
	Context       *string // optional context
}

func (e ErrInvalidValue) Error() string {
	message := fmt.Sprintf("invalid value %q for option %q", e.Value, e.ParameterName)

	if e.Context != nil {
		message += fmt.Sprintf(" %s", *e.Context)
	}

	return message
}

// Ptr simply returns a pointer to a literal. e.g. Ptr(0.5)
// This is useful when passing literals to functions which require pointers to basic types.
func Ptr[T any](v T) *T {
	return &v
}

// info is the basic info about a parameter
type info struct {
	name        string
	description string
}

func (i info) Name() string        { return i.name }
func (i info) Description() string { return i.description }

// Bool is a boolean parameter
type Bool struct {
	info
}

// NewBool creates a new boolean param
func NewBool(name, description string) Bool {
	return Bool{
		info: info{name, description},
	}
}

// Str is a string parameter
type Str struct {
	info

	validValues []string
}

func (s Str) ValidValues() []string { return s.validValues }

// NewStr creates a new string param with optional list of valid values
func NewStr(name, description string, validValues []string) Str {
	return Str{
		info:        info{name, description},
		validValues: validValues,
	}
}

// Int is an int parameter with optional min and max constraints
type Int struct {
	info

	min *int
	max *int
}

func (i Int) Min() *int { return i.min }
func (i Int) Max() *int { return i.max }

// NewInt creates a new int param with optional min/max constraints
func NewInt(name, description string, min, max *int) Int {
	return Int{
		info{name, description},
		min, max,
	}
}

// Float is a float parameter with optional min and max constraints
type Float struct {
	info

	min *float64
	max *float64
}

func (f Float) Min() *float64 { return f.min }
func (f Float) Max() *float64 { return f.max }

// NewFloat creates a new float param with optional min/max constraints
func NewFloat(name, description string, min, max *float64) Float {
	return Float{
		info{name, description},
		min, max,
	}
}

// ParamInterface provides an interface to a parameter
type ParamInterface interface {
	Name() string
	Description() string
}

// List is a slice of ParamInterface
type List []ParamInterface

// infoMap maps a name to the parameter's info for easy lookup
type infoMap map[string]ParamInterface

type parameters struct {
	params infoMap // store as a map for easy lookup
}

// Add a ParametersInterface to a struct to store and validate parameters
type ParametersInterface interface {
	ParameterList() List

	ValidateParam(param *keyvalue.KeyValue) error
}

func NewParameters(params List) ParametersInterface {
	iMap := make(infoMap, len(params))

	for _, param := range params {
		iMap[param.Name()] = param
	}

	return parameters{params: iMap}
}

// Parameters returns a slice of parameters sorted by name
func (p parameters) ParameterList() List {
	params := maps.Values(p.params)

	compareParamInterface := func(a, b ParamInterface) int {
		return strings.Compare(a.Name(), b.Name())
	}

	slices.SortFunc(params, compareParamInterface)

	return params
}

// ValidateParam given an actr param will validate it against our parameter info
func (p parameters) ValidateParam(param *keyvalue.KeyValue) (err error) {
	paramInfo := p.parameterInfo(param.Key)
	if paramInfo == nil {
		return ErrUnrecognizedOption{Option: param.Key}
	}

	value := param.Value
	valueType := value.Type()

	switch pInfo := paramInfo.(type) {
	case Str:
		if value.Str == nil {
			return ErrInvalidType{
				FoundType:    valueType,
				ExpectedType: "string",
			}
		}

		valid := pInfo.validValues
		if (len(valid) > 0) && !slices.Contains(valid, *value.Str) {
			context := fmt.Sprintf("(expected one of: %s)", strings.Join(valid, ", "))

			return ErrInvalidValue{
				ParameterName: param.Key,
				Value:         *value.Str,
				Context:       &context,
			}
		}

	case Int:
		if value.Number == nil {
			return ErrInvalidType{
				FoundType:    valueType,
				ExpectedType: "number",
			}
		}

		val := int(*value.Number)

		err = compareMinMax(val, pInfo.min, pInfo.max)
		if err != nil {
			return
		}

	case Float:
		if value.Number == nil {
			return ErrInvalidType{
				FoundType:    valueType,
				ExpectedType: "number",
			}
		}

		val := *value.Number

		err = compareMinMax(val, pInfo.min, pInfo.max)
		if err != nil {
			return
		}

	case Bool:
		if value.ID == nil {
			return ErrInvalidType{
				FoundType:    valueType,
				ExpectedType: "true or false",
			}
		}

		if !slices.Contains(keyvalue.BooleanValues, *value.ID) {
			context := fmt.Sprintf("(expected one of: %s)", strings.Join(keyvalue.BooleanValues, ", "))

			return ErrInvalidValue{
				ParameterName: param.Key,
				Value:         *value.ID,
				Context:       &context,
			}
		}
	}

	if !value.IsSet() {
		return keyvalue.ErrInvalidType{ExpectedType: keyvalue.Number}
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

func convertNumberToStr[T number](num T) *string {
	var str string

	switch any(num).(type) {
	case int:
		str = strconv.Itoa(int(num))

	case float64:
		str = numbers.Float64Str(float64(num))
	}

	return &str
}

func compareMinMax[T number](value T, min, max *T) error {
	if (min != nil) && (max != nil) &&
		((value < *min) || (value > *max)) {
		return ErrValueOutOfRange{
			Min: convertNumberToStr[T](*min),
			Max: convertNumberToStr[T](*max),
		}
	}

	if min != nil && (value < *min) {
		return ErrValueOutOfRange{
			Min: convertNumberToStr[T](*min),
		}
	}

	if max != nil && (value > *max) {
		return ErrValueOutOfRange{
			Max: convertNumberToStr[T](*max),
		}
	}

	return nil
}
