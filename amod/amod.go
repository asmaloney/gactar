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

	addConfig(model, amod.Config)
	initialize(model, amod.Init)
	addProductions(model, amod.Productions)

	return
}

func addConfig(model *actr.Model, config *configSection) {
	for _, name := range config.Buffers.Identifiers {
		var buffer actr.Buffer

		buffer.Name = name

		model.Buffers = append(model.Buffers, &buffer)
	}

	for _, mem := range config.Memories {
		var memory actr.Memory

		memory.Name = mem.Name

		for _, field := range mem.Fields.Fields {
			if field.Key == "buffer" {
				bufferName := field.Value

				buffer := model.LookupBuffer(bufferName)
				if buffer == nil {
					fmt.Printf("buffer not found for memory '%s': %s\n", mem.Name, bufferName)
					continue
				} else {
					memory.Buffer = buffer
				}
			}
		}

		model.Memories = append(model.Memories, &memory)
	}
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
		var prod actr.Production

		prod.Name = production.Name

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
