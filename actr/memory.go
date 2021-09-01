package actr

type Memory struct {
	Name   string
	Buffer BufferInterface

	// The following optional fields came from the ccm framework.
	// TODO: determine if they apply to others.
	Latency   *float64
	Threshold *float64
	MaxTime   *float64
	FinstSize *int     // not sure what the 'f' is in finst?
	FinstTime *float64 // not sure what the 'f' is in finst?
}

// LookupMemory looks up the named memory in the model and returns it (or nil if it does not exist).
func (model Model) LookupMemory(memoryName string) *Memory {
	for _, mem := range model.Memories {
		if mem.Name == memoryName {
			return mem
		}
	}

	return nil
}
