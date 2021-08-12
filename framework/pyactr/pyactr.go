package pyactr

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"gitlab.com/asmaloney/gactar/actr"
	"gitlab.com/asmaloney/gactar/framework"
)

type PyACTR struct {
	model     *actr.Model
	className string
	tmpPath   string
}

// New simply creates a new PyACTR instance and sets the tmp path.
func New(cli *cli.Context) (p *PyACTR, err error) {

	p = &PyACTR{tmpPath: "tmp"}

	return
}

func (p *PyACTR) Initialize() (err error) {
	_, err = framework.CheckForExecutable("python3")
	if err != nil {
		return
	}

	framework.PythonIdentify()

	err = framework.PythonCheckForPackage("pyactr")
	if err != nil {
		return
	}

	err = os.MkdirAll(p.tmpPath, os.ModePerm)
	if err != nil {
		return
	}

	return
}

func (p *PyACTR) SetModel(model *actr.Model) (err error) {
	if model.Name == "" {
		err = fmt.Errorf("model is missing name")
		return
	}

	p.model = model
	p.className = fmt.Sprintf("gactar_pyactr_%s", strings.Title(p.model.Name))

	return
}

func (p *PyACTR) Run(initialGoal string) (output []byte, err error) {
	return
}

func (p *PyACTR) WriteModel(path string) (outputFileName string, err error) {
	outputFileName = fmt.Sprintf("%s.py", p.className)
	if path != "" {
		outputFileName = fmt.Sprintf("%s/%s", path, outputFileName)
	}

	f, err := os.Create(outputFileName)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString("# This file is generated by gactar\n")
	f.WriteString(fmt.Sprintf("# %s\n", p.model.Description))
	f.WriteString("\n")

	f.WriteString("import pyactr as actr\n\n")
	f.WriteString(fmt.Sprintf("%s = actr.ACTRModel()\n\n", p.className))

	for _, buffer := range p.model.Buffers {
		// Note that this is not quite what pyactr is expecting. It has the concept of a class and some slots
		// where the class seems to be what we have in the first slot.
		f.WriteString(fmt.Sprintf("actr.chunktype(\"%s\", \"%s\")\n", buffer.Name, strings.Join(buffer.SlotNames, ", ")))
	}
	f.WriteString("\n")

	f.WriteString(fmt.Sprintf("dm = %s.decmem\n\n", p.className))

	// initialize
	if len(p.model.Initializers) > 0 {
		for _, init := range p.model.Initializers {
			f.WriteString("dm.add(actr.chunkstring(string=\"\"\"\n")

			buffer := init.Memory.Buffer

			f.WriteString(fmt.Sprintf("\tisa\t%s\n", buffer.Name))

			slots := strings.Split(init.Text, " ")

			for i, slot := range buffer.SlotNames {
				f.WriteString(fmt.Sprintf("\t%s\t%s\n", slot, slots[i]))
			}

			f.WriteString("\"\"\"))\n")
		}
	}

	f.WriteString("\n")

	// productions
	for _, production := range p.model.Productions {

		f.WriteString(fmt.Sprintf("%s.productionstring(name=\"%s\", string=\"\"\"\n", p.className, production.Name))
		for _, match := range production.Matches {
			outputMatch(f, match)
		}

		f.WriteString("\t==>\n")

		if production.DoStatements != nil {
			for _, statement := range production.DoStatements {
				outputStatement(f, statement)
			}
		}

		f.WriteString("\"\"\")\n")
	}

	f.WriteString("\n")

	// goal (this is just for testing - it will be created dynamically)
	f.WriteString(fmt.Sprintf("%s.goal.add(actr.chunkstring(string=\"\"\"\n", p.className))

	f.WriteString("\tisa\tgoal\n")
	f.WriteString("\tslot_1\tcountFrom\n")
	f.WriteString("\tslot_2\t2\n")
	f.WriteString("\tslot_3\t4\n")
	f.WriteString("\tslot_4\tstarting\n")

	f.WriteString("\"\"\"))\n")
	f.WriteString("\n")

	// run
	f.WriteString("if __name__ == \"__main__\":\n")
	f.WriteString(fmt.Sprintf("\tsim = %s.simulation()\n", p.className))
	f.WriteString("\tsim.run()\n")

	return
}

func outputMatch(f *os.File, match *actr.Match) {
	text := "g"
	if (match.Memory != nil) || (match.Buffer.Name == "retrieve") {
		text = "retrieval"
	}

	f.WriteString(fmt.Sprintf("\t=%s>\n", text))
	f.WriteString(fmt.Sprintf("\tisa\t%s\n", match.Buffer.Name))

	for i, slot := range match.Buffer.SlotNames {
		patternSlot := match.Pattern.Slots[i]

		outputPatternSlot(f, slot, patternSlot)
	}
}

func outputPatternSlot(f *os.File, slotName string, patternSlot *actr.PatternSlot) {
	value := ""

	for _, item := range patternSlot.Items {
		if item.ID != nil {
			value = *item.ID
		} else if item.Var != nil {
			if *item.Var == "?" {
				return
			}

			if item.Negated {
				value += "~"
			}
			value += "="
			value += strings.TrimPrefix(*item.Var, "?")
		}
	}

	f.WriteString(fmt.Sprintf("\t%s\t%s\n", slotName, value))
}

func outputStatement(f *os.File, s *actr.Statement) {
	if s.Set != nil {
		buffer := s.Set.Buffer

		text := "g"
		if buffer.Name == "retrieve" {
			text = "retrieval"
		}

		f.WriteString(fmt.Sprintf("\t=%s>\n", text))
		f.WriteString(fmt.Sprintf("\tisa\t%s\n", buffer.Name))

		if s.Set.Slot != nil {
			slotName := ""

			if s.Set.Slot.ArgNum != nil {
				slotName = buffer.SlotNames[*s.Set.Slot.ArgNum]
			} else if s.Set.Slot.Name != nil {
				slotName = *s.Set.Slot.Name
			}
			if s.Set.ID != nil {
				f.WriteString(fmt.Sprintf("\t%s\t=%s\n", slotName, *s.Set.ID))
			} else if s.Set.Number != nil {
				f.WriteString(fmt.Sprintf("\t%s\t%s\n", slotName, *s.Set.Number))
			}
		} else if s.Set.Pattern != nil {
			for i, slot := range s.Set.Buffer.SlotNames {
				patternSlot := s.Set.Pattern.Slots[i]

				outputPatternSlot(f, slot, patternSlot)
			}
		} else {
			f.WriteString("# writing text not yet handled\n")
		}
	} else if s.Recall != nil {
		memoryName := s.Recall.Memory.Buffer.Name
		f.WriteString("\t+retrieval>\n")
		f.WriteString(fmt.Sprintf("\tisa\t%s\n", memoryName))

		for i, slot := range s.Recall.Memory.Buffer.SlotNames {
			patternSlot := s.Recall.Pattern.Slots[i]

			outputPatternSlot(f, slot, patternSlot)
		}
	} else if s.Clear != nil {
		// for _, name := range s.Clear.BufferNames {
		f.WriteString(fmt.Sprintf("\t~%s>\n", "g"))
		// }
	}
}
