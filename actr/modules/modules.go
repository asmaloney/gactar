// Package modules implements a module interface and several built-in ACT-R modules.
package modules

import (
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
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

	Params ParamInfoMap // Parameters accepted by this module (may be empty)

	MultipleInit bool
}

// Interface provides an interface for the ACT-R concept of a "module".
type Interface interface {
	ModuleName() string
	ModuleVersion() string
	ModuleDescription() string

	HasBuffers() bool
	Buffers() buffer.List

	HasParameters() bool
	Parameters() []ParamInterface
	ParameterInfo(name string) ParamInterface

	ValidateParam(param *params.Param) error
	SetParam(param *params.Param) error

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

func (m Module) HasParameters() bool {
	return len(m.Params) > 0
}

// ParameterInfo returns the detailed info about a specific parameter given by "name"
func (m Module) ParameterInfo(name string) ParamInterface {
	info, ok := m.Params[name]
	if ok {
		return info
	}

	return nil
}

// Parameters returns a slice of parameters sorted by name
func (m Module) Parameters() []ParamInterface {
	params := maps.Values(m.Params)

	slices.SortFunc[ParamInterface](params, func(a, b ParamInterface) bool { return a.GetName() < b.GetName() })

	return params
}

// ValidateParam given an actr param will validate it against our modules parameters
func (m Module) ValidateParam(param *params.Param) (err error) {

	paramInfo := m.ParameterInfo(param.Key)
	if paramInfo == nil {
		return params.ErrUnrecognizedParam
	}

	min := paramInfo.GetMin()
	max := paramInfo.GetMax()

	value := param.Value

	// we currently only have numbers
	if value.Number == nil {
		return params.ErrInvalidType{ExpectedType: params.Number}
	}

	if (min != nil) && (max != nil) &&
		((*value.Number < *min) || (*value.Number > *max)) {
		return params.ErrOutOfRange{
			Min: min,
			Max: max,
		}
	}

	if min != nil && (*value.Number < *min) {
		return params.ErrOutOfRange{
			Min: min,
		}
	}

	if max != nil && (*value.Number > *max) {
		return params.ErrOutOfRange{
			Max: max,
		}
	}

	return
}

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
