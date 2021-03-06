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

func NewProcedural() *Procedural {
	return &Procedural{
		Module: Module{Name: "procedural"},
	}
}

func (p *Procedural) SetParam(param *params.Param) (err error) {
	value := param.Value

	switch param.Key {
	case "default_action_time":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		if *value.Number < 0 {
			return params.ErrMustBePositive
		}

		p.DefaultActionTime = value.Number

	default:
		return params.ErrUnrecognizedParam
	}

	return
}
