package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

// ExtraBuffers module is used to declare one or more extra goal-style buffers in the model.
type ExtraBuffers struct {
	Module
}

func NewExtraBuffers() *ExtraBuffers {
	return &ExtraBuffers{
		Module: Module{
			Name:        "extra_buffers",
			Version:     BuiltIn,
			Description: "allows declaration of one or more extra goal-style buffers in the model",
		},
	}
}

func (eb *ExtraBuffers) SetParam(param *params.Param) (err error) {
	newBuffer := buffer.Buffer{
		Name:         param.Key,
		MultipleInit: false,
	}

	eb.BufferList = append(eb.BufferList, newBuffer)
	return
}
