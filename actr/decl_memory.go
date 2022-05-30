package actr

// DeclMemory is a module which provides declarative memory.
type DeclMemory struct {
	BufferInterface

	// Setting these kinds of parameters is going to be tricky.
	// The way each framework handles timing is different.

	// Latency
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

func NewDeclMemory() *DeclMemory {
	return &DeclMemory{
		BufferInterface: Buffer{Name: "retrieval", MultipleInit: true},
	}
}

func (d DeclMemory) GetModuleName() string {
	return "memory"
}
