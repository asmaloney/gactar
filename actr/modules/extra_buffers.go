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
			Name: "extra_buffers",
		},
	}
}

func (eb *ExtraBuffers) SetParam(param *params.Param) (err error) {
	newBuffer := buffer.Buffer{
		Name:         param.Key,
		MultipleInit: false,
	}

	eb.Buffers = append(eb.Buffers, newBuffer)
	return
}
