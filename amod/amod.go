package amod

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/errorlist"
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

	// "chunk" needs square brackets in order to parse properly.
	if !strings.HasPrefix(chunk, "[") {
		chunk = "[" + chunk
	}
	if !strings.HasSuffix(chunk, "]") {
		chunk += "]"
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

	addGACTAR(model, config.GACTAR, &errs)
	addModules(model, config.Modules, &errs)
	addChunks(model, config.ChunkDecls, &errs)

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

func addGACTAR(model *actr.Model, list []*field, errs *errorListWithContext) {
	if list == nil {
		return
	}

	for _, field := range list {
		switch field.Key {
		case "log_level":
			if (field.Value.Str == nil) || !actr.ValidLogLevel(*field.Value.Str) {
				errs.Addc(&field.Pos, "log_level '%s' must be 'min', 'info', 'or 'detail'", field.Value.String())
				continue
			}

			model.LogLevel = actr.ACTRLogLevel(*field.Value.Str)

		default:
			errs.Addc(&field.Pos, "unrecognized field in gactar section: '%s'", field.Key)
		}
	}
}

func addModules(model *actr.Model, modules []*module, errs *errorListWithContext) {
	if modules == nil {
		return
	}

	for _, module := range modules {
		switch module.Name {
		case "imaginal":
			addImaginal(model, module.InitFields, errs)
		case "memory":
			addMemory(model, module.InitFields, errs)
		default:
			errs.Addc(&module.Pos, "unrecognized module in config: '%s'", module.Name)
		}
	}
}

func addImaginal(model *actr.Model, fields []*field, errs *errorListWithContext) {
	imaginal := model.CreateImaginal()
	if imaginal == nil {
		errs.Add("could not create imaginal buffer on model")
		return
	}

	for _, field := range fields {
		switch field.Key {
		case "delay":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "imaginal delay '%s' must be a number", field.Value.String())
				continue
			}

			if *field.Value.Number < 0 {
				errs.Addc(&field.Pos, "imaginal delay '%s' must be a positive number", field.Value.String())
				continue
			}

			imaginal.Delay = *field.Value.Number

		default:
			errs.Addc(&field.Pos, "unrecognized field '%s' in imaginal config", field.Key)
		}
	}
}

func addMemory(model *actr.Model, mem []*field, errs *errorListWithContext) {
	if mem == nil {
		return
	}

	for _, field := range mem {
		switch field.Key {
		case "latency":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory latency '%s' must be a number", field.Value.String())
				continue
			}

			model.Memory.Latency = field.Value.Number

		case "threshold":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory threshold '%s' must be a number", field.Value.String())
				continue
			}

			model.Memory.Threshold = field.Value.Number

		case "max_time":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory max_time '%s' must be a number", field.Value.String())
				continue
			}

			model.Memory.MaxTime = field.Value.Number

		case "finst_size":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory finst_size '%s' must be a number", field.Value.String())
				continue
			}

			size := int(*field.Value.Number)
			model.Memory.FinstSize = &size

		case "finst_time":
			if field.Value.Number == nil {
				errs.Addc(&field.Pos, "memory finst_time '%s' must be a number", field.Value.String())
				continue
			}

			model.Memory.FinstTime = field.Value.Number

		default:
			errs.Addc(&field.Pos, "unrecognized field '%s' in memory", field.Key)
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

		slotNames := []string{}
		for _, slot := range chunk.Slots {
			slotNames = append(slotNames, slot.Slot)
		}

		aChunk := actr.Chunk{
			Name:           chunk.Name,
			SlotNames:      slotNames,
			NumSlots:       len(chunk.Slots),
			AMODLineNumber: chunk.Pos.Line,
		}

		model.Chunks = append(model.Chunks, &aChunk)
	}
}

func addInit(model *actr.Model, init *initSection) (err error) {
	if init == nil {
		return
	}

	errs := errorListWithContext{}

	for _, initialization := range init.Initializations {

		err := validateInitialization(model, initialization)
		if err != nil {
			errs.AddErrorIfNotNil(err)
			continue
		}

		name := initialization.Name
		buffer := model.LookupBuffer(name)

		if buffer != nil {
			pattern, err := createChunkPattern(model, initialization.InitPattern)
			if err != nil {
				errs.AddErrorIfNotNil(err)
				continue
			}

			init := actr.Initializer{
				Buffer:         buffer,
				Pattern:        pattern,
				AMODLineNumber: init.Pos.Line,
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
					Memory:         model.Memory,
					Pattern:        pattern,
					AMODLineNumber: init.Pos.Line,
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
			Name:           production.Name,
			Description:    production.Description,
			AMODLineNumber: production.Pos.Line,
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
				Memory:  model.Memory,
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
			Memory:  model.Memory,
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
