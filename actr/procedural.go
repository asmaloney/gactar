package actr

type Procedural struct {
	BufferInterface // unused

	// "default_action_time": time that it takes to fire a production (seconds)
	// ccm: 0.05
	// pyactr: 0.05
	// vanilla: 0.05
	DefaultActionTime *float64
}

func NewProcedural() *Procedural {
	return &Procedural{BufferInterface: Buffer{}}
}

func (Procedural) GetModuleName() string {
	return "procedural"
}
