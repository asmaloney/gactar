// Package actr implements our internal, parsed version of the amod file which is passed
// to a Framework to generate their code.
package actr

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/modules"
	"github.com/asmaloney/gactar/actr/params"
	"github.com/asmaloney/gactar/util/container"
)

type options struct {
	// "log_level": one of 'min', 'info', or 'detail'
	LogLevel ACTRLogLevel

	// "trace_activations": output detailed info about activations
	TraceActivations bool
}

// Model represents a basic ACT-R model.
// This is used as input to a Framework where it can be run or output to a file.
// (This is incomplete w.r.t. all of ACT-R's capabilities.)
type Model struct {
	Name        string
	Description string
	Authors     []string
	Examples    []*Pattern

	Modules    []modules.ModuleInterface
	Memory     *modules.DeclarativeMemory // memory is always present
	Goal       *modules.Goal              // goal is always present
	Procedural *modules.Procedural        // procedural is always present
	Chunks     []*Chunk

	// ExplicitChunks are ones which are named in the initializer like this:
	// 	castle [meaning: 'castle']
	// We keep track of these to remove them from the implicit list & to check for duplicates.
	ExplicitChunks []string

	// ImplicitChunks are chunks which aren't declared, but need to be created by some frameworks.
	// e.g. by default vanilla will create them and emit a warning:
	// 	#|Warning: Creating chunk SHARK with no slots |#
	// These chunk names come from the initializations.
	// We keep track of them so we can create them explicitly to avoid the warnings.
	ImplicitChunks []string

	Initializers []*Initializer
	Productions  []*Production

	options
}

type Initializer struct {
	Module modules.ModuleInterface
	Buffer buffer.BufferInterface

	ChunkName *string // optional chunk name
	Pattern   *Pattern

	AMODLineNumber int
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

func (model *Model) AddImplicitChunk(chunkName string) {
	model.ImplicitChunks = append(model.ImplicitChunks, chunkName)
}

// AddImplicitChunksFromPattern walks a pattern and adds any IDs to our list of implicit chunks.
func (model *Model) AddImplicitChunksFromPattern(pattern *Pattern) {
	if pattern == nil {
		return
	}

	for _, slot := range pattern.Slots {
		if slot.ID != nil {
			model.AddImplicitChunk(*slot.ID)
		}
	}
}

// AddInitializer adds the initializer to our list and adds any IDs to ImplicitChunks
// so we can (possibly) create them in the framework output.
func (model *Model) AddInitializer(initializer *Initializer) {
	model.Initializers = append(model.Initializers, initializer)

	if initializer.ChunkName != nil {
		model.ExplicitChunks = append(model.ExplicitChunks, *initializer.ChunkName)
	}

	model.AddImplicitChunksFromPattern(initializer.Pattern)
}

// LookupInitializer returns an initializer or nil if the buffer does not have one.
func (model Model) LookupInitializer(buffer string) *Initializer {
	for _, init := range model.Initializers {
		if init.Module.HasBuffer(buffer) {
			return init
		}
	}

	return nil
}

func (model Model) HasImplicitChunks() bool {
	return len(model.ImplicitChunks) > 0
}

// FinalizeImplicitChunks will remove already-declared chunks from the implicit list and
// make the list unique.
func (model *Model) FinalizeImplicitChunks() {
	if !model.HasImplicitChunks() {
		return
	}

	list := container.UniqueAndSorted(model.ImplicitChunks)

	chunkNameList := []string{}
	for _, chunk := range model.Chunks {
		chunkNameList = append(chunkNameList, chunk.Name)
	}

	chunkNameList = append(chunkNameList, model.ExplicitChunks...)

	for _, chunkName := range chunkNameList {
		list = container.FindAndDelete(list, chunkName)
	}

	model.ImplicitChunks = list
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

// CreateExtraBuffers creates the "extra_buffers" module and adds it to the list.
func (model *Model) CreateExtraBuffers() *modules.ExtraBuffers {
	eb := modules.NewExtraBuffers()
	model.Modules = append(model.Modules, eb)
	return eb
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
		names := module.BufferNames()
		if len(names) > 0 {
			list = append(list, names...)
		}
	}

	return
}

// LookupBuffer looks up the named buffer in the model and returns it (or nil if it does not exist).
func (model Model) LookupBuffer(bufferName string) buffer.BufferInterface {
	for _, module := range model.Modules {
		buff := module.LookupBuffer(bufferName)
		if buff != nil {
			return buff
		}
	}

	return nil
}

func (model *Model) SetParam(param *params.Param) (err error) {
	value := param.Value

	switch param.Key {
	case "log_level":
		if (value.Str == nil) || !ValidLogLevel(*value.Str) {
			return params.ErrInvalidOption{Expected: ACTRLoggingLevels}
		}

		model.LogLevel = ACTRLogLevel(*value.Str)

	case "trace_activations":
		boolVal, err := value.AsBool()
		if err != nil {
			return err
		}

		model.TraceActivations = boolVal

	default:
		return params.ErrUnrecognizedParam
	}

	return
}
