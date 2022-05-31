package actr

// We need to take apart patterns such as [add: ?num1 ?num2 ?count!?num2 ?sum]
// so we can verify variable use.

type Pattern struct {
	Chunk *Chunk
	Slots []*PatternSlot
}

type PatternSlot struct {
	Items []*PatternSlotItem
}

type PatternSlotItem struct {
	// The item is one of the following:
	Nil      bool
	Wildcard bool
	ID       *string
	Var      *string
	Num      *string // we don't need to treat this as a number anywhere, so keep as a string

	Negated bool // this item is negated
}

func (p PatternSlot) String() (str string) {
	for _, item := range p.Items {
		if item.Negated {
			str += "!"
		}

		if item.Wildcard {
			str += "*"
		} else if item.Nil {
			str += "nil"
		} else if item.ID != nil {
			str += *item.ID
		} else if item.Var != nil {
			str += *item.Var
		} else if item.Num != nil {
			str += *item.Num
		}
	}

	return
}

func (p *PatternSlot) AddItem(item *PatternSlotItem) {
	p.Items = append(p.Items, item)
}

func (p Pattern) String() (str string) {
	str = "[" + p.Chunk.Name + ": "

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

func (p Pattern) LookupVariable(varName string) *PatternSlotItem {
	for _, slot := range p.Slots {
		for _, item := range slot.Items {
			if item.Var == nil {
				continue
			}

			if *item.Var == varName {
				return item
			}
		}
	}

	return nil
}
