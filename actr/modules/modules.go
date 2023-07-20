// Package modules implements several ACT-R modules.
package modules

import (
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
	ParameterNames() []string
	ParameterInfo(name string) *ParamInfo
	SetParam(param *params.Param) (err error)
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

func (m Module) ParameterInfo(name string) *ParamInfo {
	info, ok := m.Params[name]
	if ok {
		return &info
	}

	return nil
}

func (m Module) ParameterNames() (names []string) {
	for _, param := range m.Params {
		names = append(names, param.Name)
	}

	return names
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
	// TODO: With go 1.21 we can get the keys all at once with Keys
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
