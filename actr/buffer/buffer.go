// Package buffer implements ACT-R buffers.
package buffer

import (
	"fmt"
	"slices"
	"strings"

	"github.com/asmaloney/gactar/actr/param"
	"github.com/asmaloney/gactar/util/keyvalue"
)

// validStates is a list of valid buffer states to use when matching
// TODO: needs review and correction
// See: https://github.com/asmaloney/gactar/discussions/221
var validStates = []string{
	"empty", // the buffer does not contain a chunk and the failure flag is clear
	"full",  // there is currently a chunk in the buffer
}

const (
	BuiltIn    = true
	NotBuiltIn = false
)

type buffer struct {
	name string

	// Keeps track of whether this is a built-in buffer or not
	builtIn bool

	// "spreading_activation": see "Spreading Activation" in "ACT-R 7.26 Reference Manual" pg. 290
	spreadingActivation float64
	// The defaultSpreadingActivation is set when creating the buffer and may be overridden in the config.
	defaultSpreadingActivation float64

	// Any parameters this buffer has
	parameters param.ParametersInterface

	// Request parameters (if any)
	requestParameters param.ParametersInterface
}

type Interface interface {
	Name() string
	IsBuiltIn() bool

	DefaultSpreadingActivation() float64
	SpreadingActivation() float64

	Parameters() param.ParametersInterface
	SetParam(param *keyvalue.KeyValue) error

	RequestParameters() param.ParametersInterface
	SetRequestParam(param *keyvalue.KeyValue) error
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

func NewBuffer(name string, builtIn bool, spreadingActivation float64, requestParameters param.ParametersInterface) Interface {
	spreadingActivationParam := param.NewFloat(
		"spreading_activation",
		"spreading activation weight",
		param.Ptr(0.0), nil,
	)

	parameters := param.NewParameters(param.List{
		spreadingActivationParam,
	})

	newBuffer := &buffer{
		name:                       name,
		builtIn:                    builtIn,
		defaultSpreadingActivation: spreadingActivation,
		parameters:                 parameters,
		requestParameters:          requestParameters,
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

func (b buffer) Name() string                          { return b.name }
func (b buffer) IsBuiltIn() bool                       { return b.builtIn }
func (b buffer) SpreadingActivation() float64          { return b.spreadingActivation }
func (b buffer) DefaultSpreadingActivation() float64   { return b.defaultSpreadingActivation }
func (b buffer) Parameters() param.ParametersInterface { return b.parameters }

func (b *buffer) SetParam(param *keyvalue.KeyValue) (err error) {
	err = b.Parameters().ValidateParam(param)
	if err != nil {
		return
	}

	value := param.Value

	if param.Key == "spreading_activation" {
		b.spreadingActivation = *value.Number
	}

	return
}

func (b buffer) RequestParameters() param.ParametersInterface { return b.requestParameters }

func (b buffer) SetRequestParam(param *keyvalue.KeyValue) (err error) {
	err = b.RequestParameters().ValidateParam(param)
	if err != nil {
		return
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
