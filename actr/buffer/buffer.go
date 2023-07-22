// Package buffer implements ACT-R buffers.
package buffer

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

// validBufferStates is a list of the valid buffer states to use with the _status chunk
// TODO: needs review and correction
// See: https://github.com/asmaloney/gactar/discussions/221
var validBufferStates = map[string]bool{
	// buffer
	"empty": true,
	"full":  true,

	// state
	"busy":  true,
	"error": true,
}

type Buffer struct {
	Name string

	MultipleInit bool
}

type Interface interface {
	BufferName() string
	AllowsMultipleInit() bool
}

// List is a list of buffers that we provide operations on
type List []Buffer

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

func (b Buffer) AllowsMultipleInit() bool {
	return b.MultipleInit
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

	return &m[index]
}

// Lookup looks up a buffer by name (returns nil if not found)
func (b List) Lookup(name string) Interface {
	for _, buff := range b {
		if buff.Name == name {
			return buff
		}
	}

	return nil
}

// IsValidBufferState checks if 'state' is a valid buffer state.
func IsValidBufferState(state string) bool {
	v, ok := validBufferStates[state]
	return v && ok
}

// ValidBufferStatesStr returns a list of (sorted) valid buffer states. Used for error output.
func ValidBufferStatesStr() string {
	keys := make([]string, 0, len(validBufferStates))
	for k := range validBufferStates {
		keys = append(keys, fmt.Sprintf("'%s'", k))
	}

	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
