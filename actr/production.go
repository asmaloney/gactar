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
	Buffer *Buffer // buffer
	Memory *Memory // OR memory

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

// Value holds something that may be printed.
type Value struct {
	Var    *string
	ID     *string
	Str    *string
	Number *string
}

// PrintStatement outputs the string, id, or number to stdout.
type PrintStatement struct {
	Values *[]*Value
}

// RecallStatement is used to pull information from a memory.
type RecallStatement struct {
	Pattern *Pattern
	Memory  *Memory
}

// WriteStatement will send the list of strings, ids, and numbers to the text output.
type WriteStatement struct {
	Values         *[]*Value
	TextOutputName string
}

type SetValue struct {
	Var    *string // set to this Var OR
	Number *string // OR this number (no need to store as actual number at the moment)
	Str    *string // OR this string
}

func (s SetValue) String() string {
	if s.Var != nil {
		return *s.Var
	} else if s.Number != nil {
		return *s.Number
	} else if s.Str != nil {
		return "'" + *s.Str + "'"
	}

	return ""
}

type SetSlot struct {
	Name      string
	SlotIndex int // (this slot index in the chunk)
	Value     *SetValue
}

// SetStatement will set a slot or the entire contents of the named buffer to a string or a pattern.
// There are two forms:
//	(1) set (SetSlot) of (Buffer) to (SetValue)
//	(2) set (Buffer) to (Pattern)
type SetStatement struct {
	Slots *[]SetSlot // (1) set this slot (optional)
	Chunk *Chunk     // (1) if we are setting slots, point at the chunk they reference for easy lookup

	Buffer *Buffer // (1 & 2) buffer we are manipulating

	Pattern *Pattern // (2) pattern if we are setting the whole buffer
}

func (p Production) LookupMatchByBuffer(bufferName string) *Match {
	for _, m := range p.Matches {
		if m.Buffer.Name == bufferName {
			return m
		}
	}

	return nil
}

// LookupSetStatementByBuffer is used when combining several set consecutive statements on one buffer.
// So this:
//		set foo on goal to 1
//		set bar on goal to 10
// is treated like this:
//		set foo, bar on goal to 1, 10
func (p Production) LookupSetStatementByBuffer(bufferName string) *SetStatement {
	if len(p.DoStatements) == 0 {
		return nil
	}

	last := p.DoStatements[len(p.DoStatements)-1]

	if (last.Set == nil) || (last.Set.Slots == nil) {
		return nil
	}

	if last.Set.Buffer.Name == bufferName {
		return last.Set
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

func (s *SetStatement) AddSlot(slot *SetSlot) {
	if s.Slots == nil {
		newSlots := []SetSlot{}
		s.Slots = &newSlots
	}

	*s.Slots = append(*s.Slots, *slot)
}
