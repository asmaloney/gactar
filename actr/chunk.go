package actr

import "github.com/asmaloney/gactar/util/container"

type Chunk struct {
	TypeName  string
	SlotNames []string
	NumSlots  int

	AMODLineNumber int // line number in the amod file of the this chunk declaration
}

func IsInternalChunkType(name string) bool {
	return name[0] == '_'
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
	return container.Contains(slot, chunk.SlotNames)
}

// SlotIndex returns the slot index (indexed from 1) of the slot name or -1 if not found.
func (chunk Chunk) SlotIndex(slot string) int {
	return container.GetIndex1(slot, chunk.SlotNames)
}
