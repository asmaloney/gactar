// Package modules implements a module interface and several built-in ACT-R modules.
package modules

import (
	"strings"

	"golang.org/x/exp/slices"

	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/param"
	"github.com/asmaloney/gactar/util/keyvalue"
)

const BuiltIn = "built-in"

var validStates = []string{
	"busy",
	"error",
	"free",
}

// Module is an ACT-R module
type Module struct {
	Name        string
	Version     string
	Description string

	BufferList buffer.List // The buffers this module provides (may be empty)

	MultipleInit bool

	param.ParametersInterface
}

// Interface provides an interface for the ACT-R concept of a "module".
type Interface interface {
	ModuleName() string
	ModuleVersion() string
	ModuleDescription() string

	HasBuffers() bool
	Buffers() buffer.List

	Parameters() param.ParametersInterface
	SetParam(param *keyvalue.KeyValue) error

	AllowsMultipleInit() bool
}

func (m Module) ModuleName() string {
	return m.Name
}

func (m Module) ModuleVersion() string {
	return m.Version
}

func (m Module) ModuleDescription() string {
	return m.Description
}

func (m Module) HasBuffers() bool {
	return len(m.BufferList) > 0
}

func (m Module) Buffers() buffer.List {
	return m.BufferList
}

func (m Module) Parameters() param.ParametersInterface {
	return m.ParametersInterface
}

func (m Module) SetParam(param *keyvalue.KeyValue) error { return nil }

// AllowsMultipleInit returns whether this module allows more than one initialization.
// e.g. goal would only allow one, whereas declarative memory would allow multiple.
func (m Module) AllowsMultipleInit() bool {
	return m.MultipleInit
}

// AllModules returns a slice of all the modules
func AllModules() (modules []Interface) {
	modules = append(modules, NewExtraBuffers())
	modules = append(modules, NewGoal())
	modules = append(modules, NewImaginal())
	modules = append(modules, NewDeclarativeMemory())
	modules = append(modules, NewProcedural())

	return
}

// ModuleNames returns a slice of all the module names
func ModuleNames() (names []string) {
	for _, module := range AllModules() {
		names = append(names, module.ModuleName())
	}

	return names
}

// FindModule finds a module by name or returns nil if not found
func FindModule(name string) Interface {
	for _, module := range AllModules() {
		if module.ModuleName() == name {
			return module
		}
	}

	return nil
}

// IsValidState checks if 'state' is a valid module state.
func IsValidState(state string) bool {
	return slices.Contains(validStates, state)
}

// ValidStatesStr returns a list of valid module states. Used for error output.
func ValidStatesStr() string {
	return strings.Join(validStates, ", ")
}
