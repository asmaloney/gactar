package actr

import (
	"fmt"
	"strings"
)

// User cannot create chunks with these names. Perhaps needs to be expanded with other keywords?
var reservedChunkNames = map[string]bool{
	"_status":   true,
	"goal":      true,
	"imaginal":  true,
	"memory":    true,
	"retrieval": true,
}

type Chunk struct {
	Name      string
	SlotNames []string
	NumSlots  int
}

func IsInternalChunkName(name string) bool {
	return name[0] == '_'
}

func ReservedChunkNameExists(name string) bool {
	v, ok := reservedChunkNames[name]
	return v && ok
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

func (c Chunk) String() (str string) {
	return fmt.Sprintf("%s( %s )", c.Name, strings.Join(c.SlotNames, " "))
}

func (c Chunk) IsInternal() bool {
	return c.Name[0] == '_'
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