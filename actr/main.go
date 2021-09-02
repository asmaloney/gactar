package actr

// Model represents a basic ACT-R model.
// This is used as input to a Framework where it can be run or output to a file.
// (This is incomplete w.r.t. all of ACT-R's capabilities.)
type Model struct {
	Name         string
	Description  string
	Examples     []*Pattern
	Chunks       []*Chunk
	Buffers      []BufferInterface
	Memory       *Memory
	Initializers []*Initializer
	Productions  []*Production
	LogLevel     ACTRLogLevel
	HasImaginal  bool
}

type Initializer struct {
	Buffer BufferInterface // buffer...
	Memory *Memory         // ... OR memory

	Pattern *Pattern
}

func (model *Model) Initialize() {
	// Internal chunk for handling buffer and memory status
	model.Chunks = []*Chunk{
		{
			Name:      "_status",
			SlotNames: []string{"status"},
			NumSlots:  1,
		},
	}

	retrieval := &Buffer{Name: "retrieval"}
	model.Buffers = []BufferInterface{
		retrieval,
		&Buffer{Name: "goal"},
	}

	model.Memory = &Memory{
		Name:   "memory",
		Buffer: retrieval,
	}

	model.LogLevel = "info"
}

// HasInitializer checks if the model has an initialization for the buffer.
func (model Model) HasInitializer(buffer string) bool {
	for _, init := range model.Initializers {
		if init.Memory != nil {
			continue
		}

		if init.Buffer.GetName() == buffer {
			return true
		}
	}

	return false
}

// HasPrintStatement checks if this model uses the print statement.
// This is used to include extra code to handle printing in some frameworks.
func (model Model) HasPrintStatement() bool {
	for _, production := range model.Productions {
		if production.DoStatements != nil {
			for _, statement := range production.DoStatements {
				if statement.Print != nil {
					return true
				}
			}
		}
	}

	return false
}
