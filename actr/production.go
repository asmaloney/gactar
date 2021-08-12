package actr

// Production stores information on how to match buffers and perform some operations.
// It uses a small language to modify buffers and memories upon successful matches.
type Production struct {
	Name         string
	Matches      []*Match
	DoPython     []string
	DoStatements []*Statement
}

type Match struct {
	Buffer  *Buffer
	Memory  *Memory
	Text    *string
	Pattern *Pattern
}

type Statement struct {
	Clear  *ClearStatement
	Print  *PrintStatement
	Recall *RecallStatement
	Set    *SetStatement
	Write  *WriteStatement
}

// ClearStatement clears a list of buffers.
type ClearStatement struct {
	BufferNames []string
}

// PrintStatement outputs the string, id, or number to stdout.
type PrintStatement struct {
	Args []string // the strings, identifiers, or numbers to print
}

// RecallStatement is used to pull information from a memory.
type RecallStatement struct {
	Pattern    *Pattern
	Memory  *Memory
}

// WriteStatement will send the list of strings, ids, and numbers to the text output.
type WriteStatement struct {
	Args           []string // the strings, identifiers, or numbers to write
	TextOutputName string
}

// Slot represents either a slot position number or a name
type Slot struct {
	ArgNum *int
	Name   *string
}

// SetStatement will set a slot or the entire contents of the named buffer to a string or a pattern.
type SetStatement struct {
	Slot       *Slot  // set this slot
	Buffer *Buffer // of this buffer

	ID      *string  // to this ID OR
	Number  *string  // OR this number (no need to store as actual number at the moment)
	String  *string  // OR this string
	Pattern *Pattern // OR this pattern
}

func (p Production) LookupMatchByBuffer(bufferName string) *Match {
	for _, m := range p.Matches {
		if m.Buffer.Name == bufferName {
			return m
		}
	}

	return nil
}

// LookupMatchByVariable checks all matches for a variable by name.
// This is pretty inefficient, but given the small number of matches
// in a production, it's probably not worth doing anything more complicated.
// We could store all the vars used in all the matches on the Match itself
// and look it up there.
func (p Production) LookupMatchByVariable(varName string) *Match {
	for _, m := range p.Matches {
		if m.Pattern == nil {
			return nil
		}

		patternItem := m.Pattern.LookupVariable(varName)
		if patternItem != nil {
			return m
		}
	}

	return nil
}
