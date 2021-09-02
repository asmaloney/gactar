package actr

type Memory struct {
	Name   string
	Buffer BufferInterface

	// The following optional fields came from the ccm framework.
	// TODO: determine if they apply to others.
	Latency   *float64
	Threshold *float64
	MaxTime   *float64
	FinstSize *int // finst == "fingers of instantiation"
	FinstTime *float64
}
