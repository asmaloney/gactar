package actr

import "fmt"

type Pattern struct {
	AnyChunk bool

	Chunk *Chunk
	Slots []*PatternSlot
}

type PatternVar struct {
	Name        *string
	Constraints []*Constraint // any constraints on this var
}

type PatternSlot struct {
	// The item is one of the following:
	Nil      bool
	Wildcard bool

	ID  *string
	Str *string
	Var *PatternVar
	Num *string // we don't need to treat this as a number anywhere, so keep as a string

	Negated bool // this item is negated
}

func (p PatternSlot) String() (str string) {
	if p.Negated {
		str += "!"
	}

	switch {
	case p.Wildcard:
		str += "*"

	case p.Nil:
		str += "nil"

	case p.ID != nil:
		str += *p.ID

	case p.Str != nil:
		str += fmt.Sprintf("'%s'", *p.Str)

	case p.Var != nil:
		str += *p.Var.Name

	case p.Num != nil:
		str += *p.Num
	}

	return
}

func (p Pattern) String() (str string) {
	str = "[" + p.Chunk.TypeName + ": "

	numSlots := len(p.Slots)

	for i, slot := range p.Slots {
		str += slot.String()
		if i < numSlots-1 {
			str += " "
		}
	}

	str += "]"

	return
}

func (p *Pattern) AddSlot(slot *PatternSlot) {
	p.Slots = append(p.Slots, slot)
}

func (p Pattern) LookupVariable(varName string) *PatternVar {
	for _, slot := range p.Slots {
		if slot.Var == nil {
			continue
		}

		if *slot.Var.Name == varName {
			return slot.Var
		}
	}

	return nil
}
