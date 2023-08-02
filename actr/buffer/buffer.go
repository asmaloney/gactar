// Package buffer implements ACT-R buffers.
package buffer

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
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

type Buffer struct {
	Name string

	// Optional list of request parameter names
	RequestParameters []string
}

type Interface interface {
	BufferName() string

	HasRequestParameters() bool
	RequestParameterNames() []string
	IsValidRequestKey(key string) bool
	ValidateRequestParameter(key, value string) error
}

// List is a list of buffers that we provide operations on
type List []Interface

// ListInterface defines functions we can call on a List
type ListInterface interface {
	Count() int
	Names() []string
	Has(name string) bool
	At(index int) *Buffer
	Lookup(name string) Interface
}

func (b Buffer) BufferName() string {
	return b.Name
}

// HasRequestParameters returns whether this buffer has any additional request parameters available.
func (b Buffer) HasRequestParameters() bool {
	return len(b.RequestParameters) > 0
}

// RequestParameterNames returns a list of request parameter names.
func (b Buffer) RequestParameterNames() []string {
	return b.RequestParameters
}

// IsValidRequestKey checks if 'key' is a valid request parameter for this buffer.
func (b Buffer) IsValidRequestKey(key string) bool {
	return slices.Contains(b.RequestParameters, key)
}

// ValidateRequestParameter checks if 'param' is a valid request parameter for this buffer.
// Other buffers should implement this to check types and so on.
func (b Buffer) ValidateRequestParameter(key, value string) error {
	return nil
}

func (b Buffer) String() string {
	return b.Name
}

// Count returns the number of buffers in the list
func (b List) Count() int {
	return len(b)
}

// Names returns the list of buffer names
func (b List) Names() (names []string) {
	for _, buff := range b {
		names = append(names, buff.BufferName())
	}

	return
}

// Has returns true if the buffer "name" exists
func (b List) Has(name string) bool {
	names := b.Names()

	return slices.Contains(names, name)
}

// At returns the buffer at "index" or nil if out of range
func (m List) At(index int) Interface {
	if index < 0 || index > len(m) {
		return nil
	}

	return m[index]
}

// Lookup looks up a buffer by name (returns nil if not found)
func (b List) Lookup(name string) Interface {
	for _, buff := range b {
		if buff.BufferName() == name {
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
