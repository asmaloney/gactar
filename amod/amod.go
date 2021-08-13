package amod

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"gitlab.com/asmaloney/gactar/actr"
	"gitlab.com/asmaloney/gactar/errorlist"
)

var debugging bool = false

type errorListWithContext struct {
	errorlist.Errors
}

func (err *errorListWithContext) Addc(pos *lexer.Position, e string, a ...interface{}) {
	str := fmt.Sprintf(e, a...)
	err.Addf("%s (line %d)", str, pos.Line)
}

// SetDebug turns debugging on and off. This will output the tokens as they are generated.
func SetDebug(debug bool) {
	debugging = debug
}

// OutputEBNF outputs the extended Backus–Naur form (EBNF) of the amod grammar to stdout.
// See: https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form
func OutputEBNF() {
	fmt.Println(parser.String())
}

// GenerateModel generates a model from the text in the buffer.
func GenerateModel(buffer string) (model *actr.Model, err error) {
	r := strings.NewReader(buffer)

	amod, err := parse(r)
	if err != nil {
		return
	}

	return generateModel(amod)
}

// GenerateModelFromFile generates a model from the file 'fileName'.
func GenerateModelFromFile(fileName string) (model *actr.Model, err error) {
	amod, err := parseFile(fileName)
	if err != nil {
		return
	}

	return generateModel(amod)
}

// generateModel runs through the parsed structures and creates an actr.Model from them
func generateModel(amod *amodFile) (model *actr.Model, err error) {
	model = &actr.Model{
		Name:        amod.Model.Name,
		Description: amod.Model.Description,
		Examples:    amod.Model.Examples,
	}

	model.CreateBuffersAndMemory()

	err = addConfig(model, amod.Config)
	if err != nil {
		return
	}

	err = initialize(model, amod.Init)
	if err != nil {
		return
	}

	err = addProductions(model, amod.Productions)
	if err != nil {
		return
	}

	err = validateInitializers(model, amod.Init)
	if err != nil {
		return
	}

	return
}

func addConfig(model *actr.Model, config *configSection) (err error) {
	if config == nil {
		return
	}

	errs := errorListWithContext{}

	addACTR(model, config.ACTR, &errs)
	addChunks(model, config.Chunks, &errs)
	addMemory(model, config.Memory, &errs)
	addTextOutputs(model, config.TextOutputs, &errs)

	return errs.ErrorOrNil()
}

func addACTR(model *actr.Model, list *fieldList, errs *errorListWithContext) {
	if list == nil {
		return
	}

	for _, field := range list.Fields {
		switch field.Key {
		case "log":
			if field.Value.Number != nil {
				model.Logging = (*field.Value.Number != 0)
			} else {
				model.Logging = (strings.ToLower(*field.Value.String) == "true")
			}
		default:
			errs.Addc(&field.Pos, "unrecognized field in actr section: '%s'", field.Key)
		}
	}
}

func addChunks(model *actr.Model, chunks []*chunk, errs *errorListWithContext) {
	if chunks == nil {
		return
	}

	for _, chunk := range chunks {
		c := model.LookupChunk(chunk.Name)
		if c != nil {
			errs.Addc(&chunk.Pos, "duplicate chunk name: '%s'", chunk.Name)
			continue
		}

		aChunk := actr.Chunk{
			Name:      chunk.Name,
			SlotNames: chunk.SlotNames,
			NumSlots:  len(chunk.SlotNames),
		}

		model.Chunks = append(model.Chunks, &aChunk)
	}
}

func addMemory(model *actr.Model, mem *memory, errs *errorListWithContext) {
	if mem == nil {
		return
	}

	memory := model.LookupMemory("memory")
	if memory == nil {
		errs.Add("could not find memory on model")
		return
	}

	for _, field := range mem.Fields.Fields {
		switch field.Key {
		case "latency":
			if field.Value.String != nil {
				errs.Addc(&field.Pos, "memory latency '%s' should not be a string", *field.Value.String)
				continue
			}

			memory.Latency = field.Value.Number

		case "threshold":
			if field.Value.String != nil {
				errs.Addc(&field.Pos, "memory threshold '%s' should not be a string", *field.Value.String)
				continue
			}

			memory.Threshold = field.Value.Number

		case "max_time":
			if field.Value.String != nil {
				errs.Addc(&field.Pos, "memory max_time '%s' should not be a string", *field.Value.String)
				continue
			}

			memory.MaxTime = field.Value.Number

		case "finst_size":
			if field.Value.String != nil {
				errs.Addc(&field.Pos, "memory finst_size '%s' should not be a string", *field.Value.String)
				continue
			}

			size := int(*field.Value.Number)
			memory.FinstSize = &size

		case "finst_time":
			if field.Value.String != nil {
				errs.Addc(&field.Pos, "memory finst_time '%s' should not be a string", *field.Value.String)
				continue
			}

			memory.FinstTime = field.Value.Number

		default:
			errs.Addc(&field.Pos, "unrecognized field '%s' in memory", field.Key)
		}
	}
}

func addTextOutputs(model *actr.Model, list *identList, errs *errorListWithContext) {
	if list == nil {
		return
	}

	for _, name := range list.Identifiers {
		textOutput := actr.TextOutput{
			Name: name,
		}

		model.TextOutputs = append(model.TextOutputs, &textOutput)
	}
}

func initialize(model *actr.Model, init *initSection) (err error) {
	if init == nil {
		return
	}

	errs := errorListWithContext{}

	for _, initializer := range init.Initializers {
		memory := model.LookupMemory("memory")
		if memory == nil {
			errs.Addc(&initializer.Pos, "memory not found")
			continue
		}

		if initializer.Items == nil {
			errs.Addc(&initializer.Pos, "no memory initializers found")
			continue
		}

		for _, item := range initializer.Items.Strings {
			init := actr.Initializer{
				Memory: memory,
				Text:   item,
			}

			model.Initializers = append(model.Initializers, &init)
		}
	}

	return errs.ErrorOrNil()
}

func addProductions(model *actr.Model, productions *productionSection) (err error) {
	if productions == nil {
		return
	}

	errs := errorListWithContext{}

	for _, production := range productions.Productions {
		prod := actr.Production{
			Name: production.Name,
		}

		for _, item := range production.Match.Items {
			name := item.Name

			buffer := model.LookupBuffer(name)
			memory := model.LookupMemory(name)

			if (buffer == nil) && (memory == nil) {
				errs.Addc(&item.Pos, "buffer or memory '%s' not found in production '%s'", name, prod.Name)
				continue
			}

			if item.Pattern == nil {
				errs.Addc(&item.Pos, "invalid pattern for '%s' in production '%s'", name, prod.Name)
				continue
			}

			pattern, err := createChunkPattern(model, item.Pattern)
			if err != nil {
				errs.AddErrorIfNotNil(err)

			} else {
				prod.Matches = append(prod.Matches, &actr.Match{
					Buffer:  buffer,
					Memory:  memory,
					Pattern: pattern,
				})
			}
		}

		if production.Do.PyCode != nil {
			prod.DoPython = append(prod.DoPython, *production.Do.PyCode...)
		} else if production.Do.Statements != nil {
			for _, statement := range *production.Do.Statements {
				err := addStatement(model, statement, &prod)
				errs.AddErrorIfNotNil(err)
			}
		}

		model.Productions = append(model.Productions, &prod)
	}

	return errs.ErrorOrNil()
}

func createChunkPattern(model *actr.Model, cp *pattern) (*actr.Pattern, error) {
	errs := errorListWithContext{}

	chunk := model.LookupChunk(cp.ChunkName)
	if chunk == nil {
		errs.Addc(&cp.Pos, "could not find chunk named '%s'", cp.ChunkName)
		return nil, errs
	}

	pattern := actr.Pattern{
		Chunk: chunk,
	}
	for _, s := range cp.Slots {

		slot := actr.PatternSlot{}

		for _, item := range s.Items {
			if item.ID != nil {
				slot.AddItem(&actr.PatternSlotItem{ID: item.ID})
			} else if item.Num != nil {
				slot.AddItem(&actr.PatternSlotItem{Num: item.Num})
			} else if item.Var != nil {
				slot.AddItem(&actr.PatternSlotItem{Var: item.Var})
			} else if item.NotVar != nil {
				slot.AddItem(&actr.PatternSlotItem{Var: item.NotVar, Negated: true})
			}
		}
		pattern.AddSlot(&slot)
	}

	return &pattern, nil
}

func addStatement(model *actr.Model, statement *statement, production *actr.Production) (err error) {
	var s *actr.Statement

	if statement.Set != nil {
		s, err = addSetStatement(model, statement.Set, production)
	} else if statement.Recall != nil {
		s, err = addRecallStatement(model, statement.Recall, production)
	} else if statement.Clear != nil {
		s, err = addClearStatement(model, statement.Clear, production)
	} else if statement.Print != nil {
		s, err = addPrintStatement(model, statement.Print, production)
	} else if statement.Write != nil {
		s, err = addWriteStatement(model, statement.Write, production)
	} else {
		err = fmt.Errorf("statement type not handled: %T", statement)
		return err
	}

	if err != nil {
		return err
	}

	production.DoStatements = append(production.DoStatements, s)

	return nil
}

func addSetStatement(model *actr.Model, set *setStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateSetStatement(set, model, production)
	if err != nil {
		return nil, err
	}

	buffer := model.LookupBuffer(set.BufferName)

	s := actr.Statement{
		Set: &actr.SetStatement{
			Buffer: buffer,
		},
	}

	if set.Slot != nil {
		s.Set.Slot = set.Slot

		// find slot index in chunk
		match := production.LookupMatchByBuffer(buffer.Name)
		if match == nil {
			err = fmt.Errorf("could not find buffer match '%s' in production '%s'", buffer.Name, production.Name)
			return nil, err
		}

		s.Set.Chunk = match.Pattern.Chunk

		slotName := *set.Slot
		index := match.Pattern.Chunk.GetSlotIndex(slotName)
		if index == -1 {
			err = fmt.Errorf("could not find slot named '%s' in buffer match '%s' in production '%s'", slotName, buffer.Name, production.Name)
			return nil, err
		}

		s.Set.SlotIndex = index
	}

	if set.Pattern != nil {
		pattern, err := createChunkPattern(model, set.Pattern)
		if err != nil {
			return nil, err
		}

		s.Set.Pattern = pattern
	} else if set.ID != nil {
		s.Set.ID = set.ID
	} else if set.Number != nil {
		s.Set.Number = set.Number
	} else if set.String != nil {
		s.Set.String = set.String
	}

	return &s, nil
}

func addRecallStatement(model *actr.Model, recall *recallStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateRecallStatement(recall, model, production)
	if err != nil {
		return nil, err
	}

	pattern, err := createChunkPattern(model, recall.Pattern)
	if err != nil {
		return nil, err
	}

	s := actr.Statement{
		Recall: &actr.RecallStatement{
			Pattern: pattern,
			Memory:  model.LookupMemory("memory"),
		},
	}

	return &s, nil
}

func addClearStatement(model *actr.Model, clear *clearStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateClearStatement(clear, model, production)
	if err != nil {
		return nil, err
	}

	s := actr.Statement{
		Clear: &actr.ClearStatement{
			BufferNames: clear.BufferNames,
		},
	}

	return &s, nil
}

func addPrintStatement(model *actr.Model, print *printStatement, production *actr.Production) (*actr.Statement, error) {
	err := validatePrintStatement(print, model, production)
	if err != nil {
		return nil, err
	}

	p := actr.PrintStatement{}
	if print.Args != nil {
		p.Args = print.Args.Strings()
	}

	s := actr.Statement{Print: &p}

	return &s, nil
}

func addWriteStatement(model *actr.Model, write *writeStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateWriteStatement(write, model, production)
	if err != nil {
		return nil, err
	}

	s := actr.Statement{
		Write: &actr.WriteStatement{
			Args:           write.Args.Strings(),
			TextOutputName: write.TextOutputName},
	}

	return &s, nil
}
