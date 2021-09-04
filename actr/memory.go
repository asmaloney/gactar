package actr

type Memory struct {
	Name   string
	Buffer BufferInterface

	// Setting these kinds of parameters is going to be tricky.
	// The way each framework handles timing is different.
	// For example pyactr has an additional "rule_firing" setting that affects timing (see
	// "Computational Cognitive Modeling and Linguistic Theory" section 8.4 pg. 201).

	Latency   *float64 // latency factor
	Threshold *float64 // retrieval threshold
	MaxTime   *float64 // maximum time to run (in sim time, not real-time)
	FinstSize *int     // finst == "fingers of instantiation"
	FinstTime *float64
}
