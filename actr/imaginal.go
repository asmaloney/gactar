package actr

type Imaginal struct {
	Buffer

	Delay float64 // non-negative time (in seconds) and defaults to .2
}

func (model *Model) CreateImaginal() *Imaginal {
	// This uses the defaults as per ACT-R docs:
	// 	http://act-r.psy.cmu.edu/actr7.x/reference-manual.pdf page 276

	imaginal := &Imaginal{
		Buffer: Buffer{Name: "imaginal"},

		Delay: 0.2,
	}

	model.Buffers = append(model.Buffers, imaginal)
	model.HasImaginal = true

	return imaginal
}

// GetImaginal gets the imaginal buffer (or returns nil if it does not exist).
func (model Model) GetImaginal() *Imaginal {
	buffer := model.LookupBuffer("imaginal")
	if buffer == nil {
		return nil
	}

	imaginal, ok := buffer.(*Imaginal)
	if !ok {
		return nil
	}

	return imaginal
}
