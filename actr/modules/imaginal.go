package modules

import "github.com/asmaloney/gactar/actr/buffer"

// Imaginal is a module which provides the ACT-R "imaginal" buffer.
type Imaginal struct {
	buffer.BufferInterface

	Delay float64 // non-negative time (in seconds) and defaults to .2
}

func NewImaginal() *Imaginal {
	// This uses the defaults as per ACT-R docs:
	// 	http://act-r.psy.cmu.edu/actr7.x/reference-manual.pdf page 276
	return &Imaginal{
		BufferInterface: buffer.Buffer{Name: "imaginal", MultipleInit: false},
		Delay:           0.2,
	}
}

func (i Imaginal) ModuleName() string {
	return "imaginal"
}

func (i *Imaginal) SetParam(param *Param) (err ParamError) {
	value := param.Value

	switch param.Key {
	case "delay":
		if value.Number == nil {
			return NumberRequired
		}

		if *value.Number < 0 {
			return NumberMustBePositive
		}

		i.Delay = *value.Number

	default:
		return UnrecognizedParam
	}

	return
}
