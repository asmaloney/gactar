package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/param"

	"github.com/asmaloney/gactar/util/keyvalue"
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

// NewImaginal creates and returns a new Imaginal module
func NewImaginal() *Imaginal {
	delay := param.NewFloat(
		"delay",
		"time it takes a request to the buffer to complete (seconds)",
		param.Ptr(0.0), nil,
	)

	parameters := param.NewParameters(param.InfoMap{
		"delay": delay,
	})

	imaginalBuffer := buffer.NewBuffer("imaginal", 1.0, nil)

	return &Imaginal{
		Module: Module{
			Name:                "imaginal",
			Version:             BuiltIn,
			Description:         "provides a goal style buffer with a delay and an action buffer for manipulating the imaginal chunk",
			BufferList:          buffer.List{imaginalBuffer},
			ParametersInterface: parameters,
			MultipleInit:        false,
		},
	}
}

// SetParam is called to set our module's parameter from the parameter in the code ("param")
func (i *Imaginal) SetParam(param *keyvalue.KeyValue) (err error) {
	err = i.ValidateParam(param)
	if err != nil {
		return
	}

	value := param.Value

	if param.Key == "delay" {
		i.Delay = value.Number
	}

	return
}
