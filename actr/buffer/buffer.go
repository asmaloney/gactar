// Package buffer implements ACT-R buffers.
package buffer

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/asmaloney/gactar/actr/param"
	"github.com/asmaloney/gactar/util/keyvalue"
)

// validStates is a list of valid buffer states to use when matching
// TODO: needs review and correction
// See: https://github.com/asmaloney/gactar/discussions/221
var validStates = []string{
	"empty",
	"full",
}

type ErrInvalidRequestParameterValue struct {
	ParameterName string
	Value         string
	Context       *string // optional context
}

func (e ErrInvalidRequestParameterValue) Error() string {
	message := fmt.Sprintf("invalid value %q for parameter %q", e.Value, e.ParameterName)

	if e.Context != nil {
		message += fmt.Sprintf(" %s", *e.Context)
	}

	return message
}

type buffer struct {
	name string

	// "spreading_activation": see "Spreading Activation" in "ACT-R 7.26 Reference Manual" pg. 290
	// The default is set when creating the buffer and may be overridden in the config.
	spreadingActivation float64

	// List of valid request parameter names (if any)
	requestParameters []string

	param.ParametersInterface
}

type Interface interface {
	Name() string

	SetSpreadingActivation(activation float64)
	SpreadingActivation() float64

	SetValidRequestParameters(params []string)
	HasRequestParameters() bool
	RequestParameterKeys() []string
	IsValidRequestKey(key string) bool
	ValidateRequestParameter(key, value string) error

	Parameters() param.ParametersInterface
	SetParam(param *keyvalue.KeyValue) error
}

// List is a list of buffers that we provide operations on
type List []Interface

// ListInterface defines functions we can call on a List
type ListInterface interface {
	Count() int
	Names() []string
	Has(name string) bool
	At(index int) *Interface
	Lookup(name string) Interface
}

func NewBuffer(name string, spreadingActivation float64) Interface {
	spreadingActivationParam := param.NewFloat(
		"spreading_activation",
		"spreading activation weight",
		param.Ptr(0.0), nil,
	)

	parameters := param.NewParameters(param.InfoMap{
		"spreading_activation": spreadingActivationParam,
	})

	newBuffer := &buffer{
		name:                name,
		ParametersInterface: parameters,
	}

	err := newBuffer.SetParam(&keyvalue.KeyValue{
		Key:   "spreading_activation",
		Value: keyvalue.Value{Number: &spreadingActivation},
	})

	if err != nil {
		panic(fmt.Sprintf("INTERNAL: invalid default value for %q on buffer %q (%s)",
			"spreading_activation",
			name,
			err,
		))
	}

	return newBuffer
}

func (b buffer) Name() string {
	return b.name
}

func (b *buffer) SetSpreadingActivation(activation float64) { b.spreadingActivation = activation }

func (b buffer) SpreadingActivation() float64 {
	return b.spreadingActivation
}

// SetValidRequestParameters sets the list of valid request parameter keys.
func (b *buffer) SetValidRequestParameters(params []string) {
	b.requestParameters = params
}

// RequestParameterKeys returns a list of request parameter keys.
func (b buffer) RequestParameterKeys() []string {
	return b.requestParameters
}

// HasRequestParameters returns whether this buffer has any additional request parameters available.
func (b buffer) HasRequestParameters() bool {
	return len(b.requestParameters) > 0
}

// IsValidRequestKey checks if 'key' is a valid request parameter for this buffer.
func (b buffer) IsValidRequestKey(key string) bool {
	return slices.Contains(b.requestParameters, key)
}

// ValidateRequestParameter checks if 'param' is a valid request parameter for this buffer.
// Other buffers should implement this to check types and so on.
func (b buffer) ValidateRequestParameter(key, value string) error { return nil }

func (b buffer) Parameters() param.ParametersInterface {
	return b.ParametersInterface
}

func (b buffer) SetParam(param *keyvalue.KeyValue) (err error) {
	err = b.ValidateParam(param)
	if err != nil {
		return
	}

	value := param.Value

	if param.Key == "spreading_activation" {
		b.spreadingActivation = *value.Number
	}

	return
}

// Count returns the number of buffers in the list
func (l List) Count() int {
	return len(l)
}

// Names returns the list of buffer names
func (l List) Names() (names []string) {
	for _, buff := range l {
		names = append(names, buff.Name())
	}

	return
}

// Has returns true if the buffer "name" exists
func (l List) Has(name string) bool {
	names := l.Names()

	return slices.Contains(names, name)
}

// At returns the buffer at "index" or nil if out of range
func (l List) At(index int) Interface {
	if index < 0 || index > len(l) {
		return nil
	}

	return l[index]
}

// Lookup looks up a buffer by name (returns nil if not found)
func (l List) Lookup(name string) Interface {
	for _, buff := range l {
		if buff.Name() == name {
			return buff
		}
	}

	return nil
}

// IsValidState checks if 'state' is a valid buffer state.
func IsValidState(state string) bool {
	return slices.Contains(validStates, state)
}

// ValidStatesStr returns a list of (sorted) valid buffer states. Used for error output.
func ValidStatesStr() string {
	return strings.Join(validStates, ", ")
}
