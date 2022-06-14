package amod

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/actr/modules"
	"github.com/asmaloney/gactar/actr/params"

	"github.com/asmaloney/gactar/util/issues"
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

func (e *ErrParseChunk) Error() string {
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

	var p pattern

	r := strings.NewReader(chunk)

	log := newLog()

	err := patternParser.Parse("", r, &p)
	if err != nil {
		pErr, ok := err.(participle.Error)
		if ok {
			err = &ErrParseChunk{Message: pErr.Message()}
		}

		return nil, err
	}

	err = validatePattern(model, log, &p)
	if err != nil {
		err = &ErrParseChunk{Message: log.FirstEntry()}
		return nil, err
	}

	return createChunkPattern(model, log, &p)
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

	addGACTAR(model, log, config.GACTAR)
	addModules(model, log, config.Modules)
	addChunks(model, log, config.ChunkDecls)
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

func addGACTAR(model *actr.Model, log *issueLog, list []*field) {
	if list == nil {
		return
	}

	for _, field := range list {
		value := field.Value

		options, err := model.SetParam(&params.Param{
			Key: field.Key,
			Value: params.Value{
				ID:     value.ID,
				Str:    value.Str,
				Number: value.Number,
			},
		})

		if err != params.NoError {
			switch err {
			case params.InvalidOption:
				log.errorT(value.Tokens, "%s ('%s') must be one of %q", field.Key, value.String(), strings.Join(options, ", "))
				continue

			case params.UnrecognizedParam:
				log.errorTR(field.Tokens, 0, 1, "unrecognized field in gactar section: '%s'", field.Key)
				continue

			default:
				log.errorT(field.Tokens, "internal: unhandled error (%d) in gactar section: '%s'", err, field.Key)
				continue
			}
		}
	}
}

func addModules(model *actr.Model, log *issueLog, modules []*module) {
	if modules == nil {
		return
	}

	for _, module := range modules {
		switch module.Name {
		case "goal":
			addGoal(model, log, module.InitFields)
		case "imaginal":
			addImaginal(model, log, module.InitFields)
		case "memory":
			addMemory(model, log, module.InitFields)
		case "procedural":
			addProcedural(model, log, module.InitFields)
		default:
			log.errorT(module.Tokens, "unrecognized module in config: '%s'", module.Name)
		}
	}
}

func setModuleParams(module modules.ModuleInterface, log *issueLog, fields []*field) {
	if len(fields) == 0 {
		return
	}

	moduleName := module.ModuleName()

	for _, field := range fields {
		value := field.Value

		err := module.SetParam(&params.Param{
			Key: field.Key,
			Value: params.Value{
				ID:     value.ID,
				Str:    value.Str,
				Number: value.Number,
			},
		})

		if err != params.NoError {
			switch err {
			case params.NumberRequired:
				log.errorT(value.Tokens, "%s %s '%s' must be a number", moduleName, field.Key, value.String())
				continue

			case params.NumberMustBePositive:
				log.errorT(value.Tokens, "%s %s '%s' must be a positive number", moduleName, field.Key, value.String())
				continue

			case params.UnrecognizedParam:
				log.errorTR(field.Tokens, 0, 1, "unrecognized field '%s' in %s config", field.Key, moduleName)
				continue

			default:
				log.errorT(field.Tokens, "internal: unhandled error (%d) in %s config: '%s'", err, moduleName, field.Key)
				continue
			}
		}
	}
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

func addChunks(model *actr.Model, log *issueLog, chunks []*chunkDecl) {
	if chunks == nil {
		return
	}

	for _, chunk := range chunks {
		err := validateChunk(model, log, chunk)
		if err != nil {
			continue
		}

		aChunk := actr.Chunk{
			Name:           chunk.Name,
			SlotNames:      chunk.Slots,
			NumSlots:       len(chunk.Slots),
			AMODLineNumber: chunk.Tokens[0].Pos.Line,
		}

		model.Chunks = append(model.Chunks, &aChunk)
	}
}

func addInit(model *actr.Model, log *issueLog, init *initSection) {
	if init == nil {
		return
	}

	for _, initialization := range init.Initializations {
		err := validateInitialization(model, log, initialization)
		if err != nil {
			continue
		}

		name := initialization.Name
		moduleInterface := model.LookupModule(name)

		for _, init := range initialization.InitPatterns {
			pattern, err := createChunkPattern(model, log, init)
			if err != nil {
				continue
			}

			init := actr.Initializer{
				Module:         moduleInterface,
				Pattern:        pattern,
				AMODLineNumber: init.Tokens[0].Pos.Line,
			}

			model.Initializers = append(model.Initializers, &init)
		}
	}
}

func addProductions(model *actr.Model, log *issueLog, productions *productionSection) {
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

		for _, match := range production.Match.Items {
			pattern, err := createChunkPattern(model, log, match.Pattern)
			if err != nil {
				continue
			}

			name := match.Name
			actrMatch := actr.Match{
				Buffer:  model.LookupBuffer(name),
				Pattern: pattern,
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
						Buffer:   actrMatch.Buffer,
						SlotName: pattern.Chunk.SlotName(index),
					}
					prod.VarIndexMap[name] = varIndex
				}
			}

			if match.When != nil {
				for _, expr := range *match.When.Expressions {
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
		s, err = addSetStatement(model, log, statement.Set, production)

	case statement.Recall != nil:
		s, err = addRecallStatement(model, log, statement.Recall, production)

	case statement.Clear != nil:
		s, err = addClearStatement(model, log, statement.Clear, production)

	case statement.Print != nil:
		s, err = addPrintStatement(model, log, statement.Print, production)

	case statement.Stop != nil:
		s, err = addStopStatement(model, log, statement.Stop, production)

	default:
		return ErrStatementNotHandled
	}

	if err != nil {
		return err
	}

	if s != nil {
		production.DoStatements = append(production.DoStatements, s)
	}

	return nil
}

func addSetStatement(model *actr.Model, log *issueLog, set *setStatement, production *actr.Production) (*actr.Statement, error) {
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
		bufferName := buffer.BufferName()

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

func addRecallStatement(model *actr.Model, log *issueLog, recall *recallStatement, production *actr.Production) (*actr.Statement, error) {
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
			MemoryName: model.Memory.ModuleName(),
		},
	}

	return &s, nil
}

func addClearStatement(model *actr.Model, log *issueLog, clear *clearStatement, production *actr.Production) (*actr.Statement, error) {
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

func addPrintStatement(model *actr.Model, log *issueLog, print *printStatement, production *actr.Production) (*actr.Statement, error) {
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

func addStopStatement(model *actr.Model, log *issueLog, stop *stopStatement, production *actr.Production) (*actr.Statement, error) {
	return &actr.Statement{
		Stop: &actr.StopStatement{},
	}, nil
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
