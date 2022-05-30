package actr

// Imaginal is a module which provides the ACT-R "imaginal" buffer.
type Imaginal struct {
	BufferInterface

	Delay float64 // non-negative time (in seconds) and defaults to .2
}

func NewImaginal() *Imaginal {
	// This uses the defaults as per ACT-R docs:
	// 	http://act-r.psy.cmu.edu/actr7.x/reference-manual.pdf page 276
	return &Imaginal{
		BufferInterface: Buffer{Name: "imaginal", MultipleInit: false},
		Delay:           0.2,
	}
}

func (i Imaginal) GetModuleName() string {
	return "imaginal"
}
