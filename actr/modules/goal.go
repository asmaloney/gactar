package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

// Goal is a module which provides the ACT-R "goal" buffer.
type Goal struct {
	Module

	// "spreading_activation": see "Spreading Activation" in "ACT-R 7.26 Reference Manual" pg. 290
	// 	ccm (DMSpreading.weight): 1.0
	// 	pyactr (buffer_spreading_activation): 1.0
	// 	vanilla (:ga): 1.0
	SpreadingActivation *float64
}

func NewGoal() *Goal {
	return &Goal{
		Module: Module{
			Name: "goal",
			Buffers: []buffer.Buffer{
				{Name: "goal", MultipleInit: false},
			},
		},
	}
}

func (g *Goal) SetParam(param *params.Param) (err error) {
	value := param.Value

	switch param.Key {
	case "spreading_activation":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		if *value.Number < 0 {
			return params.ErrMustBePositive
		}

		g.SpreadingActivation = value.Number

	default:
		return params.ErrUnrecognizedParam
	}

	return
}
