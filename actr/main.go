package actr

// Model represents a basic ACT-R model.
// This is used as input to a Framework where it can be run or output to a file.
// (This is incomplete w.r.t. all of ACT-R's capabilities.)
type Model struct {
	Name         string
	Description  string
	Examples     []string
	Buffers      []*Buffer
	Memories     []*Memory
	TextOutputs  []*TextOutput
	Initializers []*Initializer
	Productions  []*Production
	Logging      bool
}

type Buffer struct {
	Name string
}

type Memory struct {
	Name      string
	Buffer    *Buffer // required
	Latency   *float64
	Threshold *float64
	MaxTime   *float64
	FinstSize *int     // not sure what the 'f' is in finst?
	FinstTime *float64 // not sure what the 'f' is in finst?
}

type TextOutput struct {
	Name string
}

type Initializer struct {
	MemoryName string
	Text       string
}

// LookupBuffer looks up the named buffer in the model and returns it (or nil if it does not exist).
func (model *Model) LookupBuffer(bufferName string) *Buffer {
	for _, buf := range model.Buffers {
		if buf.Name == bufferName {
			return buf
		}
	}

	return nil
}

// LookupMemory looks up the named memory in the model and returns it (or nil if it does not exist).
func (model *Model) LookupMemory(memoryName string) *Memory {
	for _, mem := range model.Memories {
		if mem.Name == memoryName {
			return mem
		}
	}

	return nil
}

// LookupTextOutput looks up the named text output in the model and returns it (or nil if it does not exist).
func (model *Model) LookupTextOutput(textOutputName string) *TextOutput {
	for _, textOutput := range model.TextOutputs {
		if textOutput.Name == textOutputName {
			return textOutput
		}
	}

	return nil
}

// BufferOrMemoryExists looks up the named buffer or memory in the model and returns it (or nil if it does not exist).
func (model *Model) BufferOrMemoryExists(name string) bool {
	buffer := model.LookupBuffer(name)
	if buffer != nil {
		return true
	}

	memory := model.LookupMemory(name)

	return memory != nil
}
