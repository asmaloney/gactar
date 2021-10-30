package actr

// Goal is a module which provides the ACT-R "goal" buffer.
type Goal struct {
	BufferInterface
}

func NewGoal() *Goal {
	return &Goal{
		BufferInterface: Buffer{Name: "goal"},
	}
}

func (g Goal) GetModuleName() string {
	return "goal"
}

func (g *Goal) AllowsMultipleInit() bool {
	return false
}
