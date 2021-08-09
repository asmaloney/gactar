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
					errs.Addc(&field.Pos, "buffer should not be a number in memory '%s': %v", mem.Name, *field.Value.Number)
					continue
				}

				bufferName := field.Value.String

				buffer := model.LookupBuffer(*bufferName)
				if buffer == nil {
					errs.Addc(&field.Pos, "buffer not found for memory '%s': %s", mem.Name, *bufferName)
					continue
				} else {
					memory.Buffer = buffer
				}

			case "latency":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "latency should not be a string in memory '%s': %v", mem.Name, *field.Value.String)
					continue
				}

				memory.Latency = field.Value.Number

			case "threshold":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "threshold should not be a string in memory '%s': %v", mem.Name, *field.Value.String)
					continue
				}

				memory.Threshold = field.Value.Number

			case "max_time":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "max_time should not be a string in memory '%s': %v", mem.Name, *field.Value.String)
					continue
				}

				memory.MaxTime = field.Value.Number

			case "finst_size":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "finst_size should not be a string in memory '%s': %v", mem.Name, *field.Value.String)
					continue
				}

				size := int(*field.Value.Number)
				memory.FinstSize = &size

			case "finst_time":
				if field.Value.String != nil {
					errs.Addc(&field.Pos, "finst_time should not be a string in memory '%s': %v", mem.Name, *field.Value.String)
					continue
				}

				memory.FinstTime = field.Value.Number

			default:
				errs.Addc(&field.Pos, "Unrecognized field in memory '%s': '%s'", memory.Name, field.Key)
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
				errs.Addc(&item.Pos, "buffer or memory not found for production '%s': %s", prod.Name, name)
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
	for _, item := range cp.Items {
		if item.ID != nil {
			pattern.AddID(item.ID)
		} else if item.Var != nil {
			pattern.AddVar(item.Var)
		} else if item.Num != nil {
			pattern.AddNum(item.Num)
		} else if item.Field != nil {
			field := actr.PatternField{}

			if item.Field.Field != nil {
				field.Name = item.Field.Field
			}

			for _, f := range item.Field.Items {
				if f.ID != nil {
					field.Items = append(field.Items, actr.PatternFieldItem{ID: f.ID})
				} else if f.Num != nil {
					field.Items = append(field.Items, actr.PatternFieldItem{Num: f.Num})
				} else if f.NotID != nil {
					field.Items = append(field.Items, actr.PatternFieldItem{ID: f.NotID, Negated: true})
				} else if f.OptionalID != nil {
					field.Items = append(field.Items, actr.PatternFieldItem{ID: f.OptionalID, Optional: true})
				} else if f.NotOptionalID != nil {
					field.Items = append(field.Items, actr.PatternFieldItem{ID: f.NotOptionalID, Negated: true, Optional: true})
				}
			}
			pattern.AddField(&field)
		}
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
	errs := errorListWithContext{}

	name := set.BufferName
	buffer := model.LookupBuffer(name)
	if buffer == nil {
		errs.Addc(&set.Pos, "buffer not found in production '%s': '%s'", production.Name, name)
		return nil, errs
	}

	s := actr.Statement{
		Set: &actr.SetStatement{
			BufferName: name,
		},
	}

	if set.Pattern != nil {
		pattern := createChunkPattern(set.Pattern)
		s.Set.Pattern = pattern
	} else if set.Arg != nil {
		arg := set.Arg
		s.Set.Text = &arg.Arg
	}

	if set.Field != nil {
		if set.Field.ArgNum != nil {
			argNum := int(*set.Field.ArgNum)
			s.Set.Field = &actr.SetField{
				ArgNum: &argNum,
			}
		} else if set.Field.Name != nil {
			s.Set.Field = &actr.SetField{
				FieldName: set.Field.Name,
			}
		}
	}
	return &s, nil
}

func addRecallStatement(model *actr.Model, recall *recallStatement, production *actr.Production) (*actr.Statement, error) {
	errs := errorListWithContext{}

	name := recall.MemoryName
	memory := model.LookupMemory(name)
	if memory == nil {
		errs.Addc(&recall.Pos, "memory not found in production '%s': '%s'", production.Name, name)
		return nil, errs
	}

	pattern := createChunkPattern(recall.Pattern)

	s := actr.Statement{
		Recall: &actr.RecallStatement{
			Pattern:    pattern,
			MemoryName: name,
		},
	}

	return &s, nil
}

func addClearStatement(model *actr.Model, clear *clearStatement, production *actr.Production) (*actr.Statement, error) {
	errs := errorListWithContext{}
	bufferNames := clear.BufferNames

	for _, name := range bufferNames {
		buffer := model.LookupBuffer(name)
		if buffer == nil {
			errs.Addc(&clear.Pos, "buffer not found in production '%s': '%s'", production.Name, name)
			continue
		}
	}
	if errs.ErrorOrNil() != nil {
		return nil, errs
	}

	s := actr.Statement{
		Clear: &actr.ClearStatement{
			BufferNames: bufferNames,
		},
	}

	return &s, nil
}

func addPrintStatement(model *actr.Model, print *printStatement, production *actr.Production) (*actr.Statement, error) {
	s := actr.Statement{
		Print: &actr.PrintStatement{
			Args: print.Args.Strings(),
		},
	}

	return &s, nil
}

func addWriteStatement(model *actr.Model, write *writeStatement, production *actr.Production) (*actr.Statement, error) {
	errs := errorListWithContext{}

	name := write.TextOutputName
	textOutput := model.LookupTextOutput(name)
	if textOutput == nil {
		errs.Addc(&write.Pos, "text output not found in production '%s': '%s'", production.Name, name)
		return nil, errs
	}

	s := actr.Statement{
		Write: &actr.WriteStatement{
			Args:           write.Args.Strings(),
			TextOutputName: name},
	}

	return &s, nil
}
