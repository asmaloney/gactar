package actr

type BufferInterface interface {
	GetName() string
}

type Buffer struct {
	Name string
}

func (b Buffer) GetName() string {
	return b.Name
}

func (b Buffer) String() string {
	return b.Name
}

// LookupBuffer looks up the named buffer in the model and returns it (or nil if it does not exist).
func (model Model) LookupBuffer(bufferName string) BufferInterface {
	for _, buf := range model.Buffers {
		if buf.GetName() == bufferName {
			return buf
		}
	}

	return nil
}
