package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

// Goal is a module which provides the ACT-R "goal" buffer.
type Goal struct {
	buffer.BufferInterface

	// "spreading_activation": see "Spreading Activation" in "ACT-R 7.26 Reference Manual" pg. 290
	// ccm: 1.0
	// pyactr: 1.0
	// vanilla: 1.0
	SpreadingActivation *float64
}

func NewGoal() *Goal {
	return &Goal{
		BufferInterface: &buffer.Buffer{Name: "goal", MultipleInit: false},
	}
}

func (g Goal) ModuleName() string {
	return "goal"
}

func (g *Goal) SetParam(param *params.Param) (err params.ParamError) {
	value := param.Value

	switch param.Key {
	case "spreading_activation":
		if value.Number == nil {
			return params.NumberRequired
		}

		if *value.Number < 0 {
			return params.NumberMustBePositive
		}

		g.SpreadingActivation = value.Number

	default:
		return params.UnrecognizedParam
	}

	return
}
