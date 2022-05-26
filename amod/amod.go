package amod

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/issues"
)

var debugging bool = false

type ParseError struct{}

func (e ParseError) Error() string {
	return "Failed to parse amod file"
}

type CompileError struct{}

func (e CompileError) Error() string {
	return "Failed to compile amod file"
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
func GenerateModel(buffer string) (model *actr.Model, log *Log, err error) {
	r := strings.NewReader(buffer)

	log = &Log{issues.New()}

	amod, err := parse(r)
	if err != nil {
		pErr, ok := err.(participle.Error)
		if ok {
			location := issues.Location{
				Line:        pErr.Position().Line,
				ColumnStart: pErr.Position().Column,
				ColumnEnd:   pErr.Position().Column,
			}
			log.Error(&location, pErr.Message())
		} else {
			log.Error(&issues.Location{}, err.Error())
		}

		err = ParseError{}
		return
	}

	model, err = generateModel(amod, log)
	return
}

// GenerateModelFromFile generates a model from the file 'fileName'.
func GenerateModelFromFile(fileName string) (model *actr.Model, log *Log, err error) {
	log = &Log{issues.New()}

	amod, err := parseFile(fileName)
	if err != nil {
		pErr, ok := err.(participle.Error)
		if ok {
			location := issues.Location{
				Line:        pErr.Position().Line,
				ColumnStart: pErr.Position().Column,
				ColumnEnd:   pErr.Position().Column,
			}
			log.Error(&location, pErr.Message())
		} else {
			log.Error(&issues.Location{}, err.Error())
		}

		err = ParseError{}
		return
	}

	model, err = generateModel(amod, log)
	return
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

	log := Log{issues.New()}

	err := patternParser.Parse("", r, &p)
	if err != nil {
		pErr, ok := err.(participle.Error)
		if ok {
			err = errors.New(pErr.Message())
		}

		return nil, err
	}

	err = validatePattern(model, &log, &p)
	if err != nil {
		err = errors.New(log.FirstEntry())
		return nil, err
	}

	return createChunkPattern(model, &log, &p)
}

// generateModel runs through the parsed structures and creates an actr.Model from them
func generateModel(amod *amodFile, log *Log) (model *actr.Model, err error) {
	model = &actr.Model{
		Name:        amod.Model.Name,
		Description: amod.Model.Description,
		Authors:     amod.Model.Authors,
	}

	model.Initialize()

	addConfig(model, log, amod.Config)
	addExamples(model, log, amod.Model.Examples)
	addInit(model, log, amod.Init)
	addProductions(model, log, amod.Productions)

	if log.HasError() {
		return nil, CompileError{}
	}

	return
}

func addConfig(model *actr.Model, log *Log, config *configSection) {
	if config == nil {
		return
	}

	addGACTAR(model, log, config.GACTAR)
	addModules(model, log, config.Modules)
	addChunks(model, log, config.ChunkDecls)
}

func addExamples(model *actr.Model, log *Log, examples []*pattern) {
	if len(examples) == 0 {
		return
	}

	for _, example := range examples {
		err := validatePattern(model, log, example)
		if err != nil {
			continue
		}

		pattern, err := createChunkPattern(model, log, example)
		if err != nil {
			continue
		}

		model.Examples = append(model.Examples, pattern)
	}
}

func addGACTAR(model *actr.Model, log *Log, list []*field) {
	if list == nil {
		return
	}

	for _, field := range list {
		value := field.Value

		switch field.Key {
		case "log_level":
			if (value.Str == nil) || !actr.ValidLogLevel(*value.Str) {
				log.ErrorT(value.Tokens, "log_level '%s' must be 'min', 'info', 'or 'detail'", value.String())
				continue
			}

			model.LogLevel = actr.ACTRLogLevel(*value.Str)

		default:
			log.ErrorTR(field.Tokens, 0, 1, "unrecognized field in gactar section: '%s'", field.Key)
		}
	}
}

func addModules(model *actr.Model, log *Log, modules []*module) {
	if modules == nil {
		return
	}

	for _, module := range modules {
		switch module.Name {
		case "imaginal":
			addImaginal(model, log, module.InitFields)
		case "memory":
			addMemory(model, log, module.InitFields)
		default:
			log.ErrorT(module.Tokens, "unrecognized module in config: '%s'", module.Name)
		}
	}
}

func addImaginal(model *actr.Model, log *Log, fields []*field) {
	imaginal := model.CreateImaginal()

	for _, field := range fields {
		value := field.Value

		switch field.Key {
		case "delay":
			if value.Number == nil {
				log.ErrorT(value.Tokens, "imaginal delay '%s' must be a number", value.String())
				continue
			}

			if *value.Number < 0 {
				log.ErrorT(value.Tokens, "imaginal delay '%s' must be a positive number", value.String())
				continue
			}

			imaginal.Delay = *value.Number

		default:
			log.ErrorTR(field.Tokens, 0, 1, "unrecognized field '%s' in imaginal config", field.Key)
		}
	}
}

func addMemory(model *actr.Model, log *Log, mem []*field) {
	if mem == nil {
		return
	}

	for _, field := range mem {
		value := field.Value

		switch field.Key {
		case "latency":
			if value.Number == nil {
				log.ErrorT(value.Tokens, "memory latency '%s' must be a number", value.String())
				continue
			}

			model.Memory.Latency = value.Number

		case "threshold":
			if value.Number == nil {
				log.ErrorT(value.Tokens, "memory threshold '%s' must be a number", value.String())
				continue
			}

			model.Memory.Threshold = value.Number

		case "max_time":
			if field.Value.Number == nil {
				log.ErrorT(value.Tokens, "memory max_time '%s' must be a number", value.String())
				continue
			}

			model.Memory.MaxTime = value.Number

		case "finst_size":
			if value.Number == nil {
				log.ErrorT(value.Tokens, "memory finst_size '%s' must be a number", value.String())
				continue
			}

			size := int(*value.Number)
			model.Memory.FinstSize = &size

		case "finst_time":
			if value.Number == nil {
				log.ErrorT(value.Tokens, "memory finst_time '%s' must be a number", value.String())
				continue
			}

			model.Memory.FinstTime = value.Number

		default:
			log.ErrorTR(field.Tokens, 0, 1, "unrecognized field '%s' in memory", field.Key)
		}
	}
}

func addChunks(model *actr.Model, log *Log, chunks []*chunkDecl) {
	if chunks == nil {
		return
	}

	for _, chunk := range chunks {
		err := validateChunk(model, log, chunk)
		if err != nil {
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
			AMODLineNumber: chunk.Tokens[0].Pos.Line,
		}

		model.Chunks = append(model.Chunks, &aChunk)
	}
}

func addInit(model *actr.Model, log *Log, init *initSection) {
	if init == nil {
		return
	}

	for _, initialization := range init.Initializations {
		err := validateInitialization(model, log, initialization)
		if err != nil {
			continue
		}

		name := initialization.Name
		module := model.LookupModule(name)

		for _, init := range initialization.InitPatterns {
			pattern, err := createChunkPattern(model, log, init)
			if err != nil {
				continue
			}

			init := actr.Initializer{
				Buffer:         module,
				Pattern:        pattern,
				AMODLineNumber: init.Tokens[0].Pos.Line,
			}

			model.Initializers = append(model.Initializers, &init)
		}
	}
}

func addProductions(model *actr.Model, log *Log, productions *productionSection) {
	if productions == nil {
		return
	}

	for _, production := range productions.Productions {
		prod := actr.Production{
			Name:           production.Name,
			Description:    production.Description,
			VarIndexMap:    map[string]actr.VarIndex{},
			AMODLineNumber: production.Tokens[0].Pos.Line,
		}

		err := validateMatch(production.Match, model, log, &prod)
		if err != nil {
			continue
		}

		for _, item := range production.Match.Items {
			pattern, err := createChunkPattern(model, log, item.Pattern)
			if err != nil {
				continue
			}

			name := item.Name
			match := actr.Match{
				Buffer:  model.LookupBuffer(name),
				Pattern: pattern,
			}

			prod.Matches = append(prod.Matches, &match)

			for index, slot := range pattern.Slots {
				item := slot.Items[0]

				if item.Var == nil {
					continue
				}

				// Track the buffer and slot name the variable refers to
				varItem := *item.Var
				if _, ok := prod.VarIndexMap[varItem]; !ok {
					varIndex := actr.VarIndex{
						Var:      varItem,
						Buffer:   match.Buffer,
						SlotName: pattern.Chunk.SlotName(index),
					}
					prod.VarIndexMap[varItem] = varIndex
				}
			}
		}

		if production.Do.Statements != nil {
			validateDo(log, production)

			for _, statement := range *production.Do.Statements {
				addStatement(model, log, statement, &prod)
			}

			validateVariableUsage(log, production.Match, production.Do)
		}

		model.Productions = append(model.Productions, &prod)
	}
}

func createChunkPattern(model *actr.Model, log *Log, cp *pattern) (*actr.Pattern, error) {
	chunk := model.LookupChunk(cp.ChunkName)
	if chunk == nil {
		log.ErrorTR(cp.Tokens, 1, 2, "could not find chunk named '%s'", cp.ChunkName)
		return nil, CompileError{}
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

func addStatement(model *actr.Model, log *Log, statement *statement, production *actr.Production) (err error) {
	var s *actr.Statement

	if statement.Set != nil {
		s, err = addSetStatement(model, log, statement.Set, production)
	} else if statement.Recall != nil {
		s, err = addRecallStatement(model, log, statement.Recall, production)
	} else if statement.Clear != nil {
		s, err = addClearStatement(model, log, statement.Clear, production)
	} else if statement.Print != nil {
		s, err = addPrintStatement(model, log, statement.Print, production)
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

func addSetStatement(model *actr.Model, log *Log, set *setStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateSetStatement(set, model, log, production)
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
		bufferName := buffer.GetBufferName()

		// find slot index in chunk
		match := production.LookupMatchByBuffer(bufferName)

		s.Set.Chunk = match.Pattern.Chunk

		slotName := *set.Slot
		index := match.Pattern.Chunk.GetSlotIndex(slotName)
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
		pattern, err := createChunkPattern(model, log, set.Pattern)
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

func addRecallStatement(model *actr.Model, log *Log, recall *recallStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateRecallStatement(recall, model, log, production)
	if err != nil {
		return nil, err
	}

	pattern, err := createChunkPattern(model, log, recall.Pattern)
	if err != nil {
		return nil, err
	}

	s := actr.Statement{
		Recall: &actr.RecallStatement{
			Pattern:    pattern,
			MemoryName: model.Memory.GetModuleName(),
		},
	}

	return &s, nil
}

func addClearStatement(model *actr.Model, log *Log, clear *clearStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateClearStatement(clear, model, log, production)
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

func addPrintStatement(model *actr.Model, log *Log, print *printStatement, production *actr.Production) (*actr.Statement, error) {
	err := validatePrintStatement(print, model, log, production)
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
