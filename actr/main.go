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
	Name         string
	Matches      []*Match
	DoPython     []string
	DoStatements []*Statement
}

type Match struct {
	Name string
	Text string
}

type Statement struct {
	Print  *PrintStatement
	Recall *RecallStatement
	Set    *SetStatement
	Write  *WriteStatement
}

type PrintStatement struct {
	Args []string // the strings, identifiers, or numbers to print
}

type RecallStatement struct {
	Contents   string
	MemoryName string
}

type WriteStatement struct {
	Args           []string // the strings, identifiers, or numbers to write
	TextOutputName string
}

type ArgOrField struct {
	ArgNum    *int
	FieldName *string
}

type SetStatement struct {
	ArgOrField *ArgOrField
	BufferName string
	Contents   string
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
