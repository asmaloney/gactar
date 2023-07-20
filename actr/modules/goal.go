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
	// 	pyactr (buffer_spreading_activation): no default
	// 	vanilla (:ga): 0.0
	SpreadingActivation *float64
}

func NewGoal() *Goal {
	return &Goal{
		Module: Module{
			Name:        "goal",
			Version:     BuiltIn,
			Description: "provides a goal buffer for the model",
			BufferList: buffer.List{
				{Name: "goal", MultipleInit: false},
			},
			Params: ParamInfoMap{
				"spreading_activation": {"spreading_activation", "spreading activation weight"},
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
