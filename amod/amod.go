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
	fmt.Println(amodParser.String())
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

// ParseChunk is used to parse goals when given as input from a user.
func ParseChunk(model *actr.Model, chunk string) (*actr.Pattern, error) {
	if chunk == "" {
		return nil, nil
	}

	// "chunk" needs backticks in order to parse properly.
	if !strings.HasPrefix(chunk, "`") {
		chunk = "`" + chunk + "`"
	}

	var p pattern

	r := strings.NewReader(chunk)

	err := patternParser.Parse("", r, &p)
	if err != nil {
		return nil, err
	}

	err = validatePattern(model, &p)
	if err != nil {
		return nil, err
	}

	return createChunkPattern(model, &p)
}

// generateModel runs through the parsed structures and creates an actr.Model from them
func generateModel(amod *amodFile) (model *actr.Model, err error) {
	model = &actr.Model{
		Name:        amod.Model.Name,
		Description: amod.Model.Description,
	}

	model.Initialize()

	err = addConfig(model, amod.Config)
	if err != nil {
		return
	}

	err = addExamples(model, amod.Model.Examples)
	if err != nil {
		return
	}

	err = addInit(model, amod.Init)
	if err != nil {
		return
	}

	err = addProductions(model, amod.Productions)
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
	addChunks(model, config.ChunkDecls, &errs)
	addMemory(model, config.MemoryDecl, &errs)

	return errs.ErrorOrNil()
}

func addExamples(model *actr.Model, examples []*pattern) (err error) {
	if len(examples) == 0 {
		return
	}

	errs := errorListWithContext{}

	for _, example := range examples {
		err = validatePattern(model, example)
		if err != nil {
			errs.AddErrorIfNotNil(err)
			continue
		}

		pattern, err := createChunkPattern(model, example)
		if err != nil {
			errs.AddErrorIfNotNil(err)
			continue
		}

		model.Examples = append(model.Examples, pattern)
	}

	return errs.ErrorOrNil()
}

func addACTR(model *actr.Model, list []*field, errs *errorListWithContext) {
	if list == nil {
		return
	}

	for _, field := range list {
		switch field.Key {
		case "log":
			if field.Value.Number != nil {
				model.Logging = (*field.Value.Number != 0)
			} else if field.Value.ID != nil {
				model.Logging = (strings.ToLower(*field.Value.ID) == "true")
			} else if field.Value.Str != nil {
				model.Logging = (strings.ToLower(*field.Value.Str) == "true")
			}
		default:
			errs.Addc(&field.Pos, "unrecognized field in actr section: '%s'", field.Key)
		}
	}
}

func addChunks(model *actr.Model, chunks []*chunkDecl, errs *errorListWithContext) {
	if chunks == nil {
		return
	}

	for _, chunk := range chunks {
		err := validateChunk(model, chunk)
		if err != nil {
			errs.AddErrorIfNotNil(err)
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

func addMemory(model *actr.Model, mem []*field, errs *errorListWithContext) {
	if mem == nil {
		return
	}

	memory := model.LookupMemory("memory")
	if memory == nil {
		errs.Add("could not find memory on model")
		return
	}

	for _, field := range mem {
		switch field.Key {
		case "latency":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory latency '%s' must be a number", field.Value.String())
				continue
			}

			memory.Latency = field.Value.Number

		case "threshold":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory threshold '%s' must be a number", field.Value.String())
				continue
			}

			memory.Threshold = field.Value.Number

		case "max_time":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory max_time '%s' must be a number", field.Value.String())
				continue
			}

			memory.MaxTime = field.Value.Number

		case "finst_size":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory finst_size '%s' must be a number", field.Value.String())
				continue
			}

			size := int(*field.Value.Number)
			memory.FinstSize = &size

		case "finst_time":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory finst_time '%s' must be a number", field.Value.String())
				continue
			}

			memory.FinstTime = field.Value.Number

		default:
			errs.Addc(&field.Pos, "unrecognized field '%s' in memory", field.Key)
		}
	}
}

func addInit(model *actr.Model, init *initSection) (err error) {
	if init == nil {
		return
	}

	errs := errorListWithContext{}

	for _, initialization := range init.Initializations {
		name := initialization.Name
		buffer := model.LookupBuffer(name)
		memory := model.LookupMemory(name)

		err := validateInitialization(model, initialization)
		if err != nil {
			errs.AddErrorIfNotNil(err)
			continue
		}

		if buffer != nil {
			pattern, err := createChunkPattern(model, initialization.InitPattern)
			if err != nil {
				errs.AddErrorIfNotNil(err)
				continue
			}

			init := actr.Initializer{
				Buffer:  buffer,
				Pattern: pattern,
			}

			model.Initializers = append(model.Initializers, &init)
		} else { // memory
			for _, init := range initialization.InitPatterns {
				pattern, err := createChunkPattern(model, init)
				if err != nil {
					errs.AddErrorIfNotNil(err)
					continue
				}

				init := actr.Initializer{
					Memory:  memory,
					Pattern: pattern,
				}

				model.Initializers = append(model.Initializers, &init)
			}
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
			Name:        production.Name,
			Description: production.Description,
		}

		err := validateMatch(production.Match, model, &prod)
		if err != nil {
			errs.AddErrorIfNotNil(err)
			continue
		}

		for _, item := range production.Match.Items {
			pattern, err := createChunkPattern(model, item.Pattern)
			if err != nil {
				errs.AddErrorIfNotNil(err)
				continue
			}

			name := item.Name

			prod.Matches = append(prod.Matches, &actr.Match{
				Buffer:  model.LookupBuffer(name),
				Memory:  model.LookupMemory(name),
				Pattern: pattern,
			})

		}

		if production.Do.Statements != nil {
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
			newItem := &actr.PatternSlotItem{
				Negated: item.Not,
			}

			if item.Nil != nil {
				newItem.Nil = true
			} else if item.ID != nil {
				newItem.ID = item.ID
			} else if item.Num != nil {
				newItem.Num = item.Num
			} else if item.Var != nil {
				newItem.Var = item.Var
			}

			slot.AddItem(newItem)
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
	} else {
		err = fmt.Errorf("statement type not handled: %T", statement)
		return err
	}

	if err != nil {
		return err
	}

	if s != nil {
		production.DoStatements = append(production.DoStatements, s)
	}

	return nil
}

func addSetStatement(model *actr.Model, set *setStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateSetStatement(set, model, production)
	if err != nil {
		return nil, err
	}

	buffer := model.LookupBuffer(set.BufferName)

	s := actr.Statement{Set: &actr.SetStatement{
		Buffer: buffer,
	}}
	createNewStatement := true
	setStatement := production.LookupSetStatementByBuffer(set.BufferName)
	if setStatement != nil {
		// If we found one, use its set statement
		createNewStatement = false
		s.Set = setStatement
	}

	if set.Slot != nil {
		bufferName := buffer.GetName()

		// find slot index in chunk
		match := production.LookupMatchByBuffer(bufferName)
		if match == nil {
			err = fmt.Errorf("could not find buffer match '%s' in production '%s'", bufferName, production.Name)
			return nil, err
		}

		s.Set.Chunk = match.Pattern.Chunk

		slotName := *set.Slot
		index := match.Pattern.Chunk.GetSlotIndex(slotName)
		if index == -1 {
			err = fmt.Errorf("could not find slot named '%s' in buffer match '%s' in production '%s'", slotName, bufferName, production.Name)
			return nil, err
		}

		value := &actr.SetValue{}

		if set.Value.Var != nil {
			varName := strings.TrimPrefix(*set.Value.Var, "?")
			value.Var = &varName
		} else if set.Value.Nil != nil {
			value.Nil = *set.Value.Nil
		} else if set.Value.Number != nil {
			value.Number = set.Value.Number
		} else if set.Value.Str != nil {
			value.Str = set.Value.Str
		}

		newSlot := &actr.SetSlot{
			Name:      *set.Slot,
			SlotIndex: index,
			Value:     value,
		}

		s.Set.AddSlot(newSlot)

	} else if set.Pattern != nil {
		pattern, err := createChunkPattern(model, set.Pattern)
		if err != nil {
			return nil, err
		}

		s.Set.Pattern = pattern
	}

	if !createNewStatement {
		return nil, nil
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
		p.Values = convertArgs(print.Args)
	}

	s := actr.Statement{Print: &p}

	return &s, nil
}

func convertArgs(args []*arg) *[]*actr.Value {
	actrValues := []*actr.Value{}

	for _, v := range args {
		newValue := actr.Value{}

		if v.Var != nil {
			newValue.Var = v.Var
		} else if v.ID != nil {
			newValue.ID = v.ID
		} else if v.Str != nil {
			newValue.Str = v.Str
		} else if v.Number != nil {
			newValue.Number = v.Number
		}

		actrValues = append(actrValues, &newValue)
	}

	return &actrValues
}
