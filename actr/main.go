// Package actr implements our internal, parsed version of the amod file which is passed
// to a Framework to generate their code.
package actr

// Model represents a basic ACT-R model.
// This is used as input to a Framework where it can be run or output to a file.
// (This is incomplete w.r.t. all of ACT-R's capabilities.)
type Model struct {
	Name         string
	Description  string
	Authors      []string
	Examples     []*Pattern
	Chunks       []*Chunk
	Modules      []ModuleInterface
	Memory       *DeclMemory // memory is always present, so keep track of it instead of looking it up
	Initializers []*Initializer
	Productions  []*Production
	LogLevel     ACTRLogLevel
}

type Initializer struct {
	Buffer         BufferInterface
	Pattern        *Pattern
	AMODLineNumber int // line number in the amod file of this initialization
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

	model.Modules = append(model.Modules, NewGoal())

	model.Memory = NewDeclMemory()
	model.Modules = append(model.Modules, model.Memory)

	model.LogLevel = "info"
}

// LookupInitializer returns an initializer or nil if the buffer does not have one.
func (model Model) LookupInitializer(buffer string) *Initializer {
	for _, init := range model.Initializers {
		if init.Buffer.GetBufferName() == buffer {
			return init
		}
	}

	return nil
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

// CreateImaginal creates the imaginal module and adds it to the list.
func (model *Model) CreateImaginal() *Imaginal {
	imaginal := NewImaginal()
	model.Modules = append(model.Modules, imaginal)
	return imaginal
}

// GetImaginal gets the imaginal module (or returns nil if it does not exist).
func (model Model) GetImaginal() *Imaginal {
	module := model.LookupModule("imaginal")
	if module == nil {
		return nil
	}

	imaginal, ok := module.(*Imaginal)
	if !ok {
		return nil
	}

	return imaginal
}

// LookupModule looks up the named module in the model and returns it (or nil if it does not exist).
func (model Model) LookupModule(moduleName string) ModuleInterface {
	for _, module := range model.Modules {
		if module.GetModuleName() == moduleName {
			return module
		}
	}

	return nil
}

// LookupBuffer looks up the named buffer in the model and returns it (or nil if it does not exist).
func (model Model) LookupBuffer(bufferName string) BufferInterface {
	for _, module := range model.Modules {
		if module.GetBufferName() == bufferName {
			return module
		}
	}

	return nil
}
