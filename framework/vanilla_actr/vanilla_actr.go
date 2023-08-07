// Package vanilla_actr provides functions to output the internal actr data structures in Lisp
// suitable for running using the ACT-R code, and to run those models on the Clozure Common Lisp compiler.
package vanilla_actr

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/framework"

	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/executil"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/issues"
	"github.com/asmaloney/gactar/util/lisp"
	"github.com/asmaloney/gactar/util/numbers"
)

func init() {
	// We only support 64-bit. Nobody still uses 32-bit, right?
	osArch := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOARCH == "386" {
		fmt.Println("ERROR: I don't know how to set the Clozure Common Lisp compiler for", osArch)
		return
	}

	cclExecutableName, err := lisp.GetExecutableName()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Store the name of the cclExecutable - it's different depending on platform & architecture.
	Info.ExecutableName = cclExecutableName
}

var Info framework.Info = framework.Info{
	Name:          "vanilla",
	Language:      "commonlisp",
	FileExtension: "lisp",
	// ExecutableName: have to set this in init()
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
func New(settings *cli.Settings) (v *VanillaACTR, err error) {
	v = &VanillaACTR{
		tmpPath: settings.TempPath,
		envPath: os.Getenv("VIRTUAL_ENV"),
	}

	err = framework.Setup(&Info)
	if err != nil {
		v = nil
		return
	}

	return
}

func (VanillaACTR) Info() *framework.Info {
	return &Info
}

func (VanillaACTR) ValidateModel(model *actr.Model) (log *issues.Log) {
	log = issues.New()
	return
}

func (v *VanillaACTR) SetModel(model *actr.Model) (err error) {
	if model.Name == "" {
		err = framework.ErrModelMissingName
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

	if Info.ExecutableName == "" {
		err = &framework.ErrExecutableNotSet{Name: "Clozure Common Lisp"}
		return
	}

	// run it!
	output, err := executil.ExecCommand(Info.ExecutableName, "--batch", "--quiet", "--load", runFile)
	output = removePreamble(output)
	if err != nil {
		err = &executil.ErrExecuteCommand{Output: output}
		return
	}

	result.Output = []byte(output)

	return
}

// WriteModel converts the internal actr.Model to Lisp and writes it to a file.
func (v *VanillaACTR) WriteModel(path string, initialBuffers framework.InitialBuffers) (outputFileName string, err error) {
	outputFileName = fmt.Sprintf("%s.lisp", v.modelName)
	if path != "" {
		outputFileName = fmt.Sprintf("%s/%s", path, outputFileName)
	}

	err = filesystem.RemoveFile(outputFileName)
	if err != nil {
		return "", err
	}

	_, err = v.GenerateCode(initialBuffers)
	if err != nil {
		return
	}

	err = v.WriteFile(outputFileName)
	if err != nil {
		return
	}

	return
}

// GenerateCode converts the internal actr.Model to Lisp code.
func (v *VanillaACTR) GenerateCode(initialBuffers framework.InitialBuffers) (code []byte, err error) {
	patterns, err := framework.ParseInitialBuffers(v.model, initialBuffers)
	if err != nil {
		return
	}

	goal := patterns["goal"]

	err = v.InitWriterHelper()
	if err != nil {
		return
	}

	v.writeHeader()

	v.Writeln("(clear-all)")

	v.Writeln("")

	// add any extra buffers
	extraBuffers := v.model.LookupModule("extra_buffers")
	if extraBuffers != nil {
		v.Writeln(`(require-compiled "GOAL-STYLE-MODULE")`)
		v.Writeln("")
		v.Writeln(";; define a goal-style module for each extra buffer")
		for _, buff := range extraBuffers.Buffers().Names() {
			v.Writeln(
				`(define-module %[1]s (%[1]s) nil
	:version "1.0"
	:documentation "Extra buffer: %[1]s"
	:query goal-style-query
	:request goal-style-request
	:buffer-mod goal-style-mod-request)
	`, buff)
		}
	}

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

	if memory.Decay != nil {
		v.Writeln("\t:bll %s", numbers.Float64Str(*memory.Decay))
	}

	if memory.MaxSpreadStrength != nil {
		v.Writeln("\t:mas %s", numbers.Float64Str(*memory.MaxSpreadStrength))

		goalActivation := v.model.Goal.Buffer().SpreadingActivation()
		if goalActivation != 0.0 {
			v.Writeln("\t:ga %s", numbers.Float64Str(goalActivation))
		}
	}

	if memory.InstantaneousNoise != nil {
		v.Writeln("\t:ans %s", numbers.Float64Str(*memory.InstantaneousNoise))
	}

	if memory.MismatchPenalty != nil {
		v.Writeln("\t:mp %s", numbers.Float64Str(*memory.MismatchPenalty))
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

	// random
	if v.model.RandomSeed != nil {
		v.Writeln("(sgp :seed (%d 0))\n", *v.model.RandomSeed)
	}

	// chunks
	for _, chunk := range v.model.Chunks {
		if chunk.IsInternal() {
			continue
		}

		v.Writeln(";; amod line %d", chunk.AMODLineNumber)
		v.Writeln("(chunk-type %s %s)", chunk.TypeName, strings.Join(chunk.SlotNames, " "))
	}
	v.Writeln("")

	v.writeInitializers(goal)

	v.writeSimilarities()

	v.writeProductions()

	// Useful for debugging - output the contents of the imaginal buffer and the dm
	// v.Writeln("(buffer-chunk imaginal)")
	// v.Writeln("(dm)")

	v.Writeln("(goal-focus goal)")

	v.Writeln(")")

	code = v.GetContents()
	return
}

func (v VanillaACTR) writeHeader() {
	v.Writeln(";;; Generated by gactar %s", framework.GactarVersion)
	v.Writeln(";;;           on %s", framework.TimeNow().Format("2006-01-02 @ 15:04:05"))
	v.Writeln(";;;   https://github.com/asmaloney/gactar")
	v.Writeln("")
	v.Writeln(";;; *** NOTE: This is a generated file. Any changes may be overwritten.")
	v.Writeln("")

	if v.model.Description != "" {
		v.Write(";;; %s\n\n", v.model.Description)
	}

	v.writeAuthors()
}

func (v VanillaACTR) writeAuthors() {
	if len(v.model.Authors) == 0 {
		return
	}

	v.Writeln(";;; Authors:")

	for _, author := range v.model.Authors {
		v.Write(";;;\t\t%s\n", author)
	}

	v.Writeln("")
}

func (v VanillaACTR) writeImplicitChunks() {
	if !v.model.HasImplicitChunks() {
		return
	}

	v.Writeln(" ;; declare implicit chunks without slots to avoid warnings")
	v.SetLineLen(80)
	for _, chunkName := range v.model.ImplicitChunks {
		v.Write(" (%s)", chunkName)
	}
	v.ResetLineLen()

	v.Writeln("\n")
}

func (v VanillaACTR) writeBufferInitializer(bufferName string, lineNumber int, pattern *actr.Pattern) {
	v.Writeln(";; initialize our %q buffer", bufferName)
	if lineNumber != 0 {
		v.Writeln(";; amod line %d", lineNumber)
	}
	v.Writeln("(set-buffer-chunk '%s '(", bufferName)
	v.outputPattern(pattern, 1)
	v.Writeln("))")
	v.Writeln("")
}

func (v VanillaACTR) writeInitializers(goal *actr.Pattern) {
	// First write out our declarative memory
	v.Writeln(";; initialize our declarative memory")
	v.Writeln("(add-dm")

	v.writeImplicitChunks()

	factNum := 0
	for _, init := range v.model.Initializers {
		moduleName := init.Module.ModuleName()

		if moduleName == "memory" {
			v.Writeln(" ;; amod line %d", init.AMODLineNumber)
			if init.ChunkName != nil {
				v.Writeln(" (%s", *init.ChunkName)
			} else {
				v.Writeln(" (%s_%d", init.Pattern.Chunk.TypeName, factNum)
				factNum++
			}

			v.outputPattern(init.Pattern, 1)
			v.Writeln(" )")
		} else if moduleName == "goal" {
			// allow the user-set goal to override the initializer
			if goal != nil {
				v.Writeln(" ;; goal set by user")
				v.Writeln(" (goal")
				v.outputPattern(goal, 1)
				v.Writeln(" )")
			} else {
				v.Writeln(" ;; amod line %d", init.AMODLineNumber)
				v.Writeln(" (goal")
				v.outputPattern(init.Pattern, 1)
				v.Writeln(" )")
			}
		}
	}

	v.Writeln(")\n")

	// now everything else
	for _, init := range v.model.Initializers {
		module := init.Module
		moduleName := module.ModuleName()

		switch {
		case moduleName == "memory":
			continue

		case moduleName == "goal":
			continue

		// for extra buffers, we use the buffer name
		case moduleName == "extra_buffers":
			v.writeBufferInitializer(init.Buffer.Name(), init.AMODLineNumber, init.Pattern)

		default:
			v.writeBufferInitializer(moduleName, init.AMODLineNumber, init.Pattern)
		}
	}
}

func (v VanillaACTR) writeSimilarities() {
	if len(v.model.Similarities) == 0 {
		return
	}

	v.Writeln("(set-similarities")

	for _, similar := range v.model.Similarities {
		v.Writeln("    ;; amod line %d", similar.AMODLineNumber)
		v.Writeln("    (%s %s %s)", similar.ChunkOne, similar.ChunkTwo, numbers.Float64Str(similar.Value))
	}

	v.Writeln(")\n")
}

func (v VanillaACTR) writeProductions() {
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
}

func (v VanillaACTR) outputPattern(pattern *actr.Pattern, tabs int) {
	tabbedItems := framework.KeyValueList{}
	tabbedItems.Add("isa", pattern.Chunk.TypeName)

	for i, slot := range pattern.Slots {
		slotName := pattern.Chunk.SlotNames[i]
		addPatternSlot(&tabbedItems, slotName, slot)
	}

	v.TabWrite(tabs, tabbedItems)
}

func (v VanillaACTR) outputMatch(match *actr.Match) {
	tabbedItems := framework.KeyValueList{}

	// check for case where we need to combine module & buffer checks
	if (match.BufferState != nil) && (match.ModuleState != nil) {
		bufferName := match.BufferState.Buffer.Name()

		v.Writeln("\t?%s>", bufferName)
		tabbedItems.Add("buffer", match.BufferState.State)
		tabbedItems.Add("state", match.ModuleState.State)
		v.TabWrite(2, tabbedItems)

		return
	}

	switch {
	case match.BufferPattern != nil:
		bufferName := match.BufferPattern.Buffer.Name()

		v.Writeln("\t=%s>", bufferName)
		v.outputPattern(match.BufferPattern.Pattern, 2)

	case match.BufferState != nil:
		bufferName := match.BufferState.Buffer.Name()

		v.Writeln("\t?%s>", bufferName)
		tabbedItems.Add("buffer", match.BufferState.State)
		v.TabWrite(2, tabbedItems)

	case match.ModuleState != nil:
		bufferName := match.ModuleState.Buffer.Name()

		v.Writeln("\t?%s>", bufferName)
		tabbedItems.Add("state", match.ModuleState.State)
		v.TabWrite(2, tabbedItems)
	}
}

func addPatternSlot(tabbedItems *framework.KeyValueList, slotName string, slot *actr.PatternSlot) {
	if slot.Wildcard {
		return
	}

	value := ""
	slotStr := ""

	if slot.Negated {
		slotStr = "- "
	}

	switch {
	case slot.Nil:
		value = "empty"

	case slot.ID != nil:
		value = *slot.ID

	case slot.Str != nil:
		value = fmt.Sprintf("%q", *slot.Str)

	case slot.Num != nil:
		value = *slot.Num

	case slot.Var != nil:
		varName := strings.TrimPrefix(*slot.Var.Name, "?")
		value = fmt.Sprintf("=%s", varName)
	}

	slotStr += slotName

	tabbedItems.Add(slotStr, value)

	// Check for constraints on a var and output them
	if slot.Var != nil {
		if len(slot.Var.Constraints) > 0 {
			for _, constraint := range slot.Var.Constraints {
				slotStr := ""

				if constraint.Comparison == actr.NotEqual {
					slotStr = "- "
				}

				slotStr += slotName

				if constraint.RHS.Var != nil {
					varName := strings.TrimPrefix(*constraint.RHS.Var, "?")
					value = fmt.Sprintf("=%s", varName)

					tabbedItems.Add(slotStr, value)
				} else {
					tabbedItems.Add(slotStr, constraint.RHS.String())
				}
			}
		}
	}
}

func (p VanillaACTR) outputRequestParameters(params map[string]string, tabs int) {
	tabbedItems := framework.KeyValueList{}

	for _, param := range maps.Keys(params) {
		if param != "recently_retrieved" {
			continue
		}

		name := param

		if param == "recently_retrieved" {
			name = ":recently-retrieved"
		}

		tabbedItems.Add(name, params[param])
	}

	p.TabWrite(tabs, tabbedItems)
}

func (v VanillaACTR) outputStatement(s *actr.Statement) {
	switch {
	case s.Set != nil:
		buffer := s.Set.Buffer

		v.Writeln("\t=%s>", buffer.Name())

		if s.Set.Slots != nil {
			tabbedItems := framework.KeyValueList{}
			tabbedItems.Add("isa", s.Set.Chunk.TypeName)

			for _, slot := range *s.Set.Slots {
				slotName := slot.Name

				switch {
				case slot.Value.Nil != nil:
					tabbedItems.Add(slotName, "empty")

				case slot.Value.Var != nil:
					tabbedItems.Add(slotName, fmt.Sprintf("=%s", *slot.Value.Var))

				case slot.Value.ID != nil:
					tabbedItems.Add(slotName, *slot.Value.ID)

				case slot.Value.Number != nil:
					tabbedItems.Add(slotName, *slot.Value.Number)

				case slot.Value.Str != nil:
					tabbedItems.Add(slotName, fmt.Sprintf(`%q`, *slot.Value.Str))
				}
			}
			v.TabWrite(2, tabbedItems)
		} else if s.Set.Pattern != nil {
			v.outputPattern(s.Set.Pattern, 2)
		}

	case s.Recall != nil:
		v.Writeln("\t+retrieval>")
		v.outputPattern(s.Recall.Pattern, 2)
		v.outputRequestParameters(s.Recall.RequestParameters, 2)

	case s.Print != nil:
		outputArgs := createOutputArgs(s.Print.Values)
		v.Write("\t!output!\t(%s)\n", outputArgs)

	case s.Clear != nil:
		for _, name := range s.Clear.BufferNames {
			v.Writeln("\t-%s>", name)
		}

	case s.Stop != nil:
		v.Writeln("\t!stop!")
	}
}

// createOutputArgs creates a string suitable for use in an !output! statement
// !output! is explained in:
//
//	ACT-R 7.21+ Reference Manual pg. 235
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
	err = v.InitWriterHelper()
	if err != nil {
		return
	}

	v.Writeln(`(load "%s/actr/load-single-threaded-act-r.lisp")`, v.envPath)
	v.Writeln(`(load "%s")`, modelFile)

	// TODO: We should be able to set this somewhere.
	// 10.0 is an arbitrary length of time.
	v.Writeln(`(run 10.0)`)

	outputFile = fmt.Sprintf("%s_run.lisp", v.modelName)
	if v.tmpPath != "" {
		outputFile = fmt.Sprintf("%s/%s", v.tmpPath, outputFile)
	}

	err = v.WriteFile(outputFile)
	if err != nil {
		return
	}

	return
}

// removePreamble will remove the long preamble whenever ACT-R is loaded.
func removePreamble(text string) string {
	r := regexp.MustCompile(`(?s).+######### This is a single threaded build #########(.+)`)
	matches := r.FindAllStringSubmatch(text, -1)
	if len(matches) == 1 {
		text = strings.TrimSpace(matches[0][1])
	}

	return text
}
