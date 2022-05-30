package actr

import (
	"fmt"
	"sort"
	"strings"
)

// ValidBufferStates is a list of the valid buffer states to use with the _status chunk
var ValidBufferStates = map[string]bool{
	"empty": true,
	"full":  true,
	"busy":  true,
	"error": true,
}

type BufferInterface interface {
	GetBufferName() string
	AllowsMultipleInit() bool
}

type Buffer struct {
	Name string

	MultipleInit bool
}

func (b Buffer) GetBufferName() string {
	return b.Name
}

func (b Buffer) AllowsMultipleInit() bool {
	return b.MultipleInit
}

func (b Buffer) String() string {
	return b.Name
}

// IsValidBufferState checks if 'state' is a valid buffer state.
func IsValidBufferState(state string) bool {
	v, ok := ValidBufferStates[state]
	return v && ok
}

// ValidBufferStatesStr returns a list of (sorted) valid buffer states. Used for error output.
func ValidBufferStatesStr() string {
	keys := make([]string, 0, len(ValidBufferStates))
	for k := range ValidBufferStates {
		keys = append(keys, fmt.Sprintf("'%s'", k))
	}

	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
