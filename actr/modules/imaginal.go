package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

// Imaginal is a module which provides the ACT-R "imaginal" buffer.
type Imaginal struct {
	Module

	// "delay": how long it takes a request to the buffer to complete (seconds)
	// 	ccm (ImaginalModule.delay): 0.2
	// 	pyactr (Goal.delay): 0.2
	// 	vanilla (:imaginal-delay): 0.2
	Delay *float64
}

func NewImaginal() *Imaginal {
	delay := NewParamFloat(
		"delay",
		"time it takes a request to the buffer to complete (seconds)",
		Ptr(0.0), nil,
	)

	return &Imaginal{
		Module: Module{
			Name:        "imaginal",
			Version:     BuiltIn,
			Description: "provides a goal style buffer with a delay and an action buffer for manipulating the imaginal chunk",
			BufferList: buffer.List{
				{Name: "imaginal", MultipleInit: false},
			},
			Params: ParamInfoMap{
				"delay": delay,
			},
		},
	}
}

func (i *Imaginal) SetParam(param *params.Param) (err error) {
	err = i.ValidateParam(param)
	if err != nil {
		return
	}

	value := param.Value

	switch param.Key {
	case "delay":
		i.Delay = value.Number
	}

	return
}
