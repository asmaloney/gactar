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

	// "mismatch_penalty": turns on partial matching and sets the penalty in the activation equation to this
	// (there are no defaults since setting it activates the capability)
	// 	ccm (Partial class)
	// 	pyactr (partial_matching & mismatch_penalty)
	// 	vanilla (:mp)
	MismatchPenalty *float64
}

// NewDeclarativeMemory creates and returns a new DeclarativeMemory module
func NewDeclarativeMemory() *DeclarativeMemory {
	decay := NewParamFloat(
		"decay",
		"the 'base-level learning' decay parameter",
		Ptr(0.0), Ptr(1.0),
	)

	finstSize := NewParamInt(
		"finst_size",
		"how many chunks are retained in memory",
		Ptr(0), nil,
	)

	finstTime := NewParamFloat(
		"finst_time",
		"how long the finst lasts in memory",
		Ptr(0.0), nil,
	)

	instNoise := NewParamFloat(
		"instantaneous_noise",
		"turns on the activation noise calculation & sets instantaneous noise",
		Ptr(0.0), nil,
	)

	latencyExponent := NewParamFloat(
		"latency_exponent",
		"latency exponent (f)",
		Ptr(0.0), nil,
	)

	latencyFactor := NewParamFloat(
		"latency_factor",
		"latency latency_factor (F)",
		Ptr(0.0), nil,
	)

	maxSpreadStrength := NewParamFloat(
		"max_spread_strength",
		"turns on the spreading activation calculation & sets the maximum associative strength",
		nil, nil,
	)

	mismatchPenalty := NewParamFloat(
		"mismatch_penalty",
		"turns on partial matching and sets the penalty in the activation equation",
		nil, nil,
	)

	retrievalThreshold := NewParamFloat(
		"retrieval_threshold",
		"retrieval threshold (τ)",
		nil, nil,
	)

	return &DeclarativeMemory{
		Module: Module{
			Name:        "memory",
			Version:     BuiltIn,
			Description: "declarative memory which stores chunks from the buffers for retrieval",
			BufferList: buffer.List{
				{Name: "retrieval", MultipleInit: true},
			},
			Params: ParamInfoMap{
				"decay":               decay,
				"finst_size":          finstSize,
				"finst_time":          finstTime,
				"instantaneous_noise": instNoise,
				"latency_exponent":    latencyExponent,
				"latency_factor":      latencyFactor,
				"max_spread_strength": maxSpreadStrength,
				"mismatch_penalty":    mismatchPenalty,
				"retrieval_threshold": retrievalThreshold,
			},
		},
	}
}

// We only have one buffer, so this is a convenience function to get its name.
func (d DeclarativeMemory) BufferName() string {
	return d.Buffers().At(0).BufferName()
}

// SetParam is called to set our module's parameter from the parameter in the code ("param")
func (d *DeclarativeMemory) SetParam(param *params.Param) (err error) {
	err = d.ValidateParam(param)
	if err != nil {
		return
	}

	value := param.Value

	switch param.Key {
	case "latency_factor":
		d.LatencyFactor = value.Number

	case "latency_exponent":
		d.LatencyExponent = value.Number

	case "retrieval_threshold":
		d.RetrievalThreshold = value.Number

	case "finst_size":
		size := int(*value.Number)
		d.FinstSize = &size

	case "finst_time":
		d.FinstTime = value.Number

	case "decay":
		d.Decay = value.Number

	case "max_spread_strength":
		d.MaxSpreadStrength = value.Number

	case "instantaneous_noise":
		d.InstantaneousNoise = value.Number

	case "mismatch_penalty":
		d.MismatchPenalty = value.Number
	}

	return
}
