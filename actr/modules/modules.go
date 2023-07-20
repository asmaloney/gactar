// Package modules implements several ACT-R modules.
package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

// Module is an ACT-R module
type Module struct {
	Name       string
	BufferList buffer.List
}

// ModuleInterface provides an interface for the ACT-R concept of a "module".
type ModuleInterface interface {
	ModuleName() string

	Buffers() buffer.List

	SetParam(param *params.Param) (err error)
}

func (m Module) ModuleName() string {
	return m.Name
}

func (m Module) Buffers() buffer.List {
	return m.BufferList
}
