// Package modules implements several ACT-R modules.
package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
	"github.com/asmaloney/gactar/util/container"
)

// ModuleInterface provides an interface for the ACT-R concept of a "module".
type ModuleInterface interface {
	ModuleName() string

	NumBuffers() int
	BufferNames() []string
	HasBuffer(name string) bool
	OnlyBuffer() *buffer.Buffer
	LookupBuffer(name string) buffer.BufferInterface

	SetParam(param *params.Param) (err params.ParamError)
}

type Module struct {
	Name    string
	Buffers []buffer.Buffer
}

func (m Module) ModuleName() string {
	return m.Name
}

func (m Module) NumBuffers() int {
	return len(m.Buffers)
}

func (m Module) BufferNames() (names []string) {
	for _, buff := range m.Buffers {
		names = append(names, buff.BufferName())
	}

	return
}

func (m Module) HasBuffer(name string) bool {
	names := m.BufferNames()

	return container.Contains(name, names)
}

func (m Module) OnlyBuffer() *buffer.Buffer {
	if m.Buffers == nil || len(m.Buffers) > 1 {
		return nil
	}

	return &m.Buffers[0]
}

func (m Module) LookupBuffer(name string) buffer.BufferInterface {
	for _, buff := range m.Buffers {
		if buff.Name == name {
			return buff
		}
	}

	return nil
}
