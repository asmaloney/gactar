package pyactr

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

	framework.PythonIdentify("pyactr")

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
	outputFile, err := p.WriteModel(p.tmpPath, initialGoal)
	if err != nil {
		return
	}

	// run it!
	cmd := exec.Command("python3", outputFile)

	output, err = cmd.CombinedOutput()
	output = removeWarning(output)
	if err != nil {
		err = fmt.Errorf("%s", string(output))
		return
	}

	return
}

func (p *PyACTR) WriteModel(path, initialGoal string) (outputFileName string, err error) {
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

	for _, chunk := range p.model.Chunks {
		f.WriteString(fmt.Sprintf("actr.chunktype(\"%s\", \"%s\")\n", chunk.Name, strings.Join(chunk.SlotNames, ", ")))
	}
	f.WriteString("\n")

	f.WriteString(fmt.Sprintf("dm = %s.decmem\n\n", p.className))

	// initialize
	if len(p.model.Initializers) > 0 {
		for _, init := range p.model.Initializers {
			f.WriteString("dm.add(actr.chunkstring(string=\"\"\"\n")

			chunkName, slots := actr.SplitStringForChunk(init.Text)
			chunk := p.model.LookupChunk(chunkName)

			if chunk == nil {
				err = fmt.Errorf("cannot find chunk named '%s' in initializer", chunkName)
				return
			}

			f.WriteString(fmt.Sprintf("\tisa\t%s\n", chunkName))

			for i, name := range chunk.SlotNames {
				value := slots[i]
				f.WriteString(fmt.Sprintf("\t%s\t%s\n", name, value))
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

	if initialGoal != "" {
		// add our goal...
		f.WriteString(fmt.Sprintf("%s.goal.add(actr.chunkstring(string=\"\"\"\n", p.className))

		chunkName, slots := actr.SplitStringForChunk(initialGoal)
		chunk := p.model.LookupChunk(chunkName)

		if chunk == nil {
			err = fmt.Errorf("cannot find chunk named '%s' from initial goal", chunkName)
			return
		}

		if len(slots) != chunk.NumSlots {
			err = fmt.Errorf("expecting %d slots for '%s', found %d", chunk.NumSlots, chunkName, len(slots))
			return
		}

		f.WriteString(fmt.Sprintf("\tisa\t%s\n", chunkName))

		for i, name := range chunk.SlotNames {
			value := slots[i]
			f.WriteString(fmt.Sprintf("\t%s\t%s\n", name, value))
		}

		f.WriteString("\"\"\"))\n")
		f.WriteString("\n")

		// ...and our code to run
		f.WriteString("if __name__ == \"__main__\":\n")
		f.WriteString(fmt.Sprintf("\tsim = %s.simulation()\n", p.className))
		f.WriteString("\tsim.run()\n")
	}

	return
}

func outputMatch(f *os.File, match *actr.Match) {
	text := "g"
	if (match.Memory != nil) || (match.Buffer.Name == "retrieve") {
		text = "retrieval"
	}

	f.WriteString(fmt.Sprintf("\t=%s>\n", text))
	f.WriteString(fmt.Sprintf("\tisa\t%s\n", match.Pattern.Chunk.Name))

	// TODO Not sure how to handle memory here.
	// e.g.  memory: `error:True`
	if match.Buffer != nil {
		for i, slot := range match.Pattern.Slots {
			slotName := match.Pattern.Chunk.SlotNames[i]
			outputPatternSlot(f, slotName, slot)
		}
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

		text := "g" // default to "goal"
		if buffer.Name == "retrieve" {
			text = "retrieval"
		}

		f.WriteString(fmt.Sprintf("\t=%s>\n", text))

		if s.Set.Slot != nil {
			f.WriteString(fmt.Sprintf("\tisa\t%s\n", s.Set.Chunk.Name))

			slotName := *s.Set.Slot

			if s.Set.ID != nil {
				f.WriteString(fmt.Sprintf("\t%s\t=%s\n", slotName, *s.Set.ID))
			} else if s.Set.Number != nil {
				f.WriteString(fmt.Sprintf("\t%s\t%s\n", slotName, *s.Set.Number))
			}
		} else if s.Set.Pattern != nil {
			f.WriteString(fmt.Sprintf("\tisa\t%s\n", s.Set.Pattern.Chunk.Name))

			for i, slot := range s.Set.Pattern.Slots {
				slotName := s.Set.Pattern.Chunk.SlotNames[i]
				outputPatternSlot(f, slotName, slot)
			}
		} else {
			f.WriteString("# writing text not yet handled\n")
		}
	} else if s.Recall != nil {
		chunk := s.Recall.Pattern.Chunk

		f.WriteString("\t+retrieval>\n")
		f.WriteString(fmt.Sprintf("\tisa\t%s\n", chunk.Name))

		for i, slot := range s.Recall.Pattern.Slots {
			slotName := chunk.SlotNames[i]
			outputPatternSlot(f, slotName, slot)
		}

	} else if s.Clear != nil {
		// for _, name := range s.Clear.BufferNames {
		f.WriteString(fmt.Sprintf("\t~%s>\n", "g"))
		// }
	}
}

// removeWarning will remove the long warning whenever pyactr is run without tkinter.
func removeWarning(text []byte) []byte {
	str := string(text)

	r := regexp.MustCompile(`(?s).+warnings.warn\("Simulation GUI is set to False."\)(.+)`)
	matches := r.FindAllStringSubmatch(str, -1)
	if len(matches) == 1 {
		str = strings.TrimSpace(matches[0][1])
	}

	return []byte(str)
}
