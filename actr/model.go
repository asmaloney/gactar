// Package actr implements our internal, parsed version of the amod file which is passed
// to a Framework to generate their code.
package actr

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/modules"
	"github.com/asmaloney/gactar/actr/param"

	"github.com/asmaloney/gactar/util/container"
	"github.com/asmaloney/gactar/util/keyvalue"
)

type Options struct {
	// "log_level": one of 'min', 'info', or 'detail'
	LogLevel ACTRLogLevel

	// "trace_activations": output detailed info about activations
	TraceActivations bool

	// "random_seed": the seed to use for generating pseudo-random numbers (allows for reproducible runs)
	// For all frameworks, if it is not set it uses current system time.
	// Use a uint32 because pyactr uses numpy and that's what its random number seed uses.
	RandomSeed *uint32
}

// Model represents a basic ACT-R model.
// This is used as input to a Framework where it can be run or output to a file.
// (This is incomplete w.r.t. all of ACT-R's capabilities.)
type Model struct {
	Name        string
	Description string
	Authors     []string
	Examples    []*Pattern

	Modules    []modules.Interface
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
	// These chunk names come from the initializations & similarities.
	// We keep track of them so we can create them explicitly to avoid the warnings.
	ImplicitChunks []string

	Initializers []*Initializer
	Similarities []*Similarity

	Productions []*Production

	Options

	// Used to validate our parameters
	parameters param.ParametersInterface
}

type Initializer struct {
	Module modules.Interface
	Buffer buffer.Interface

	ChunkName *string // optional chunk name
	Pattern   *Pattern

	AMODLineNumber int
}

type Similarity struct {
	ChunkOne string
	ChunkTwo string
	Value    float64

	AMODLineNumber int
}

func (model *Model) Initialize() {
	// Set up our built-in modules

	model.Memory = modules.NewDeclarativeMemory()
	model.Modules = append(model.Modules, model.Memory)

	model.Goal = modules.NewGoal()
	model.Modules = append(model.Modules, model.Goal)

	model.Procedural = modules.NewProcedural()
	model.Modules = append(model.Modules, model.Procedural)

	model.LogLevel = "info"

	// Declare our parameters
	loggingParam := param.NewStr(
		"log_level",
		"Level of logging output",
		ACTRLoggingLevels,
	)

	traceParam := param.NewBool(
		"trace_activations",
		"output detailed info about activations",
	)

	seedParam := param.NewInt(
		"random_seed",
		"the seed to use for generating pseudo-random numbers",
		param.Ptr(0), nil,
	)

	parameters := param.NewParameters(param.InfoMap{
		"log_level":         loggingParam,
		"trace_activations": traceParam,
		"random_seed":       seedParam,
	})

	model.parameters = parameters
}

func (m *Model) SetRunOptions(options *Options) {
	if options == nil {
		return
	}

	m.LogLevel = options.LogLevel
	m.TraceActivations = options.TraceActivations

	if options.RandomSeed != nil {
		m.RandomSeed = options.RandomSeed
	}
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
		if init.Module.Buffers().Has(buffer) {
			return init
		}
	}

	return nil
}

// AddSimilarity will add a similarity to the list and keep track of the chunk names.
func (model *Model) AddSimilarity(similar *Similarity) {
	model.Similarities = append(model.Similarities, similar)

	model.ImplicitChunks = append(model.ImplicitChunks, similar.ChunkOne, similar.ChunkTwo)
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

	for _, chunkName := range model.ExplicitChunks {
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
func (model Model) LookupModule(moduleName string) modules.Interface {
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
		names := module.Buffers().Names()
		if len(names) > 0 {
			list = append(list, names...)
		}
	}

	return
}

// LookupBuffer looks up the named buffer in the model and returns it (or nil if it does not exist).
func (model Model) LookupBuffer(bufferName string) buffer.Interface {
	for _, module := range model.Modules {
		buff := module.Buffers().Lookup(bufferName)
		if buff != nil {
			return buff
		}
	}

	return nil
}

func (model *Model) SetParam(kv *keyvalue.KeyValue) (err error) {
	err = model.parameters.ValidateParam(kv)
	if err != nil {
		return
	}

	value := kv.Value

	switch kv.Key {
	case "log_level":
		model.LogLevel = ACTRLogLevel(*value.Str)

	case "trace_activations":
		boolVal, _ := value.AsBool() // already validated
		model.TraceActivations = boolVal

	case "random_seed":
		seed := uint32(*value.Number)

		model.RandomSeed = &seed

	default:
		return param.ErrUnrecognizedOption{Option: kv.Key}
	}

	return
}
