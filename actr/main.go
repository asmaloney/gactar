package actr

// Model represents a basic ACT-R model
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

type Production struct {
	Name    string
	Matches []*Match
	Do      []string
}

type Match struct {
	Name string
	Text string
}

// LookupBuffer looks up the named buffer in the model and returns it (or nil if it does not exist)
func (model *Model) LookupBuffer(bufferName string) (buffer *Buffer) {
	for _, buf := range model.Buffers {
		if buf.Name == bufferName {
			return buf
		}
	}

	return
}

// LookupMemory looks up the named memory in the model and returns it (or nil if it does not exist)
func (model *Model) LookupMemory(memoryName string) (memory *Memory) {
	for _, mem := range model.Memories {
		if mem.Name == memoryName {
			return mem
		}
	}

	return
}

// LookupTextOutput looks up the named text output in the model and returns it (or nil if it does not exist)
func (model *Model) LookupTextOutput(textOutputName string) (textOutput *TextOutput) {
	for _, textOutput := range model.TextOutputs {
		if textOutput.Name == textOutputName {
			return textOutput
		}
	}

	return
}

// BufferOrMemoryExists looks up the named item in the model
func (model *Model) BufferOrMemoryExists(name string) bool {
	buffer := model.LookupBuffer(name)
	if buffer != nil {
		return true
	}

	memory := model.LookupMemory(name)

	return memory != nil
}
