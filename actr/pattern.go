package actr

// We need to take apart patterns such as 'add ?num1 ?num2 count:?count!?num2 sum:?sum'
// so we can verify variable use.

type Pattern struct {
	Items []*PatternItem
}

type PatternItem struct {
	ID    *string
	Var   *string
	Num   *string // we don't need to treat this as a number anywhere, so keep as a string
	Field *PatternField
}

type PatternField struct {
	Name  *string
	Items []PatternFieldItem
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

func (p Pattern) String() (str string) {
	for i, item := range p.Items {
		if item.ID != nil {
			str += *item.ID
		} else if item.Var != nil {
			str += *item.Var
		} else if item.Num != nil {
			str += *item.Num
		} else if item.Field != nil {
			str += item.Field.String()
		}

		if i != len(p.Items)-1 {
			str += " "
		}
	}

	return
}

func (p *Pattern) AddID(id *string) {
	p.Items = append(p.Items, &PatternItem{ID: id})
}

func (p *Pattern) AddVar(id *string) {
	p.Items = append(p.Items, &PatternItem{Var: id})
}

func (p *Pattern) AddNum(num *string) {
	p.Items = append(p.Items, &PatternItem{Num: num})
}

func (p *Pattern) AddField(field *PatternField) {
	p.Items = append(p.Items, &PatternItem{Field: field})
}

func (p Pattern) LookupFieldName(fieldName string) *PatternField {
	for _, item := range p.Items {
		if (item.Field == nil) || (item.Field.Name == nil) {
			continue
		}

		if *item.Field.Name == fieldName {
			return item.Field
		}
	}

	return nil
}
