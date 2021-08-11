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

	return
}

func addConfig(model *actr.Model, config *configSection) (err error) {
	if config == nil {
		return
	}

	errs := errorListWithContext{}

	addACTR(model, config.ACTR, &errs)
	addBuffers(model, config.Buffers, &errs)
	addMemories(model, config.Memories, &errs)
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
			errs.Addc(&field.Pos, "Unrecognized field in actr section: '%s'", field.Key)
		}
	}
}

func addBuffers(model *actr.Model, list *identList, errs *errorListWithContext) {
	if list == nil {
		return
	}

	for _, name := range list.Identifiers {
		buffer := actr.Buffer{
			Name: name,
		}

		model.Buffers = append(model.Buffers, &buffer)
	}
}

func addMemories(model *actr.Model, list *memoryList, errs *errorListWithContext) {
	if list == nil {
		return
	}

	for _, mem := range list.Memories {
		memory := actr.Memory{
			Name: mem.Name,
		}

		for _, field := range mem.Fields.Fields {
			switch field.Key {
			case "buffer":
				if field.Value.Number != nil {
					errs.Addc(&field.Pos, "buffer '%v' should not be a number in memory '%s'", *field.Value.Number, mem.Name)
					continue
				}

				bufferName := field.Value.String

				buffer := model.LookupBuffer(*bufferName)
				if buffer == nil {
					errs.Addc(&field.Pos, "buffer '%s' not found for memory '%s'", *bufferName, mem.Name)
					continue
				} else {
					memory.Buffer = buffer
				}

			case "latency":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "latency '%s' should not be a string in memory '%s'", *field.Value.String, mem.Name)
					continue
				}

				memory.Latency = field.Value.Number

			case "threshold":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "threshold '%s' should not be a string in memory '%s'", *field.Value.String, mem.Name)
					continue
				}

				memory.Threshold = field.Value.Number

			case "max_time":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "max_time '%s' should not be a string in memory '%s'", *field.Value.String, mem.Name)
					continue
				}

				memory.MaxTime = field.Value.Number

			case "finst_size":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "finst_size '%s' should not be a string in memory '%s'", *field.Value.String, mem.Name)
					continue
				}

				size := int(*field.Value.Number)
				memory.FinstSize = &size

			case "finst_time":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "finst_time '%s' should not be a string in memory '%s'", *field.Value.String, mem.Name)
					continue
				}

				memory.FinstTime = field.Value.Number

			default:
				errs.Addc(&field.Pos, "Unrecognized field '%s' in memory '%s'", field.Key, memory.Name)
			}
		}

		model.Memories = append(model.Memories, &memory)
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
		name := initializer.Name
		memory := model.LookupMemory(name)
		if memory == nil {
			errs.Addc(&initializer.Pos, "memory not found for initialization '%s'", name)
			continue
		}

		if initializer.Items == nil {
			errs.Addc(&initializer.Pos, "no memory initializers for memory '%s'", name)
			continue
		}

		for _, item := range initializer.Items.Strings {
			init := actr.Initializer{
				MemoryName: memory.Name,
				Text:       item,
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

			exists := model.BufferOrMemoryExists(name)
			if !exists {
				errs.Addc(&item.Pos, "buffer or memory '%s' not found in production '%s'", name, prod.Name)
				continue
			}

			if item.Pattern != nil {
				pattern := createChunkPattern(item.Pattern)

				prod.Matches = append(prod.Matches, &actr.Match{
					Name:    name,
					Pattern: pattern,
				})
			} else if item.Text != nil {
				prod.Matches = append(prod.Matches, &actr.Match{
					Name: name,
					Text: item.Text,
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

func createChunkPattern(cp *pattern) *actr.Pattern {
	pattern := actr.Pattern{}
	for _, f := range cp.Fields {

		field := actr.PatternField{}

		if f.Name != nil {
			field.Name = f.Name
		}

		for _, f := range f.Items {
			if f.ID != nil {
				field.AddItem(&actr.PatternFieldItem{ID: f.ID})
			} else if f.Num != nil {
				field.AddItem(&actr.PatternFieldItem{Num: f.Num})
			} else if f.Var != nil {
				field.AddItem(&actr.PatternFieldItem{Var: f.Var})
			} else if f.NotVar != nil {
				field.AddItem(&actr.PatternFieldItem{ID: f.NotVar, Negated: true})
			} else if f.OptionalVar != nil {
				field.AddItem(&actr.PatternFieldItem{ID: f.OptionalVar, Optional: true})
			} else if f.NotOptionalVar != nil {
				field.AddItem(&actr.PatternFieldItem{ID: f.NotOptionalVar, Negated: true, Optional: true})
			}
		}
		pattern.AddField(&field)
	}

	return &pattern
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

	s := actr.Statement{
		Set: &actr.SetStatement{
			BufferName: set.BufferName,
		},
	}

	if set.Field != nil {
		field := actr.SetField{}
		if set.Field.ArgNum != nil {
			argNum := int(*set.Field.ArgNum)
			field.ArgNum = &argNum
		} else if set.Field.Name != nil {
			field.FieldName = set.Field.Name
		}

		s.Set.Field = &field
	}

	if set.Pattern != nil {
		pattern := createChunkPattern(set.Pattern)
		s.Set.Pattern = pattern
	} else if set.Arg != nil {
		arg := set.Arg
		s.Set.Text = &arg.Arg
	}

	return &s, nil
}

func addRecallStatement(model *actr.Model, recall *recallStatement, production *actr.Production) (*actr.Statement, error) {
	err := validateRecallStatement(recall, model, production)
	if err != nil {
		return nil, err
	}

	pattern := createChunkPattern(recall.Pattern)

	s := actr.Statement{
		Recall: &actr.RecallStatement{
			Pattern:    pattern,
			MemoryName: recall.MemoryName,
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
