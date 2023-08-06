package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
)

// Goal is a module which provides the ACT-R "goal" buffer.
type Goal struct {
	Module
}

func (g Goal) Buffer() buffer.Interface {
	return g.BufferList.At(0)
}

// NewGoal creates and returns a new Goal module
func NewGoal() *Goal {
	goalBuff := buffer.NewBuffer("goal", 0.0, nil)

	return &Goal{
		Module: Module{
			Name:         "goal",
			Version:      BuiltIn,
			Description:  "provides a goal buffer for the model",
			BufferList:   buffer.List{goalBuff},
			MultipleInit: false,
		},
	}
}
