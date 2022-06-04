package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

// Imaginal is a module which provides the ACT-R "imaginal" buffer.
type Imaginal struct {
	buffer.BufferInterface

	// "delay": how long it takes a request to the buffer to complete (seconds)
	// ccm: 0.2
	// pyactr: 0.2
	// vanilla: 0.2
	Delay *float64
}

func NewImaginal() *Imaginal {
	return &Imaginal{
		BufferInterface: buffer.Buffer{Name: "imaginal", MultipleInit: false},
	}
}

func (i Imaginal) ModuleName() string {
	return "imaginal"
}

func (i *Imaginal) SetParam(param *params.Param) (err params.ParamError) {
	value := param.Value

	switch param.Key {
	case "delay":
		if value.Number == nil {
			return params.NumberRequired
		}

		if *value.Number < 0 {
			return params.NumberMustBePositive
		}

		i.Delay = value.Number

	default:
		return params.UnrecognizedParam
	}

	return
}
