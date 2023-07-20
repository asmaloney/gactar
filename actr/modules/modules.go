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

	Params []ParamInfo
}

// ModuleInterface provides an interface for the ACT-R concept of a "module".
type ModuleInterface interface {
	ModuleName() string
	ModuleVersion() string
	ModuleDescription() string

	Buffers() buffer.List

	HasParameters() bool
	ParameterNames() []string
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

func (m Module) ParameterNames() (names []string) {
	for _, param := range m.Params {
		names = append(names, param.Name)
	}

	return names
}

// BuiltInModules returns a slice of all the built-in modules
func BuiltInModules() (modules []ModuleInterface) {
	modules = append(modules, NewExtraBuffers())
	modules = append(modules, NewGoal())
	modules = append(modules, NewImaginal())
	modules = append(modules, NewDeclarativeMemory())
	modules = append(modules, NewProcedural())

	return
}
