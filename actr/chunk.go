package actr

import (
	"slices"

	"github.com/asmaloney/gactar/util/container"
)

// See "Default Chunks" pg. 80 of ACT-R manual.
var reservedChunkNames = []string{"busy", "clear", "empty", "error", "failure", "free", "full", "requested", "unrequested"}

type Chunk struct {
	TypeName  string
	SlotNames []string
	NumSlots  int

	AMODLineNumber int // line number in the amod file of the this chunk declaration
}

func IsInternalChunkType(name string) bool {
	return name[0] == '_'
}

// IsReservedType checks if the slot name is reserved.
// See "Default Chunks" pg. 80 of ACT-R manual.
func IsReservedType(name string) bool {
	return slices.Contains(reservedChunkNames, name)
}

// LookupChunk looks up the chunk (by type name) in the model and returns it (or nil if it does not exist).
func (model Model) LookupChunk(typeName string) *Chunk {
	for _, chunk := range model.Chunks {
		if chunk.TypeName == typeName {
			return chunk
		}
	}

	return nil
}

// SlotName returns the name of the slot given the index.
func (c Chunk) SlotName(index int) (str string) {
	return c.SlotNames[index]
}

func (c Chunk) IsInternal() bool {
	return c.TypeName[0] == '_'
}

// HasSlot checks if the slot name exists on this chunk.
func (chunk Chunk) HasSlot(slot string) bool {
	return slices.Contains(chunk.SlotNames, slot)
}

// SlotIndex returns the slot index (indexed from 1) of the slot name or -1 if not found.
func (chunk Chunk) SlotIndex(slot string) int {
	return container.GetIndex1(slot, chunk.SlotNames)
}
