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

	Negated  bool
	Optional bool
}

func (c PatternField) String() (str string) {
	if c.Name != nil {
		str = *c.Name + ":"
	}

	for _, item := range c.Items {
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

func (c Pattern) String() (str string) {
	for i, item := range c.Items {
		if item.ID != nil {
			str += *item.ID
		} else if item.Var != nil {
			str += *item.Var
		} else if item.Num != nil {
			str += *item.Num
		} else if item.Field != nil {
			str += item.Field.String()
		}

		if i != len(c.Items)-1 {
			str += " "
		}
	}

	return
}

func (c *Pattern) AddID(id *string) {
	c.Items = append(c.Items, &PatternItem{ID: id})
}

func (c *Pattern) AddVar(id *string) {
	c.Items = append(c.Items, &PatternItem{Var: id})
}

func (c *Pattern) AddNum(num *string) {
	c.Items = append(c.Items, &PatternItem{Num: num})
}

func (c *Pattern) AddField(field *PatternField) {
	c.Items = append(c.Items, &PatternItem{Field: field})
}

func (c Pattern) HasVar(id string) bool {
	for _, item := range c.Items {
		if item.Var == nil {
			continue
		}

		if *item.Var == id {
			return true
		}
	}

	return false
}
