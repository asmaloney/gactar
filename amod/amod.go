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

func (elwc *errorListWithContext) Addc(pos *lexer.Position, e string, a ...interface{}) {
	str := fmt.Sprintf(e, a...)
	elwc.Addf("%s (line %d)", str, pos.Line)
}

func SetDebug(debug bool) {
	debugging = debug
}

func OutputEBNF() {
	fmt.Println(parser.String())
}

func GenerateModel(buffer string) (model *actr.Model, err error) {
	r := strings.NewReader(buffer)

	amod, err := parse(r)
	if err != nil {
		return
	}

	return generateModel(amod)
}

func GenerateModelFromFile(fileName string) (model *actr.Model, err error) {
	amod, err := parseFile(fileName)
	if err != nil {
		return
	}

	return generateModel(amod)
}

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

			prod.Matches = append(prod.Matches, &actr.Match{
				Name: name,
				Text: item.Text,
			})
		}

		prod.Do = append(prod.Do, production.Do.Texts...)

		model.Productions = append(model.Productions, &prod)
	}

	return errs.ErrorOrNil()
}
