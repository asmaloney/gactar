package amod

import (
	"fmt"
	"strings"

	"gitlab.com/asmaloney/gactar/actr"
)

var debugging bool = false

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
	initialize(model, amod.Init)
	addProductions(model, amod.Productions)

	return
}

func addConfig(model *actr.Model, config *configSection) (err error) {
	errs := []string{}

	if config.hasACTR() {
		for _, field := range config.ACTR.Fields {
			switch field.Key {
			case "log":
				if field.Value.Number != nil {
					model.Logging = (*field.Value.Number != 0)
				} else {
					model.Logging = (strings.ToLower(*field.Value.String) == "true")
				}
			default:
				errs = append(errs, fmt.Sprintf("Unrecognized field in actr section: '%s'", field.Key))
			}
		}
	}

	if config.hasBuffers() {
		for _, name := range config.Buffers.Identifiers {
			buffer := actr.Buffer{
				Name: name,
			}

			model.Buffers = append(model.Buffers, &buffer)
		}
	}

	if config.hasMemories() {
		for _, mem := range config.Memories.Memories {
			memory := actr.Memory{
				Name: mem.Name,
			}

			for _, field := range mem.Fields.Fields {
				switch field.Key {
				case "buffer":
					if field.Value.Number != nil {
						errs = append(errs, fmt.Sprintf("buffer should not be a number in memory '%s': %v\n", mem.Name, *field.Value.Number))
						continue
					}

					bufferName := field.Value.String

					buffer := model.LookupBuffer(*bufferName)
					if buffer == nil {
						errs = append(errs, fmt.Sprintf("buffer not found for memory '%s': %s\n", mem.Name, *bufferName))
						continue
					} else {
						memory.Buffer = buffer
					}

				case "latency":
					if field.Value.String != nil {
						errs = append(errs, fmt.Sprintf("latency should not be a string in memory '%s': %v\n", mem.Name, *field.Value.String))
						continue
					}

					memory.Latency = field.Value.Number

				case "threshold":
					if field.Value.String != nil {
						errs = append(errs, fmt.Sprintf("threshold should not be a string in memory '%s': %v\n", mem.Name, *field.Value.String))
						continue
					}

					memory.Threshold = field.Value.Number

				case "max_time":
					if field.Value.String != nil {
						errs = append(errs, fmt.Sprintf("max_time should not be a string in memory '%s': %v\n", mem.Name, *field.Value.String))
						continue
					}

					memory.MaxTime = field.Value.Number

				case "finst_size":
					if field.Value.String != nil {
						errs = append(errs, fmt.Sprintf("finst_size should not be a string in memory '%s': %v\n", mem.Name, *field.Value.String))
						continue
					}

					size := int(*field.Value.Number)
					memory.FinstSize = &size

				case "finst_time":
					if field.Value.String != nil {
						errs = append(errs, fmt.Sprintf("finst_time should not be a string in memory '%s': %v\n", mem.Name, *field.Value.String))
						continue
					}

					memory.FinstTime = field.Value.Number

				default:
					errs = append(errs, fmt.Sprintf("Unrecognized field in memory '%s': '%s'", memory.Name, field.Key))
				}
			}

			model.Memories = append(model.Memories, &memory)
		}
	}

	if config.hasTextOutputs() {
		for _, name := range config.TextOutputs.Identifiers {
			textOutput := actr.TextOutput{
				Name: name,
			}

			model.TextOutputs = append(model.TextOutputs, &textOutput)
		}
	}

	if len(errs) == 0 {
		return
	}

	return fmt.Errorf(strings.Join(errs, "\n"))
}

func initialize(model *actr.Model, init *initSection) {
	for _, initializer := range init.Initializers {
		name := initializer.Name

		memory := model.LookupMemory(name)
		if memory == nil {
			fmt.Printf("memory not found for initialization '%s'\n", name)
			continue
		}

		for _, item := range initializer.Items {
			init := actr.Initializer{
				MemoryName: memory.Name,
				Text:       item.Item,
			}

			model.Initializers = append(model.Initializers, &init)
		}
	}
}

func addProductions(model *actr.Model, productions *productionSection) {
	for _, production := range productions.Productions {
		prod := actr.Production{
			Name: production.Name,
		}

		for _, item := range production.Match.Items {
			name := item.Name

			exists := model.BufferOrMemoryExists(name)
			if !exists {
				fmt.Printf("buffer or memory not found for production '%s': %s\n", prod.Name, name)
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
}
