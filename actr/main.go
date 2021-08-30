package actr

import (
	"fmt"
	"strings"
)

// User cannot create chunks with these names. Perhaps needs to be expanded with other keywords?
var reservedChunkNames = map[string]bool{
	"_status":   true,
	"goal":      true,
	"memory":    true,
	"retrieval": true,
}

// Model represents a basic ACT-R model.
// This is used as input to a Framework where it can be run or output to a file.
// (This is incomplete w.r.t. all of ACT-R's capabilities.)
type Model struct {
	Name         string
	Description  string
	Examples     []*Pattern
	Chunks       []*Chunk
	Buffers      []BufferInterface
	Memories     []*Memory // we only have one memory now, but leave as slice until we determine if we can have multiple memories
	Initializers []*Initializer
	Productions  []*Production
	Logging      bool
}

type Chunk struct {
	Name      string
	SlotNames []string
	NumSlots  int
}

type BufferInterface interface {
	GetName() string
}

type Buffer struct {
	Name string
}

func (b Buffer) GetName() string {
	return b.Name
}

func (b Buffer) String() string {
	return b.Name
}

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

type TextOutput struct {
	Name string
}

type Initializer struct {
	Buffer BufferInterface // buffer...
	Memory *Memory         // ... OR memory

	Pattern *Pattern
}

func (c Chunk) IsInternal() bool {
	return c.Name[0] == '_'
}

func IsInternalChunkName(name string) bool {
	return name[0] == '_'
}

func ReservedChunkNameExists(name string) bool {
	v, ok := reservedChunkNames[name]
	return v && ok
}

func (c Chunk) String() (str string) {
	return fmt.Sprintf("%s( %s )", c.Name, strings.Join(c.SlotNames, " "))
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

	model.Memories = []*Memory{
		{
			Name:   "memory",
			Buffer: retrieval,
		},
	}
}

// LookupChunk looks up the named chunk in the model and returns it (or nil if it does not exist).
func (model Model) LookupChunk(chunkName string) *Chunk {
	for _, chunk := range model.Chunks {
		if chunk.Name == chunkName {
			return chunk
		}
	}

	return nil
}

// LookupBuffer looks up the named buffer in the model and returns it (or nil if it does not exist).
func (model Model) LookupBuffer(bufferName string) BufferInterface {
	for _, buf := range model.Buffers {
		if buf.GetName() == bufferName {
			return buf
		}
	}

	return nil
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

// HasSlot checks if the slot name exists on this chunk.
func (chunk Chunk) HasSlot(slot string) bool {
	for _, name := range chunk.SlotNames {
		if name == slot {
			return true
		}
	}

	return false
}

// GetSlotIndex returns the slot index (indexed from 1) of the slot name or -1 if not found.
func (chunk Chunk) GetSlotIndex(slot string) int {
	for i, name := range chunk.SlotNames {
		if name == slot {
			return i + 1
		}
	}

	return -1
}
