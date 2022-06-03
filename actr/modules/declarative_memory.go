package modules

import "github.com/asmaloney/gactar/actr/buffer"

type LatencyParams struct {
	// Setting these kinds of parameters is going to be tricky.
	// The way each framework handles timing is different.

	// See "Retrieval time" in "ACT-R 7.26 Reference Manual" pg. 293

	// "latency_factor": latency factor (F)
	// ccm: 0.05
	// pyactr: 0.1
	// vanilla: 1.0
	LatencyFactor *float64

	// "latency_exponent": latency exponent (f)
	// ccm: (unsupported?)
	// pyactr: 1.0
	// vanilla: 1.0
	LatencyExponent *float64

	// "retrieval_threshold": retrieval threshold (τ)
	// ccm: 0.0
	// pyactr: 0.0
	// vanilla: 0.0
	RetrievalThreshold *float64
}

type FinstParams struct {
	// finst ("fingers of instantiation")
	// See "Declarative finsts" in "ACT-R 7.26 Reference Manual" pg. 293

	// "finst_size": how many chunks are retained in memory
	// ccm: 4
	// pyactr: 0 (sic)
	// vanilla: 4
	FinstSize *int

	// "finst_time": how long the finst lasts in memory
	// ccm: 3.0
	// pyactr: (unsupported?)
	// vanilla: 3.0
	FinstTime *float64
}

// DeclarativeMemory is a module which provides declarative memory.
type DeclarativeMemory struct {
	buffer.BufferInterface

	LatencyParams
	FinstParams

	// "max_spread_strength": turns on the spreading activation calculation & sets the maximum associative strength
	// (there are no defaults since setting it activates the capability)
	MaxSpreadStrength *float64
}

func NewDeclarativeMemory() *DeclarativeMemory {
	return &DeclarativeMemory{
		BufferInterface: buffer.Buffer{Name: "retrieval", MultipleInit: true},
	}
}

func (d DeclarativeMemory) ModuleName() string {
	return "memory"
}
