// Package amod implements parsing of amod files into actr.Models.
package amod

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/alecthomas/participle/v2"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/modules"
	"github.com/asmaloney/gactar/actr/param"

	"github.com/asmaloney/gactar/util/issues"
	"github.com/asmaloney/gactar/util/keyvalue"
)

var (
	debugging bool = false

	ErrParse               = errors.New("failed to parse amod file")
	ErrCompile             = errors.New("failed to compile amod file")
	ErrStatementNotHandled = errors.New("statement type not handled")
)

type ErrParseChunk struct {
	Message string
}

func (e ErrParseChunk) Error() string {
	return fmt.Sprintf("cannot parse chunk: %s", e.Message)
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
func GenerateModel(buffer string) (model *actr.Model, iLog *issues.Log, err error) {
	r := strings.NewReader(buffer)

	log := newLog()
	iLog = &log.Log

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

		err = ErrParse
		return
	}

	model, err = generateModel(amod, log)
	if err != nil {
		return
	}

	model.FinalizeImplicitChunks()
	return
}

// GenerateModelFromFile generates a model from the file 'fileName'.
func GenerateModelFromFile(fileName string) (model *actr.Model, iLog *issues.Log, err error) {
	log := newLog()
	iLog = &log.Log

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

		err = ErrParse
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

	log := newLog()

	patternParser, err := participle.ParserForProduction[pattern](amodParser)
	if err != nil {
		err = &ErrParseChunk{Message: err.Error()}
		return nil, err
	}

	p, err := patternParser.ParseString("", chunk)
	if err != nil {
		pErr, ok := err.(participle.Error)
		if ok {
			err = &ErrParseChunk{Message: pErr.Message()}
		}

		return nil, err
	}

	err = validatePattern(model, log, p)
	if err != nil {
		err = &ErrParseChunk{Message: log.FirstEntry()}
		return nil, err
	}

	return createChunkPattern(model, log, p)
}

// generateModel runs through the parsed structures and creates an actr.Model from them
func generateModel(amod *amodFile, log *issueLog) (model *actr.Model, err error) {
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
		return nil, ErrCompile
	}

	return
}

func addConfig(model *actr.Model, log *issueLog, config *configSection) {
	if config == nil {
		return
	}

	addGACTAR(model, log, config.GactarConfig)
	addModules(model, log, config.ModuleConfig)
	addChunks(model, log, config.ChunkConfig)
}

func addExamples(model *actr.Model, log *issueLog, examples []*pattern) {
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

func addGACTAR(model *actr.Model, log *issueLog, config *gactarConfig) {
	if config == nil {
		return
	}

	list := config.GactarFields
	if len(list) == 0 {
		return
	}

	for _, field := range list {
		value := field.Value

		kv := fieldToKeyValue(field)
		err := model.SetParam(kv)
		if err != nil {
			switch {
			// field errors
			case errors.As(err, &param.ErrUnrecognizedOption{}):
				log.errorTR(field.Tokens, 0, 1, "%v in gactar section", err)
				continue

				// value errors
			case errors.As(err, &keyvalue.ErrInvalidType{}) ||
				errors.As(err, &param.ErrInvalidType{}) ||
				errors.As(err, &param.ErrInvalidValue{}):
				log.errorTR(value.Tokens, 1, 1, "'%s' %v", field.Key, err)
				continue

			default:
				log.errorT(field.Tokens, "INTERNAL: unhandled error (%v) in gactar config: '%s'", err, field.Key)
				continue
			}
		}
	}
}

func addModules(model *actr.Model, log *issueLog, config *moduleConfig) {
	if config == nil {
		return
	}

	modules := config.Modules
	if len(modules) == 0 {
		return
	}

	for _, module := range modules {
		_ = validateFieldList(log, module.Fields)

		switch module.ModuleName {
		case "extra_buffers":
			addExtraBuffers(model, log, module.Fields)
		case "goal":
			addGoal(model, log, module.Fields)
		case "imaginal":
			addImaginal(model, log, module.Fields)
		case "memory":
			addMemory(model, log, module.Fields)
		case "procedural":
			addProcedural(model, log, module.Fields)
		default:
			log.errorT(module.Tokens, "unrecognized module in config: '%s'", module.ModuleName)
		}
	}

	_ = validateInterModuleInitDependencies(model, log, config)
}

func setBufferParams(moduleName string, buffer buffer.Interface, log *issueLog, field *field) {
	if field == nil {
		return
	}

	bufferName := buffer.Name()

	_ = validateFieldList(log, field.Value.Fields)

	for _, f := range field.Value.Fields {
		value := field.Value

		kv := fieldToKeyValue(f)
		err := buffer.SetParam(kv)

		if err != nil {
			switch {
			// field errors
			case errors.As(err, &param.ErrUnrecognizedOption{}):
				log.errorTR(field.Tokens, 0, 1, "%v in %s (%s) config", err, moduleName, bufferName)
				continue

			// value errors
			case errors.As(err, &keyvalue.ErrInvalidType{}) ||
				errors.As(err, &param.ErrInvalidType{}) ||
				errors.As(err, &param.ErrInvalidValue{}) ||
				errors.As(err, &param.ErrValueOutOfRange{}):
				log.errorTR(value.Tokens, 1, 1, "%s %q %v", moduleName, kv.Key, err)
				continue

			default:
				log.errorT(field.Tokens, "INTERNAL: unhandled error (%v) in %s (%s) config: %q", err, moduleName, bufferName, kv.Key)
				continue
			}
		}
	}
}

func setModuleParams(module modules.Interface, log *issueLog, fields []*field) {
	if len(fields) == 0 {
		return
	}

	moduleName := module.ModuleName()

	for _, field := range fields {
		value := field.Value

		kv := fieldToKeyValue(field)

		// if the key is a buffer name, init the buffer
		buffer := module.Buffers().Lookup(kv.Key)
		if buffer != nil {
			setBufferParams(moduleName, buffer, log, field)
		} else {
			// save our current buffers so we can determine if any have been created in SetParam
			saveBufferList := module.Buffers()

			err := module.SetParam(kv)

			if err != nil {
				switch {
				// field errors
				case errors.As(err, &param.ErrUnrecognizedOption{}):
					log.errorTR(field.Tokens, 0, 1, "%v in %s config", err, moduleName)
					continue

				// value errors
				case errors.As(err, &keyvalue.ErrInvalidType{}) ||
					errors.As(err, &param.ErrInvalidType{}) ||
					errors.As(err, &param.ErrInvalidValue{}) ||
					errors.As(err, &param.ErrValueOutOfRange{}):
					log.errorTR(value.Tokens, 1, 1, "%s %q %v", moduleName, field.Key, err)
					continue

				default:
					log.errorT(field.Tokens, "INTERNAL: unhandled error (%v) in %s config: %q", err, moduleName, field.Key)
					continue
				}
			}

			// check if we created any buffers through the params (e.g. extra_buffers) and set their params
			newBufferList := module.Buffers()

			if len(saveBufferList) != len(newBufferList) {
				newBuffers := newBufferList[len(saveBufferList):]
				for _, buffer := range newBuffers {
					setBufferParams(moduleName, buffer, log, field)
				}
			}
		}
	}
}

func addExtraBuffers(model *actr.Model, log *issueLog, fields []*field) {
	eb := model.CreateExtraBuffers()

	setModuleParams(eb, log, fields)
}

func addGoal(model *actr.Model, log *issueLog, fields []*field) {
	setModuleParams(model.Goal, log, fields)
}

func addImaginal(model *actr.Model, log *issueLog, fields []*field) {
	imaginal := model.CreateImaginal()

	setModuleParams(imaginal, log, fields)
}

func addMemory(model *actr.Model, log *issueLog, fields []*field) {
	setModuleParams(model.Memory, log, fields)
}

func addProcedural(model *actr.Model, log *issueLog, fields []*field) {
	setModuleParams(model.Procedural, log, fields)
}

func addChunks(model *actr.Model, log *issueLog, config *chunkConfig) {
	if config == nil {
		return
	}

	chunks := config.ChunkDecls
	if len(chunks) == 0 {
		return
	}

	for _, chunk := range chunks {
		err := validateChunk(model, log, chunk)
		if err != nil {
			continue
		}

		aChunk := actr.Chunk{
			TypeName:       chunk.TypeName,
			SlotNames:      chunk.Slots,
			NumSlots:       len(chunk.Slots),
			AMODLineNumber: chunk.Tokens[0].Pos.Line,
		}

		model.Chunks = append(model.Chunks, &aChunk)
	}
}

func addInitializers(model *actr.Model, log *issueLog, module modules.Interface, buffer buffer.Interface, init *namedInitializer) {
	// Check for duplicate initializer names.
	// Note that this can't be checked in validateInitialization because model.ExplicitChunks is not filled in yet.
	if init.ChunkName != nil && slices.Contains(model.ExplicitChunks, *init.ChunkName) {
		log.errorTR(init.Tokens, 0, 1, "duplicate chunk name %q found in initialization", *init.ChunkName)
		return
	}

	actrPattern, err := createChunkPattern(model, log, init.Pattern)
	if err != nil {
		return
	}

	model.AddInitializer(
		&actr.Initializer{
			Module:         module,
			Buffer:         buffer,
			ChunkName:      init.ChunkName,
			Pattern:        actrPattern,
			AMODLineNumber: init.Tokens[0].Pos.Line,
		},
	)
}

func addInit(model *actr.Model, log *issueLog, init *initSection) {
	if init == nil {
		return
	}

	for _, initialization := range init.Initializations {
		if initialization.ModuleInitializer != nil {
			moduleInitializer := initialization.ModuleInitializer
			err := validateModuleInitialization(model, log, moduleInitializer)
			if err != nil {
				continue
			}

			name := moduleInitializer.ModuleName
			moduleInterface := model.LookupModule(name)

			if len(moduleInitializer.InitPatterns) > 0 {
				for _, initPattern := range moduleInitializer.InitPatterns {
					addInitializers(model, log, moduleInterface, moduleInterface.Buffers().At(0), initPattern)
				}
			} else if len(moduleInitializer.BufferInitPatterns) > 0 {
				for _, bufferInit := range moduleInitializer.BufferInitPatterns {
					buff := model.LookupBuffer(bufferInit.BufferName)

					for _, initPattern := range bufferInit.InitPatterns {
						addInitializers(model, log, moduleInterface, buff, initPattern)
					}
				}
			}
		} else if initialization.SimilarityInitializer != nil {
			partialInitializer := initialization.SimilarityInitializer

			for _, similar := range partialInitializer.SimilarList {
				actrSimilar := &actr.Similarity{
					ChunkOne:       similar.ChunkOne,
					ChunkTwo:       similar.ChunkTwo,
					Value:          similar.Value,
					AMODLineNumber: similar.Tokens[0].Pos.Line,
				}

				model.AddSimilarity(actrSimilar)
			}
		}
	}
}

func addProductions(model *actr.Model, log *issueLog, productions *productionSection) {
	if productions == nil {
		return
	}

	for _, production := range productions.Productions {
		prod := actr.Production{
			Model:          model,
			Name:           production.Name,
			Description:    production.Description,
			VarIndexMap:    map[string]actr.VarIndex{},
			AMODLineNumber: production.Tokens[0].Pos.Line,
		}

		err := validateMatch(production.Match, model, log, &prod)
		if err != nil {
			continue
		}

		for _, match := range production.Match.Items {
			switch {
			case match.BufferPattern != nil:
				pattern, err := createChunkPattern(model, log, match.BufferPattern.Pattern)
				if err != nil {
					continue
				}

				name := match.BufferPattern.BufferName
				buffer := model.LookupBuffer(name)
				actrMatch := actr.Match{
					BufferPattern: &actr.BufferPatternMatch{
						Buffer:  buffer,
						Pattern: pattern,
					},
				}

				prod.Matches = append(prod.Matches, &actrMatch)

				for index, slot := range pattern.Slots {
					if slot.Var == nil {
						continue
					}

					// Track the buffer and slot name the variable refers to
					varItem := slot.Var
					name := *slot.Var.Name
					if _, ok := prod.VarIndexMap[name]; !ok {
						varIndex := actr.VarIndex{
							Var:      varItem,
							Buffer:   buffer,
							SlotName: pattern.Chunk.SlotName(index),
						}
						prod.VarIndexMap[name] = varIndex
					}
				}

				if match.BufferPattern.When != nil {
					for _, expr := range *match.BufferPattern.When.Expressions {
						comparison := actr.Equal

						if expr.Comparison.NotEqual != nil {
							comparison = actr.NotEqual
						}

						actrConstraint := actr.Constraint{
							LHS:        &expr.LHS,
							Comparison: comparison,
							RHS:        convertArg(expr.RHS),
						}

						// Add the constraint on the pattern var
						patternVar, ok := prod.VarIndexMap[expr.LHS]
						if !ok {
							// This is an error, but it is captured in validateVariableUsage() below
							continue
						}

						patternVar.Var.Constraints = append(patternVar.Var.Constraints, &actrConstraint)
					}
				}

			case match.BufferState != nil:
				name := match.BufferState.BufferName
				actrMatch := &actr.BufferStateMatch{
					Buffer: model.LookupBuffer(name),
					State:  match.BufferState.State,
				}

				// if we have a module state match already, add this buffer state match there
				match := findModuleStateMatch(&prod, name)

				if match != nil {
					match.BufferState = actrMatch
				} else {
					prod.Matches = append(prod.Matches, &actr.Match{
						BufferState: actrMatch,
					})
				}

			case match.ModuleState != nil:
				name := match.ModuleState.ModuleName
				module := model.LookupModule(name)

				// The generated code for the frameworks actually uses a buffer name, not the module name.
				// So store (one) here for convenience. If the module has multiple buffers it should not
				// matter which one we pick as the requests should be on its module.
				buffer := module.Buffers().At(0)

				actrMatch := &actr.ModuleStateMatch{
					Module: module,
					Buffer: buffer,
					State:  match.ModuleState.State,
				}

				// if we have a buffer state match already, add this module state match there
				match := findBufferStateMatch(&prod, buffer.Name())

				if match != nil {
					match.ModuleState = actrMatch
				} else {
					prod.Matches = append(prod.Matches, &actr.Match{
						ModuleState: actrMatch,
					})
				}
			}
		}

		validateDo(log, production)

		for _, statement := range *production.Do.Statements {
			err := addStatement(model, log, statement, &prod)
			if err != nil && !errors.Is(err, ErrCompile) {
				log.Error(nil, err.Error())
			}
		}

		validateVariableUsage(log, production.Match, production.Do)

		model.Productions = append(model.Productions, &prod)
	}
}

func findBufferStateMatch(prod *actr.Production, bufferName string) *actr.Match {
	for _, match := range prod.Matches {
		if match.BufferState == nil {
			continue
		}

		if match.BufferState.Buffer.Name() == bufferName {
			return match
		}
	}

	return nil
}

func findModuleStateMatch(prod *actr.Production, bufferName string) *actr.Match {
	for _, match := range prod.Matches {
		if match.ModuleState == nil {
			continue
		}

		if match.ModuleState.Buffer.Name() == bufferName {
			return match
		}
	}

	return nil
}

func argToKeyValue(key string, a *arg) *keyvalue.KeyValue {
	value := keyvalue.Value{}

	switch {
	case a.Nil != nil:
		nilStr := "nil"
		value.Str = &nilStr
	case a.Var != nil:
		value.Str = a.Var
	case a.ID != nil:
		value.Str = a.ID
	case a.Str != nil:
		value.Str = a.Str
	case a.Number != nil:
		num, _ := strconv.ParseFloat(*a.Number, 64)
		value.Number = &num
	}

	return &keyvalue.KeyValue{
		Key:   key,
		Value: value,
	}
}

func fieldToKeyValue(f *field) *keyvalue.KeyValue {
	value := f.Value

	if f.Value.OpenBrace != nil {

		fields := make([]keyvalue.KeyValue, len(value.Fields))

		for i, field := range value.Fields {
			kv := fieldToKeyValue(field)
			fields[i] = *kv
		}

		return &keyvalue.KeyValue{
			Key:   f.Key,
			Value: keyvalue.Value{Fields: &fields},
		}
	}

	return &keyvalue.KeyValue{
		Key: f.Key,
		Value: keyvalue.Value{
			ID:     value.ID,
			Str:    value.Str,
			Number: value.Number,
		},
	}
}

func createChunkPattern(model *actr.Model, log *issueLog, cp *pattern) (*actr.Pattern, error) {
	chunk := model.LookupChunk(cp.ChunkName)
	if chunk == nil {
		log.errorTR(cp.Tokens, 1, 2, "could not find chunk named '%s'", cp.ChunkName)
		return nil, ErrCompile
	}

	pattern := actr.Pattern{
		Chunk: chunk,
	}

	for _, slot := range cp.Slots {
		actrSlot := actr.PatternSlot{
			Negated: slot.Not,
		}

		switch {
		case slot.Wildcard != nil:
			actrSlot.Wildcard = true

		case slot.Nil != nil:
			actrSlot.Nil = true

		case slot.ID != nil:
			actrSlot.ID = slot.ID

		case slot.Str != nil:
			actrSlot.Str = slot.Str

		case slot.Num != nil:
			actrSlot.Num = slot.Num

		case slot.Var != nil:
			actrSlot.Var = &actr.PatternVar{Name: slot.Var}
		}

		pattern.AddSlot(&actrSlot)
	}

	return &pattern, nil
}

func addStatement(model *actr.Model, log *issueLog, statement *statement, production *actr.Production) (err error) {
	var s *actr.Statement

	switch {
	case statement.Set != nil:
		s, err = createSetStatement(model, log, statement.Set, production)

	case statement.Recall != nil:
		s, err = createRecallStatement(model, log, statement.Recall, production)

	case statement.Clear != nil:
		s, err = createClearStatement(model, log, statement.Clear, production)

	case statement.Print != nil:
		s, err = createPrintStatement(model, log, statement.Print, production)

	case statement.Stop != nil:
		s = createStopStatement()

	default:
		return ErrStatementNotHandled
	}

	if err != nil {
		return err
	}

	production.AddDoStatement(s)

	return nil
}

func createSetStatement(model *actr.Model, log *issueLog, set *setStatement, production *actr.Production) (*actr.Statement, error) {
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
		bufferName := buffer.Name()

		// find slot index in chunk
		match := production.LookupMatchByBuffer(bufferName)

		s.Set.Chunk = match.Pattern.Chunk

		slotName := *set.Slot
		index := match.Pattern.Chunk.SlotIndex(slotName)
		value := &actr.Value{}

		switch {
		case set.Value.Nil != nil:
			value.Nil = set.Value.Nil

		case set.Value.Var != nil:
			varName := strings.TrimPrefix(*set.Value.Var, "?")
			value.Var = &varName

		case set.Value.ID != nil:
			value.ID = set.Value.ID

		case set.Value.Number != nil:
			value.Number = set.Value.Number

		case set.Value.Str != nil:
			value.Str = set.Value.Str

		}

		newSlot := &actr.SetSlot{
			Name:      *set.Slot,
			SlotIndex: index,
			Value:     value,
		}

		production.AddSlotToSetStatement(s.Set, newSlot)
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

func createRecallStatement(model *actr.Model, log *issueLog, recall *recallStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateRecallStatement(recall, model, log, production)
	if err != nil {
		return nil, err
	}

	pattern, err := createChunkPattern(model, log, recall.Pattern)
	if err != nil {
		return nil, err
	}

	requestParameters := make(map[string]string)

	if recall.With != nil {
		for _, param := range *recall.With.Expressions {
			value := convertArg(param.Value)
			requestParameters[param.Param] = value.String()
		}
	}

	s := actr.Statement{
		Recall: &actr.RecallStatement{
			Pattern:           pattern,
			MemoryModuleName:  model.Memory.ModuleName(),
			RequestParameters: requestParameters,
		},
	}

	return &s, nil
}

func createClearStatement(model *actr.Model, log *issueLog, clear *clearStatement, production *actr.Production) (*actr.Statement, error) {
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

func createPrintStatement(model *actr.Model, log *issueLog, print *printStatement, production *actr.Production) (*actr.Statement, error) {
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

func createStopStatement() *actr.Statement {
	return &actr.Statement{Stop: &actr.StopStatement{}}
}

func convertArg(v *arg) (actrValue *actr.Value) {
	actrValue = &actr.Value{}

	switch {
	case v.Nil != nil:
		actrValue.Nil = v.Nil

	case v.Var != nil:
		actrValue.Var = v.Var

	case v.ID != nil:
		actrValue.ID = v.ID

	case v.Str != nil:
		actrValue.Str = v.Str

	case v.Number != nil:
		actrValue.Number = v.Number
	}

	return
}

func convertArgs(args []*arg) *[]*actr.Value {
	actrValues := []*actr.Value{}

	for _, v := range args {
		newValue := convertArg(v)

		actrValues = append(actrValues, newValue)
	}

	return &actrValues
}
