package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

type Procedural struct {
	buffer.BufferInterface // unused

	// "default_action_time": time that it takes to fire a production (seconds)
	// ccm: 0.05
	// pyactr: 0.05
	// vanilla: 0.05
	DefaultActionTime *float64
}

func NewProcedural() *Procedural {
	return &Procedural{BufferInterface: buffer.Buffer{}}
}

func (Procedural) ModuleName() string {
	return "procedural"
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

	default:
		return params.UnrecognizedParam
	}

	return
}
