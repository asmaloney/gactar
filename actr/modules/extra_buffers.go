package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

// ExtraBuffers module is used to declare one or more extra goal-style buffers in the model.
type ExtraBuffers struct {
	Module
}

// NewExtraBuffers creates and returns a new ExtraBuffers module
func NewExtraBuffers() *ExtraBuffers {
	return &ExtraBuffers{
		Module: Module{
			Name:        "extra_buffers",
			Version:     BuiltIn,
			Description: "allows declaration of one or more extra goal-style buffers in the model",
		},
	}
}

// SetParam is called to set our module's parameter from the parameter in the code ("param")
func (eb *ExtraBuffers) SetParam(param *params.Param) (err error) {
	newBuffer := buffer.Buffer{
		Name:         param.Key,
		MultipleInit: false,
	}

	eb.BufferList = append(eb.BufferList, newBuffer)
	return
}
