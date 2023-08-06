package modules

import (
	"github.com/asmaloney/gactar/actr/param"
	"github.com/asmaloney/gactar/util/keyvalue"
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
	defActionTime := param.NewFloat(
		"default_action_time",
		"time that it takes to fire a production (seconds)",
		param.Ptr(0.0), nil,
	)

	parameters := param.NewParameters(param.List{
		defActionTime,
	})

	return &Procedural{
		Module: Module{
			Name:                "procedural",
			Version:             BuiltIn,
			Description:         "handles production definition and execution",
			ParametersInterface: parameters,
		},
	}
}

// SetParam is called to set our module's parameter from the parameter in the code ("param")
func (p *Procedural) SetParam(param *keyvalue.KeyValue) (err error) {
	err = p.ValidateParam(param)
	if err != nil {
		return
	}

	value := param.Value

	if param.Key == "default_action_time" {
		p.DefaultActionTime = value.Number
	}

	return
}
