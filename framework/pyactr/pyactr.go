package pyactr

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"gitlab.com/asmaloney/gactar/actr"
	"gitlab.com/asmaloney/gactar/framework"
)

type PyACTR struct {
	framework.WriterHelper
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

	framework.IdentifyYourself("pyactr", "python3")

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

	err = p.InitWriterHelper(outputFileName)
	if err != nil {
		return
	}
	defer p.CloseWriterHelper()

	p.Writeln("# This file is generated by gactar %s", time.Now().Format("2006-01-02 15:04:05"))
	p.Writeln("# https://github.com/asmaloney/gactar")
	p.Writeln("")
	p.Writeln("# *** This is a generated file. Any changes may be overwritten.")
	p.Writeln("")
	p.Write("# %s\n\n", p.model.Description)

	p.Write("import pyactr as actr\n\n")
	p.Write("%s = actr.ACTRModel()\n\n", p.className)

	// chunks
	for _, chunk := range p.model.Chunks {
		if chunk.IsInternal() {
			continue
		}

		p.Writeln("actr.chunktype(\"%s\", \"%s\")", chunk.Name, strings.Join(chunk.SlotNames, ", "))
	}
	p.Writeln("")

	p.Write("dm = %s.decmem\n\n", p.className)

	// initialize
	for _, init := range p.model.Initializers {
		p.Writeln("dm.add(actr.chunkstring(string=\"\"\"")

		chunkName, slots := actr.SplitStringForChunk(init.Text)
		chunk := p.model.LookupChunk(chunkName)

		if chunk == nil {
			err = fmt.Errorf("cannot find chunk named '%s' in initializer", chunkName)
			return
		}

		tabbedItems := framework.KeyValueList{}
		tabbedItems.Add("isa", chunkName)

		err = tabbedItems.AddArrays(chunk.SlotNames, slots)
		if err != nil {
			return
		}

		p.TabWrite(1, tabbedItems)

		p.Writeln("\"\"\"))")
	}

	p.Writeln("")

	// productions
	for _, production := range p.model.Productions {
		p.Writeln("%s.productionstring(name=\"%s\", string=\"\"\"", p.className, production.Name)
		for _, match := range production.Matches {
			p.outputMatch(match)
		}

		p.Write("\t==>\n")

		if production.DoStatements != nil {
			for _, statement := range production.DoStatements {
				p.outputStatement(statement)
			}
		}

		p.Writeln("\"\"\")")
	}

	p.Writeln("")

	if initialGoal != "" {
		// add our goal...
		p.Write("%s.goal.add(actr.chunkstring(string=\"\"\"\n", p.className)

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

		tabbedItems := framework.KeyValueList{}
		tabbedItems.Add("isa", chunkName)

		err = tabbedItems.AddArrays(chunk.SlotNames, slots)
		if err != nil {
			return
		}

		p.TabWrite(1, tabbedItems)

		p.Writeln("\"\"\"))")
		p.Writeln("")

		// ...and our code to run
		p.Writeln("if __name__ == \"__main__\":")
		p.Writeln("\tsim = %s.simulation()", p.className)
		p.Writeln("\tsim.run()")
	}

	return
}

func (p *PyACTR) outputMatch(match *actr.Match) {
	if match.Buffer != nil {
		bufferName := match.Buffer.Name
		if bufferName == "goal" {
			bufferName = "g"
		}

		chunkName := match.Pattern.Chunk.Name

		if actr.IsInternalChunkName(chunkName) {
			if chunkName == "_status" {
				status := match.Pattern.Slots[0]
				p.Writeln("\t?%s>", bufferName)
				p.Writeln("\t\tbuffer %s", status)
			}
		} else {
			p.Writeln("\t=%s>", bufferName)

			tabbedItems := framework.KeyValueList{}
			tabbedItems.Add("isa", chunkName)

			for i, slot := range match.Pattern.Slots {
				slotName := match.Pattern.Chunk.SlotNames[i]
				addPatternSlot(&tabbedItems, slotName, slot)
			}

			p.TabWrite(2, tabbedItems)
		}
	} else if match.Memory != nil {
		bufferName := "retrieval"

		chunkName := match.Pattern.Chunk.Name
		if actr.IsInternalChunkName(chunkName) {
			if chunkName == "_status" {
				status := match.Pattern.Slots[0]
				p.Writeln("\t?%s>", bufferName)
				p.Writeln("\t\tstate %s", status)
			}
		} else {
			p.Writeln("\t=%s>", bufferName)
			p.Writeln("\t\tisa\t%s", chunkName)
		}
	}
}

func addPatternSlot(tabbedItems *framework.KeyValueList, slotName string, patternSlot *actr.PatternSlot) {
	for _, item := range patternSlot.Items {
		value := ""
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

		tabbedItems.Add(slotName, value)
	}
}

func (p *PyACTR) outputStatement(s *actr.Statement) {
	if s.Set != nil {
		buffer := s.Set.Buffer
		bufferName := buffer.Name
		if bufferName == "goal" {
			bufferName = "g"
		}

		p.Write("\t=%s>\n", bufferName)

		if s.Set.Slots != nil {
			tabbedItems := framework.KeyValueList{}
			tabbedItems.Add("isa", s.Set.Chunk.Name)

			for _, slot := range *s.Set.Slots {
				slotName := slot.Name

				if slot.Value.ID != nil {
					tabbedItems.Add(slotName, fmt.Sprintf("=%s", *slot.Value.ID))
				} else if slot.Value.Number != nil {
					tabbedItems.Add(slotName, *slot.Value.Number)
				} else if slot.Value.Str != nil {
					tabbedItems.Add(slotName, *slot.Value.Str)
				}
			}
			p.TabWrite(2, tabbedItems)
		} else if s.Set.Pattern != nil {
			tabbedItems := framework.KeyValueList{}
			tabbedItems.Add("isa", s.Set.Pattern.Chunk.Name)

			for i, slot := range s.Set.Pattern.Slots {
				slotName := s.Set.Pattern.Chunk.SlotNames[i]
				addPatternSlot(&tabbedItems, slotName, slot)
			}

			p.TabWrite(2, tabbedItems)
		} else {
			p.Writeln("# writing text not yet handled")
		}
	} else if s.Recall != nil {
		chunk := s.Recall.Pattern.Chunk

		p.Writeln("\t+retrieval>")

		tabbedItems := framework.KeyValueList{}
		tabbedItems.Add("isa", chunk.Name)

		for i, slot := range s.Recall.Pattern.Slots {
			slotName := chunk.SlotNames[i]
			addPatternSlot(&tabbedItems, slotName, slot)
		}

		p.TabWrite(2, tabbedItems)
	} else if s.Clear != nil {
		for _, name := range s.Clear.BufferNames {
			if name == "goal" {
				name = "g"
			}

			p.Writeln("\t~%s>", name)
		}
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