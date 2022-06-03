package modules

import "github.com/asmaloney/gactar/actr/buffer"

// Goal is a module which provides the ACT-R "goal" buffer.
type Goal struct {
	buffer.BufferInterface
}

func NewGoal() *Goal {
	return &Goal{
		BufferInterface: &buffer.Buffer{Name: "goal", MultipleInit: false},
	}
}

func (g Goal) ModuleName() string {
	return "goal"
}

func (g *Goal) SetParam(param *Param) (err ParamError) {
	return
}
