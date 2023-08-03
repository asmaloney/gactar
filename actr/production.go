package actr

import (
	"fmt"

	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/modules"
)

// Production stores information on how to match buffers and perform some operations.
// It uses a small language to modify states upon successful matches.
type Production struct {
	Model *Model // link back to its model so we can add implicit chunks

	Name        string
	Description *string // optional description to output as a comment in the generated code

	VarIndexMap map[string]VarIndex // track the buffer and slot name each variable refers to

	Matches      []*Match
	DoStatements []*Statement

	AMODLineNumber int // line number in the amod file of the this production
}

// VarIndex is used to track which buffer slot a variable refers to
type VarIndex struct {
	Var      *PatternVar
	Buffer   buffer.Interface
	SlotName string
}

type Comparison int

const (
	Equal Comparison = iota
	NotEqual
)

func (c Comparison) String() string {
	switch c {
	case Equal:
		return "=="
	case NotEqual:
		return "!="
	}

	return "unknown"
}

type Constraint struct {
	LHS        *string
	Comparison Comparison
	RHS        *Value
}

func (c Constraint) String() string {
	return fmt.Sprintf("%s %s %s", *c.LHS, c.Comparison, c.RHS)
}

type Match struct {
	BufferPattern *BufferPatternMatch
	BufferState   *BufferStateMatch
	ModuleState   *ModuleStateMatch
}

type BufferPatternMatch struct {
	Buffer buffer.Interface

	Pattern *Pattern
}

type BufferStateMatch struct {
	Buffer buffer.Interface

	State string
}

type ModuleStateMatch struct {
	Module modules.Interface

	// The generated code for the frameworks actually uses a buffer name, not the module name.
	// So store (one) here for convenience. If the module has multiple buffers it should not
	// matter which one we pick as the requests should be on its module.
	Buffer buffer.Interface

	State string
}

type Statement struct {
	Clear  *ClearStatement
	Print  *PrintStatement
	Recall *RecallStatement
	Set    *SetStatement
	Stop   *StopStatement
}

// ClearStatement clears a list of buffers.
type ClearStatement struct {
	BufferNames []string
}

// Value holds something that may be printed.
type Value struct {
	Nil    *bool // set this to nil
	Var    *string
	ID     *string
	Str    *string
	Number *string
}

func (v Value) String() string {
	switch {
	case v.Nil != nil:
		return "nil"
	case v.Var != nil:
		return *v.Var
	case v.ID != nil:
		return *v.ID
	case v.Str != nil:
		return *v.Str
	case v.Number != nil:
		return *v.Number
	}

	return "unknown"
}

// PrintStatement outputs the string, id, or number to stdout.
type PrintStatement struct {
	Values *[]*Value
}

// RecallStatement is used to pull information from memory.
type RecallStatement struct {
	Pattern           *Pattern
	MemoryModuleName  string
	RequestParameters map[string]string
}

type SetSlot struct {
	Name      string
	SlotIndex int // (this slot index in the chunk)
	Value     *Value
}

// SetStatement will set a slot or the entire contents of the named buffer to a string or a pattern.
// There are two forms:
//
//	(1) set (Buffer).(SetSlot) to (SetSlot.Value)
//	(2) set (Buffer) to (Pattern)
type SetStatement struct {
	Buffer buffer.Interface // buffer we are manipulating

	Slots *[]SetSlot // (1) set this slot
	Chunk *Chunk     // (1) if we are setting slots, point at the chunk they reference for easy lookup

	Pattern *Pattern // (2) pattern if we are setting the whole buffer
}

// StopStatement outputs a stop command. There are no parameters.
type StopStatement struct {
}

// AddDoStatement adds the statement to our list and adds any IDs to ImplicitChunks
// so we can (possibly) create them in the framework output.
func (p *Production) AddDoStatement(statement *Statement) {
	if statement == nil {
		return
	}

	p.DoStatements = append(p.DoStatements, statement)

	switch {
	case statement.Set != nil:
		if statement.Set.Slots != nil {
			for _, slot := range *statement.Set.Slots {
				if slot.Value.ID != nil {
					p.Model.AddImplicitChunk(*slot.Value.ID)
				}
			}
		} else if statement.Set.Pattern != nil {
			p.Model.AddImplicitChunksFromPattern(statement.Set.Pattern)
		}

	case statement.Recall != nil:
		p.Model.AddImplicitChunksFromPattern(statement.Recall.Pattern)
	}

}

func (p Production) LookupMatchByBuffer(bufferName string) *BufferPatternMatch {
	for _, m := range p.Matches {
		if m.BufferPattern == nil {
			continue
		}

		if m.BufferPattern.Buffer.BufferName() == bufferName {
			return m.BufferPattern
		}
	}

	return nil
}

// LookupSetStatementByBuffer is used when combining several set consecutive statements on one buffer.
// So this:
//
//	set goal.foo to 1
//	set goal.bar to 10
//
// is treated like this:
//
//	set foo, bar on goal to 1, 10
func (p Production) LookupSetStatementByBuffer(bufferName string) *SetStatement {
	if len(p.DoStatements) == 0 {
		return nil
	}

	last := p.DoStatements[len(p.DoStatements)-1]

	if (last.Set == nil) || (last.Set.Slots == nil) {
		return nil
	}

	if last.Set.Buffer.BufferName() == bufferName {
		return last.Set
	}

	return nil
}

// AddSlotToSetStatement is used to add slot info to an already-existing set statement.
// See LookupSetStatementByBuffer() above.
func (p *Production) AddSlotToSetStatement(statement *SetStatement, slot *SetSlot) {
	if statement.Slots == nil {
		statement.Slots = &[]SetSlot{}
	}

	*statement.Slots = append(*statement.Slots, *slot)

	if slot.Value.ID != nil {
		p.Model.AddImplicitChunk(*slot.Value.ID)
	}
}

// LookupMatchByVariable checks all matches for a variable by name.
// This is pretty inefficient, but given the small number of matches
// in a production, it's probably not worth doing anything more complicated.
// We could store all the vars used in all the matches on the Match itself
// and look it up there.
func (p Production) LookupMatchByVariable(varName string) *BufferPatternMatch {
	for _, m := range p.Matches {
		if m.BufferPattern == nil {
			continue
		}

		if m.BufferPattern.Pattern == nil {
			return nil
		}

		patternItem := m.BufferPattern.Pattern.LookupVariable(varName)
		if patternItem != nil {
			return m.BufferPattern
		}
	}

	return nil
}
