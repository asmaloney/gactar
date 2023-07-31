package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

// Goal is a module which provides the ACT-R "goal" buffer.
type Goal struct {
	Module

	// "spreading_activation": see "Spreading Activation" in "ACT-R 7.26 Reference Manual" pg. 290
	// 	ccm (DMSpreading.weight): 0.0
	// 	pyactr (buffer_spreading_activation): 0.0
	// 	vanilla (:ga): 0.0
	SpreadingActivation *float64
}

// NewGoal creates and returns a new Goal module
func NewGoal() *Goal {
	spreadingActivation := NewParamFloat(
		"spreading_activation",
		"spreading activation weight",
		Ptr(0.0), nil,
	)

	return &Goal{
		Module: Module{
			Name:        "goal",
			Version:     BuiltIn,
			Description: "provides a goal buffer for the model",
			BufferList: buffer.List{
				{Name: "goal", MultipleInit: false},
			},
			Params: ParamInfoMap{
				"spreading_activation": spreadingActivation,
			},
		},
	}
}

// SetParam is called to set our module's parameter from the parameter in the code ("param")
func (g *Goal) SetParam(param *params.Param) (err error) {
	err = g.ValidateParam(param)
	if err != nil {
		return
	}

	value := param.Value

	switch param.Key {
	case "spreading_activation":
		g.SpreadingActivation = value.Number
	}

	return
}
