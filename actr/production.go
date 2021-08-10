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
	Name    string
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
	MemoryName string
}

// WriteStatement will send the list of strings, ids, and numbers to the text output.
type WriteStatement struct {
	Args           []string // the strings, identifiers, or numbers to write
	TextOutputName string
}

type SetField struct {
	ArgNum    *int
	FieldName *string
}

// SetStatement will set a field or the entire contents of the named buffer to a string or a pattern.
type SetStatement struct {
	Field      *SetField // set this field
	BufferName string    // of this buffer

	Text    *string  // to this string OR
	Pattern *Pattern // this pattern
}

func (p Production) LookupMatchByBuffer(bufferName string) *Match {
	for _, m := range p.Matches {
		if m.Name == bufferName {
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
