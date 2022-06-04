package vanilla_actr

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/framework"

	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/issues"
	"github.com/asmaloney/gactar/util/numbers"
	"github.com/asmaloney/gactar/util/version"
)

var Info framework.Info = framework.Info{
	Name:           "vanilla",
	Language:       "commonlisp",
	FileExtension:  "lisp",
	ExecutableName: "sbcl",
}

type VanillaACTR struct {
	framework.Framework
	framework.WriterHelper
	model     *actr.Model
	modelName string
	tmpPath   string
	envPath   string
}

// New simply creates a new VanillaACTR instance and sets some paths from the context.
func New(ctx *cli.Context) (v *VanillaACTR, err error) {

	v = &VanillaACTR{
		tmpPath: ctx.Path("temp"),
		envPath: os.Getenv("VIRTUAL_ENV"),
	}

	return
}

func (VanillaACTR) Info() *framework.Info {
	return &Info
}

func (v *VanillaACTR) Initialize() (err error) {
	return framework.Setup(&Info)
}

func (VanillaACTR) ValidateModel(model *actr.Model) (log *issues.Log) {
	log = issues.New()
	return
}

func (v *VanillaACTR) SetModel(model *actr.Model) (err error) {
	if model.Name == "" {
		err = fmt.Errorf("model is missing name")
		return
	}

	v.model = model
	v.modelName = fmt.Sprintf("vanilla_%s", v.model.Name)

	return
}

func (c VanillaACTR) Model() (model *actr.Model) {
	return c.model
}

func (v *VanillaACTR) Run(initialBuffers framework.InitialBuffers) (result *framework.RunResult, err error) {
	modelFile, err := v.WriteModel(v.tmpPath, initialBuffers)
	if err != nil {
		return
	}

	// Save the current code for our result
	result = &framework.RunResult{
		FileName:      modelFile,
		GeneratedCode: v.GetContents(),
	}

	runFile, err := v.createRunFile(modelFile)
	if err != nil {
		return
	}

	// run it!
	cmd := exec.Command(runFile)
	output, err := cmd.CombinedOutput()
	output = removePreamble(output)
	if err != nil {
		err = fmt.Errorf("%s", string(output))
		return
	}

	result.Output = output

	return
}

func (v *VanillaACTR) WriteModel(path string, initialBuffers framework.InitialBuffers) (outputFileName string, err error) {
	patterns, err := framework.ParseInitialBuffers(v.model, initialBuffers)
	if err != nil {
		return
	}

	goal := patterns["goal"]

	outputFileName = fmt.Sprintf("%s.lisp", v.modelName)
	if path != "" {
		outputFileName = fmt.Sprintf("%s/%s", path, outputFileName)
	}

	err = filesystem.RemoveFile(outputFileName)
	if err != nil {
		return "", err
	}

	err = v.InitWriterHelper(outputFileName)
	if err != nil {
		return
	}
	defer func() {
		writerErr := v.CloseWriterHelper()
		if err == nil {
			err = writerErr
		} else if writerErr != nil {
			err = fmt.Errorf("%s; %w", err.Error(), writerErr)
		}
	}()

	v.writeHeader()

	v.Write("(clear-all)\n\n")

	v.Writeln("(define-model %s\n", v.modelName)

	v.Writeln("(sgp")

	// enable subsymbolic computations
	v.Writeln("\t:esc t")

	memory := v.model.Memory
	if memory.LatencyFactor != nil {
		v.Writeln("\t:lf %s", numbers.Float64Str(*memory.LatencyFactor))
	}

	if memory.LatencyExponent != nil {
		v.Writeln("\t:le %s", numbers.Float64Str(*memory.LatencyExponent))
	}

	if memory.RetrievalThreshold != nil {
		v.Writeln("\t:rt %s", numbers.Float64Str(*memory.RetrievalThreshold))
	}

	if memory.FinstSize != nil {
		v.Writeln("\t:declarative-num-finsts %d", *memory.FinstSize)
	}

	if memory.FinstTime != nil {
		v.Writeln("\t:declarative-finst-span %s", numbers.Float64Str(*memory.FinstTime))
	}

	if memory.MaxSpreadStrength != nil {
		v.Writeln("\t:mas %s", numbers.Float64Str(*memory.MaxSpreadStrength))

		goalActivation := v.model.Goal.SpreadingActivation
		if goalActivation != nil {
			v.Writeln("\t:ga %s", numbers.Float64Str(*goalActivation))
		}
	}

	procedural := v.model.Procedural
	if procedural.DefaultActionTime != nil {
		v.Writeln("\t:dat %s", numbers.Float64Str(*procedural.DefaultActionTime))
	}

	switch v.model.LogLevel {
	case "min":
		v.Writeln("\t:trace-detail low")
	case "info":
		v.Writeln("\t:trace-detail medium")
	case "detail":
		v.Writeln("\t:trace-detail high")
	}

	if v.model.TraceActivations {
		v.Writeln("\t:act t")
	}

	imaginal := v.model.ImaginalModule()
	if imaginal != nil {
		v.Writeln("\t:do-not-harvest imaginal")
		if imaginal.Delay != nil {
			v.Writeln("\t:imaginal-delay %s", numbers.Float64Str(*imaginal.Delay))
		}
	}
	v.Writeln(")\n")

	// chunks
	for _, chunk := range v.model.Chunks {
		if chunk.IsInternal() {
			continue
		}

		v.Writeln(";; amod line %d", chunk.AMODLineNumber)
		v.Writeln("(chunk-type %s %s)", chunk.Name, strings.Join(chunk.SlotNames, " "))
	}
	v.Writeln("")

	v.Writeln("(add-dm")
	for i, init := range v.model.Initializers {
		module := init.Module

		// allow the user-set goal to override the initializer
		if module.ModuleName() == "goal" && (goal != nil) {
			continue
		}

		initializer := module.ModuleName()

		if initializer == "imaginal" {
			continue
		}

		if initializer == "memory" {
			v.Writeln(" ;; amod line %d", init.AMODLineNumber)
			v.Writeln(" (fact_%d", i)
		} else {
			v.Writeln(" ;; amod line %d", init.AMODLineNumber)
			v.Writeln(" (%s", initializer)
		}

		v.outputPattern(init.Pattern, 1)
		v.Writeln(" )")
	}

	if goal != nil {
		v.Writeln(" (goal")
		v.outputPattern(goal, 1)
		v.Writeln(" )")
	}

	v.Writeln(")\n")

	// productions
	for _, production := range v.model.Productions {
		v.Writeln(";; amod line %d", production.AMODLineNumber)

		v.Writeln("(P %s", production.Name)
		if production.Description != nil {
			v.Writeln("\t\"%s\"", *production.Description)
		}

		for _, match := range production.Matches {
			v.outputMatch(match)
		}

		v.Writeln("\t==>")

		if production.DoStatements != nil {
			for _, statement := range production.DoStatements {
				v.outputStatement(statement)
			}
		}

		v.Writeln(")\n")
	}

	if imaginal != nil {
		v.Writeln(";; initialize our imaginal buffer")
		v.Writeln("(define-chunks (imaginal-init")

		// find our imaginal initializer and output it
		for _, init := range v.model.Initializers {
			if init.Module != nil {
				if init.Module.ModuleName() != "imaginal" {
					continue
				}

				v.outputPattern(init.Pattern, 1)
			}
		}
		v.Writeln("))")

		v.Writeln(`(set-buffer-chunk 'imaginal 'imaginal-init )`)
		v.Writeln("")
	}

	// Useful for debugging - output the contents of the imaginal buffer and the dm
	// v.Writeln("(buffer-chunk imaginal)")
	// v.Writeln("(dm)")

	v.Writeln("(goal-focus goal)")

	v.Writeln(")")

	return
}

func (v VanillaACTR) writeHeader() {
	v.Writeln(";;; Generated by gactar %s", version.BuildVersion)
	v.Writeln(";;;           on %s", time.Now().Format("2006-01-02 @ 15:04:05"))
	v.Writeln(";;;   https://github.com/asmaloney/gactar")
	v.Writeln("")
	v.Writeln(";;; *** NOTE: This is a generated file. Any changes may be overwritten.")
	v.Writeln("")

	if v.model.Description != "" {
		v.Write(";;; %s\n\n", v.model.Description)
	}

	v.outputAuthors()
}

func (v VanillaACTR) outputAuthors() {
	if len(v.model.Authors) == 0 {
		return
	}

	v.Writeln(";;; Authors:")

	for _, author := range v.model.Authors {
		v.Write(";;;\t\t%s\n", author)
	}

	v.Writeln("")
}

func (v VanillaACTR) outputPattern(pattern *actr.Pattern, tabs int) {
	tabbedItems := framework.KeyValueList{}
	tabbedItems.Add("isa", pattern.Chunk.Name)

	for i, slot := range pattern.Slots {
		slotName := pattern.Chunk.SlotNames[i]
		addPatternSlot(&tabbedItems, slotName, slot)
	}

	v.TabWrite(tabs, tabbedItems)
}

func (v VanillaACTR) outputMatch(match *actr.Match) {
	bufferName := match.Buffer.BufferName()
	chunkName := match.Pattern.Chunk.Name

	if actr.IsInternalChunkName(chunkName) {
		if chunkName == "_status" {
			status := match.Pattern.Slots[0]
			v.Writeln("\t?%s>", bufferName)

			if status.String() == "full" || status.String() == "empty" {
				v.Writeln("\t\tbuffer %s", status)
			} else {
				v.Writeln("\t\tstate %s", status)
			}
		}
	} else {
		v.Writeln("\t=%s>", bufferName)
		v.outputPattern(match.Pattern, 2)
	}
}

func addPatternSlot(tabbedItems *framework.KeyValueList, slotName string, patternSlot *actr.PatternSlot) {
	for _, item := range patternSlot.Items {
		if item.Wildcard {
			return
		}

		value := ""
		slot := ""

		if item.Negated {
			slot = "- "
		}

		switch {
		case item.Nil:
			value = "empty"

		case item.ID != nil:
			value = fmt.Sprintf(`"%s"`, *item.ID)

		case item.Num != nil:
			value = *item.Num

		case item.Var != nil:
			varName := strings.TrimPrefix(*item.Var, "?")
			value = fmt.Sprintf("=%s", varName)
		}

		slot += slotName

		tabbedItems.Add(slot, value)
	}
}

func (v VanillaACTR) outputStatement(s *actr.Statement) {
	switch {
	case s.Set != nil:
		buffer := s.Set.Buffer

		v.Writeln("\t=%s>", buffer.BufferName())

		if s.Set.Slots != nil {
			tabbedItems := framework.KeyValueList{}
			tabbedItems.Add("isa", s.Set.Chunk.Name)

			for _, slot := range *s.Set.Slots {
				slotName := slot.Name

				switch {
				case slot.Value.Nil:
					tabbedItems.Add(slotName, "empty")

				case slot.Value.Var != nil:
					tabbedItems.Add(slotName, fmt.Sprintf("=%s", *slot.Value.Var))

				case slot.Value.Number != nil:
					tabbedItems.Add(slotName, *slot.Value.Number)

				case slot.Value.Str != nil:
					tabbedItems.Add(slotName, fmt.Sprintf(`"%s"`, *slot.Value.Str))
				}
			}
			v.TabWrite(2, tabbedItems)
		} else if s.Set.Pattern != nil {
			v.outputPattern(s.Set.Pattern, 2)
		}

	case s.Recall != nil:
		v.Writeln("\t+retrieval>")
		v.outputPattern(s.Recall.Pattern, 2)

	case s.Print != nil:
		outputArgs := createOutputArgs(s.Print.Values)
		v.Write("\t!output!\t(%s)\n", outputArgs)

	case s.Clear != nil:
		for _, name := range s.Clear.BufferNames {
			v.Writeln("\t-%s>", name)
		}
	}
}

// createOutputArgs creates a string suitable for use in an !output! statement
// !output! is explained in:
//   ACT-R 7.21+ Reference Manual pg. 235
func createOutputArgs(values *[]*actr.Value) string {
	formatStr := `"`
	args := []string{}

	for _, v := range *values {
		switch {
		case v.Var != nil:
			formatStr += "~a"
			varName := strings.TrimPrefix(*v.Var, "?")
			args = append(args, fmt.Sprintf("=%s", varName))

		case v.Str != nil:
			formatStr += *v.Str

		case v.Number != nil:
			formatStr += *v.Number
		}
		// v.ID should not be possible because of validation
	}

	formatStr += `"`

	var argStr string
	if len(args) > 0 {
		argStr += " "

		for _, arg := range args {
			formatStr += " "
			formatStr += arg
		}
	}

	return formatStr + argStr
}

// createRunFile creates a lisp program to load ACTR and our model and then run them.
func (v VanillaACTR) createRunFile(modelFile string) (outputFile string, err error) {
	outputFile = fmt.Sprintf("%s_run.lisp", v.modelName)
	if v.tmpPath != "" {
		outputFile = fmt.Sprintf("%s/%s", v.tmpPath, outputFile)
	}

	err = v.InitWriterHelper(outputFile)
	if err != nil {
		return
	}
	defer func() {
		writerErr := v.CloseWriterHelper()
		if err == nil {
			err = writerErr
		} else if writerErr != nil {
			err = fmt.Errorf("%s; %w", err.Error(), writerErr)
		}
	}()

	v.Writeln("#!%s/bin/sbcl --script", v.envPath)
	v.Writeln(`(load "%s/actr/load-single-threaded-act-r.lisp")`, v.envPath)
	v.Writeln(`(load "%s")`, modelFile)
	v.Writeln(`(run 10.0)`)

	return
}

// removePreamble will remove the long preamble whenever ACT-R is loaded.
func removePreamble(text []byte) []byte {
	str := string(text)

	r := regexp.MustCompile(`(?s).+######### This is a single threaded build #########(.+)`)
	matches := r.FindAllStringSubmatch(str, -1)
	if len(matches) == 1 {
		str = strings.TrimSpace(matches[0][1])
	}

	return []byte(str)
}
