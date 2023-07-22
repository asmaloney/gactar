package modules

import (
	"github.com/asmaloney/gactar/actr/params"
)

type Procedural struct {
	Module

	// "default_action_time": time that it takes to fire a production (seconds)
	// 	ccm (production_time): 0.05
	// 	pyactr (rule_firing): 0.05
	// 	vanilla (:dat): 0.05
	DefaultActionTime *float64
}

// NewProcedural creates and returns a new Procedural module
func NewProcedural() *Procedural {
	defActionTime := NewParamFloat(
		"default_action_time",
		"time that it takes to fire a production (seconds)",
		Ptr(0.0), nil,
	)

	return &Procedural{
		Module: Module{
			Name:        "procedural",
			Version:     BuiltIn,
			Description: "handles production definition and execution",
			Params: ParamInfoMap{
				"default_action_time": defActionTime,
			},
		},
	}
}

// SetParam is called to set our module's parameter from the parameter in the code ("param")
func (p *Procedural) SetParam(param *params.Param) (err error) {
	err = p.ValidateParam(param)
	if err != nil {
		return
	}

	value := param.Value

	switch param.Key {
	case "default_action_time":
		p.DefaultActionTime = value.Number
	}

	return
}
