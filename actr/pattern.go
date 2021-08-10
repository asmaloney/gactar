package actr

// We need to take apart patterns such as 'add ?num1 ?num2 count:?count!?num2 sum:?sum'
// so we can verify variable use.

type Pattern struct {
	Fields []*PatternField
}

// PatternField allows for named fields - e.g. `foo:?bar!?blat`
type PatternField struct {
	Name  *string
	Items []*PatternFieldItem
}

type PatternFieldItem struct {
	ID  *string
	Var *string
	Num *string // we don't need to treat this as a number anywhere, so keep as a string

	Negated  bool // !
	Optional bool // ?
}

func (p PatternField) String() (str string) {
	if p.Name != nil {
		str = *p.Name + ":"
	}

	for _, item := range p.Items {
		if item.Negated {
			str += "!"
		}

		if item.Optional {
			str += "?"
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

func (p *PatternField) AddItem(item *PatternFieldItem) {
	p.Items = append(p.Items, item)
}

func (p Pattern) String() (str string) {
	for i, item := range p.Fields {
		str += item.String()

		if i != len(p.Fields)-1 {
			str += " "
		}
	}

	return
}

func (p *Pattern) AddID(id *string) {
	field := PatternField{}
	field.Items = append(field.Items, &PatternFieldItem{ID: id})

	p.Fields = append(p.Fields, &field)
}

func (p *Pattern) AddVar(id *string, negated, optional bool) {
	field := PatternField{}
	field.Items = append(field.Items, &PatternFieldItem{Var: id,
		Negated:  negated,
		Optional: optional,
	})

	p.Fields = append(p.Fields, &field)
}

func (p *Pattern) AddNum(num *string) {
	field := PatternField{}
	field.Items = append(field.Items, &PatternFieldItem{Num: num})

	p.Fields = append(p.Fields, &field)
}

func (p *Pattern) AddField(field *PatternField) {
	p.Fields = append(p.Fields, field)
}

func (p Pattern) LookupVariable(varName string) *PatternFieldItem {
	for _, field := range p.Fields {
		for _, item := range field.Items {
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

func (p Pattern) LookupFieldName(fieldName string) *PatternField {
	for _, field := range p.Fields {
		if field.Name == nil {
			continue
		}

		if *field.Name == fieldName {
			return field
		}
	}

	return nil
}
