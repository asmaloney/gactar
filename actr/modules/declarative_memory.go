package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

type RetrievalTimeParams struct {
	// Setting these kinds of parameters is going to be tricky.
	// The way each framework handles timing seems to be different.

	// See "Retrieval time" in "ACT-R 7.26 Reference Manual" pg. 293

	// The time that it takes to respond to a request for a chunk is:
	//		RT = F * exp( - f * A )
	//		RT: The time to retrieve the chunk in seconds
	//		A : The activation of the chunk which is being retrieved

	// If there is no chunk found or the chunk with the highest activation
	// is below the retrieval threshold, then the time required to indicate
	// a failure to retrieve any chunk is:
	// 		RT = F * exp( - f * τ )
	// 		RT: The time until the failure is noted in seconds

	// "latency_factor": latency factor (F)
	// 	ccm (latency): 0.05
	// 	pyactr (latency_factor): 0.1
	// 	vanilla (:lf): 1.0
	LatencyFactor *float64

	// "latency_exponent": latency exponent (f)
	// 	ccm: (unsupported? Based on the formulae above and the code, it seems to be fixed at 1.0.)
	// 	pyactr (latency_exponent): 1.0
	// 	vanilla (:le): 1.0
	LatencyExponent *float64

	// "retrieval_threshold": retrieval threshold (τ)
	// 	ccm (threshold): 0.0
	// 	pyactr (retrieval_threshold): 0.0
	// 	vanilla (:rt): 0.0
	RetrievalThreshold *float64
}

type FinstParams struct {
	// finst ("fingers of instantiation")
	// See "Declarative finsts" in "ACT-R 7.26 Reference Manual" pg. 293

	// "finst_size": how many chunks are retained in memory
	//	ccm (finst_size): 4
	// 	pyactr (DecMemBuffer.finst): 0 (sic)
	// 	vanilla (:declarative-num-finsts): 4
	FinstSize *int

	// "finst_time": how long the finst lasts in memory
	// 	ccm (finst_time): 3.0
	// 	pyactr: (unsupported? Always ∞ I guess?)
	// 	vanilla (:declarative-finst-span): 3.0
	FinstTime *float64
}

// DeclarativeMemory is a module which provides declarative memory.
type DeclarativeMemory struct {
	Module

	RetrievalTimeParams
	FinstParams

	// "decay": sets the "base-level learning" decay parameter
	// 	ccm (DMBaseLevel submodule 'decay'): 0.5
	// 	pyactr (decay) : 0.5
	// 	vanilla (:bll): nil (recommend 0.5 if used)
	Decay *float64

	// "max_spread_strength": turns on the spreading activation calculation & sets the maximum associative strength
	// (there are no defaults since setting it activates the capability)
	//	ccm (DMSpreading submodule)
	//	pyactr (strength_of_association)
	//	vanilla (:mas)
	MaxSpreadStrength *float64

	// "instantaneous_noise": turns on the activation noise calculation & sets instantaneous noise
	// (there are no defaults since setting it activates the capability)
	// 	ccm (DMNoise submodule 'noise')
	// 	pyactr (instantaneous_noise)
	// 	vanilla (:ans)
	InstantaneousNoise *float64
}

func NewDeclarativeMemory() *DeclarativeMemory {
	return &DeclarativeMemory{
		Module: Module{
			Name: "memory",
			Buffers: []buffer.Buffer{
				{Name: "retrieval", MultipleInit: true},
			},
		},
	}
}

// We only have one buffer, so this is a convenience function to get its name.
func (d DeclarativeMemory) BufferName() string {
	return d.Buffers[0].Name
}

func (d *DeclarativeMemory) SetParam(param *params.Param) (err error) {
	value := param.Value

	switch param.Key {
	case "latency_factor":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		if *value.Number < 0 {
			return params.ErrMustBePositive
		}

		d.LatencyFactor = value.Number

	case "latency_exponent":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		if *value.Number < 0 {
			return params.ErrMustBePositive
		}

		d.LatencyExponent = value.Number

	case "retrieval_threshold":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		d.RetrievalThreshold = value.Number

	case "finst_size":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		size := int(*value.Number)
		if size < 0 {
			return params.ErrMustBePositive
		}

		d.FinstSize = &size

	case "finst_time":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		d.FinstTime = value.Number

	case "decay":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		if *value.Number < 0 || *value.Number > 1 {
			return params.ErrOutOfRange{
				Min: 0.0,
				Max: 1.0,
			}
		}

		d.Decay = value.Number

	case "max_spread_strength":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		d.MaxSpreadStrength = value.Number

	case "instantaneous_noise":
		if value.Number == nil {
			return params.ErrInvalidType{ExpectedType: params.Number}
		}

		if *value.Number < 0 {
			return params.ErrMustBePositive
		}

		d.InstantaneousNoise = value.Number

	default:
		return params.ErrUnrecognizedParam
	}

	return
}
