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

	// "partial_matching": turn on partial matching of buffers in productions
	// (there are no defaults since setting it activates the capability)
	// 	ccm (Partial class)
	// 	pyactr (partial_matching)
	// 	vanilla (:ppm)
	PartialMatching bool
}

func NewProcedural() *Procedural {
	return &Procedural{
		Module: Module{Name: "procedural"},
	}
}

func (p *Procedural) SetParam(param *params.Param) (err params.ParamError) {
	value := param.Value

	switch param.Key {
	case "default_action_time":
		if value.Number == nil {
			return params.NumberRequired
		}

		if *value.Number < 0 {
			return params.NumberMustBePositive
		}

		p.DefaultActionTime = value.Number

	case "partial_matching":
		boolVal, err := value.AsBool()
		if err != params.NoError {
			return err
		}

		p.PartialMatching = boolVal

	default:
		return params.UnrecognizedParam
	}

	return
}
