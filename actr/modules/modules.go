// Package modules implements several ACT-R modules.
package modules

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

const BuiltIn = "built-in"

// Module is an ACT-R module
type Module struct {
	Name        string
	Version     string
	Description string

	BufferList buffer.List

	Params ParamInfoMap
}

// ModuleInterface provides an interface for the ACT-R concept of a "module".
type ModuleInterface interface {
	ModuleName() string
	ModuleVersion() string
	ModuleDescription() string

	Buffers() buffer.List

	HasParameters() bool
	Parameters() []ParamInterface
	ParameterInfo(name string) ParamInterface

	ValidateParam(param *params.Param) error
	SetParam(param *params.Param) error
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

func (m Module) Buffers() buffer.List {
	return m.BufferList
}

func (m Module) HasParameters() bool {
	return len(m.Params) > 0
}

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

// AllModules returns a slice of all the modules
func AllModules() (modules []ModuleInterface) {
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
func FindModule(name string) ModuleInterface {
	for _, module := range AllModules() {
		if module.ModuleName() == name {
			return module
		}
	}

	return nil
}
