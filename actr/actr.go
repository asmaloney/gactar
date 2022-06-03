// Package actr implements our internal, parsed version of the amod file which is passed
// to a Framework to generate their code.
package actr

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/modules"
)

// Model represents a basic ACT-R model.
// This is used as input to a Framework where it can be run or output to a file.
// (This is incomplete w.r.t. all of ACT-R's capabilities.)
type Model struct {
	Name         string
	Description  string
	Authors      []string
	Examples     []*Pattern
	Chunks       []*Chunk
	Modules      []modules.ModuleInterface
	Memory       *modules.DeclarativeMemory // memory is always present
	Goal         *modules.Goal              // goal is always present
	Procedural   *modules.Procedural        // procedural is always present
	Initializers []*Initializer
	Productions  []*Production
	LogLevel     ACTRLogLevel
}

type Initializer struct {
	Buffer         buffer.BufferInterface
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

	// Set up our built-in modules

	model.Memory = modules.NewDeclarativeMemory()
	model.Modules = append(model.Modules, model.Memory)

	model.Goal = modules.NewGoal()
	model.Modules = append(model.Modules, model.Goal)

	model.Procedural = modules.NewProcedural()
	model.Modules = append(model.Modules, model.Procedural)

	model.LogLevel = "info"
}

// LookupInitializer returns an initializer or nil if the buffer does not have one.
func (model Model) LookupInitializer(buffer string) *Initializer {
	for _, init := range model.Initializers {
		if init.Buffer.BufferName() == buffer {
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
func (model *Model) CreateImaginal() *modules.Imaginal {
	imaginal := modules.NewImaginal()
	model.Modules = append(model.Modules, imaginal)
	return imaginal
}

// ImaginalModule gets the imaginal module (or returns nil if it does not exist).
func (model Model) ImaginalModule() *modules.Imaginal {
	module := model.LookupModule("imaginal")
	if module == nil {
		return nil
	}

	imaginal, ok := module.(*modules.Imaginal)
	if !ok {
		return nil
	}

	return imaginal
}

// LookupModule looks up the named module in the model and returns it (or nil if it does not exist).
func (model Model) LookupModule(moduleName string) modules.ModuleInterface {
	for _, module := range model.Modules {
		if module.ModuleName() == moduleName {
			return module
		}
	}

	return nil
}

// BufferNames returns a slice of valid buffers.
func (model Model) BufferNames() (list []string) {
	for _, module := range model.Modules {
		name := module.BufferName()
		if name != "" {
			list = append(list, name)
		}
	}

	return
}

// LookupBuffer looks up the named buffer in the model and returns it (or nil if it does not exist).
func (model Model) LookupBuffer(bufferName string) buffer.BufferInterface {
	for _, module := range model.Modules {
		if module.BufferName() == bufferName {
			return module
		}
	}

	return nil
}
