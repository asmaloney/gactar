package actr

// We need to take apart patterns such as 'add ?num1 ?num2 count:?count!?num2 sum:?sum'
// so we can verify variable use.

type Pattern struct {
	Chunk *Chunk
	Slots []*PatternSlot
}

type PatternSlot struct {
	Items []*PatternSlotItem
}

type PatternSlotItem struct {
	ID  *string
	Var *string
	Num *string // we don't need to treat this as a number anywhere, so keep as a string

	Negated bool // !
}

func (p PatternSlot) String() (str string) {
	for _, item := range p.Items {
		if item.Negated {
			str += "!"
		}

		if item.ID != nil {
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
	str = p.Chunk.Name + "( "

	for _, slot := range p.Slots {
		str += slot.String()
		str += " "
	}

	str += ")"

	return
}

func (p *Pattern) AddID(id *string) {
	slot := PatternSlot{}
	slot.Items = append(slot.Items, &PatternSlotItem{ID: id})

	p.Slots = append(p.Slots, &slot)
}

func (p *Pattern) AddVar(id *string, negated bool) {
	slot := PatternSlot{}
	slot.Items = append(slot.Items, &PatternSlotItem{Var: id,
		Negated: negated,
	})

	p.Slots = append(p.Slots, &slot)
}

func (p *Pattern) AddNum(num *string) {
	slot := PatternSlot{}
	slot.Items = append(slot.Items, &PatternSlotItem{Num: num})

	p.Slots = append(p.Slots, &slot)
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
