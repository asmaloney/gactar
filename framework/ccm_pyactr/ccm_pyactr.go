// Package ccm_pyactr provides functions to output the internal actr data structures in Python suitable
// for running using CCM's python_actr package, and to run those models using Python.
package ccm_pyactr

import (
	_ "embed"
	"fmt"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/framework"

	"github.com/asmaloney/gactar/util/executil"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/issues"
	"github.com/asmaloney/gactar/util/numbers"
	"github.com/asmaloney/gactar/util/runoptions"
)

//go:embed ccm_print.py
var ccmPrintPython string

//go:embed gactar_ccm_activate_trace.py
var gactarActivateTraceFile string

const (
	ccmPrintFileName              = "ccm_print.py"
	ccmPrintImportName            = "ccm_print"
	gactarActivateTraceFileName   = "gactar_ccm_activate_trace.py"
	gactarActivateTraceImportName = "gactar_ccm_activate_trace"
)

var Info framework.Info = framework.Info{
	Name:           "ccm",
	Language:       "python",
	FileExtension:  "py",
	ExecutableName: "python",

	PythonRequiredPackages: []string{"python_actr"},
}

type CCMPyACTR struct {
	framework.Framework
	framework.WriterHelper

	tmpPath string

	model     *actr.Model
	className string
}

// New creates a new CCMPyACTR instance and sets the temp path.
func New(tempPath string) (c *CCMPyACTR, err error) {
	c = &CCMPyACTR{tmpPath: tempPath}

	err = framework.Setup(&Info)
	if err != nil {
		c = nil
		return
	}

	return
}

func (CCMPyACTR) Info() *framework.Info {
	return &Info
}

func (CCMPyACTR) ValidateModel(model *actr.Model) (log *issues.Log) {
	log = issues.New()

	if model.Memory.LatencyExponent != nil {
		log.Warning(nil, "ccm does not support memory module's latency_exponent")
	}

	for _, production := range model.Productions {
		if production.DoStatements != nil {
			for _, statement := range production.DoStatements {
				if (statement.Recall != nil) && (len(statement.Recall.RequestParameters) > 0) {
					keys := maps.Keys(statement.Recall.RequestParameters)
					location := issues.Location{
						Line:        production.AMODLineNumber,
						ColumnStart: 0,
						ColumnEnd:   0,
					}

					log.Warning(&location,
						"ccm does not support request parameters (%q in %q)",
						strings.Join(keys, ", "), production.Name)
				}
			}
		}
	}

	return
}

// SetModel sets our model and saves the python class name we are going to use.
func (c *CCMPyACTR) SetModel(model *actr.Model) (err error) {
	if model.Name == "" {
		err = framework.ErrModelMissingName
		return
	}

	c.model = model
	c.className = fmt.Sprintf("ccm_%s", c.model.Name)

	return
}

func (c CCMPyACTR) Model() (model *actr.Model) {
	return c.model
}

// Run generates the python code from the amod file, writes it to disk, creates a "run" file
// to actually run the model, and returns the output (stdout and stderr combined).
func (c *CCMPyACTR) Run(options *runoptions.Options) (result *framework.RunResult, err error) {
	runFile, err := c.WriteModel(c.tmpPath, options)
	if err != nil {
		return
	}

	result = &framework.RunResult{
		FileName:      runFile,
		GeneratedCode: c.GetContents(),
	}

	output, err := executil.ExecCommand(Info.ExecutableName, runFile)
	if err != nil {
		return
	}

	result.Output = []byte(output)

	return
}

// WriteModel converts the internal actr.Model to Python and writes it to a file.
func (c *CCMPyACTR) WriteModel(path string, options *runoptions.Options) (outputFileName string, err error) {
	// If our model has a print statement, then write out our support file
	if c.model.HasPrintStatement() {
		err = framework.WriteSupportFile(path, ccmPrintFileName, ccmPrintPython)
		if err != nil {
			return
		}
	}

	// If our model is tracing activations, then write out our support file
	if *options.TraceActivations {
		err = framework.WriteSupportFile(path, gactarActivateTraceFileName, gactarActivateTraceFile)
		if err != nil {
			return
		}
	}

	outputFileName = fmt.Sprintf("%s.py", c.className)
	if path != "" {
		outputFileName = fmt.Sprintf("%s/%s", path, outputFileName)
	}

	err = filesystem.RemoveFile(outputFileName)
	if err != nil {
		return "", err
	}

	_, err = c.GenerateCode(options)
	if err != nil {
		return
	}

	err = c.WriteFile(outputFileName)
	if err != nil {
		return
	}

	return
}

// GenerateCode converts the internal actr.Model to Python code.
func (c *CCMPyACTR) GenerateCode(options *runoptions.Options) (code []byte, err error) {
	patterns, err := framework.ParseInitialBuffers(c.model, options.InitialBuffers)
	if err != nil {
		return
	}

	goal := patterns["goal"]

	err = c.InitWriterHelper()
	if err != nil {
		return
	}

	c.writeHeader()

	memory := c.model.Memory

	c.writeImports(options)

	c.Write("\n\n")

	// random
	if options.RandomSeed != nil {
		c.Writeln("random.seed(%d)", *options.RandomSeed)
		c.Write("\n\n")
	}

	c.Writeln("class %s(ACTR):", c.className)

	for _, buffer := range c.model.BufferNames() {
		c.Writeln("    %s = Buffer()", buffer)
	}

	c.Writeln("")

	additionalInit := []string{}

	if memory.LatencyFactor != nil {
		additionalInit = append(additionalInit, fmt.Sprintf("latency=%s", numbers.Float64Str(*memory.LatencyFactor)))
	}

	if memory.RetrievalThreshold != nil {
		additionalInit = append(additionalInit, fmt.Sprintf("threshold=%s", numbers.Float64Str(*memory.RetrievalThreshold)))
	}

	if memory.FinstSize != nil {
		additionalInit = append(additionalInit, fmt.Sprintf("finst_size=%d", *memory.FinstSize))
	}

	if memory.FinstTime != nil {
		additionalInit = append(additionalInit, fmt.Sprintf("finst_time=%s", numbers.Float64Str(*memory.FinstTime)))
	}

	if len(additionalInit) > 0 {
		c.Writeln("    %s = Memory(%s, %s)", memory.ModuleName(), memory.BufferName(), strings.Join(additionalInit, ", "))
	} else {
		c.Writeln("    %s = Memory(%s)", memory.ModuleName(), memory.BufferName())
	}

	if *options.TraceActivations {
		c.Writeln("    trace = ActivateTrace(%s)", memory.ModuleName())
	}

	c.Writeln("")

	if memory.IsUsingBaseLevelLearning() {
		c.Writeln("    DMBaseLevel(%s, decay=%s)", memory.ModuleName(), numbers.Float64Str(*memory.Decay))
		c.Writeln("")
	}

	// Turn on DMSpreading if we are using spreading activation
	c.writeSpreadingActivation()

	// Turn on DMNoise if we have set "instantaneous_noise"
	if memory.InstantaneousNoise != nil {
		c.Writeln("    DMNoise(%s, noise=%s)", memory.ModuleName(), numbers.Float64Str(*memory.InstantaneousNoise))
		c.Writeln("")
	}
	// Turn on Partial if we have set "mismatch_penalty"
	if memory.MismatchPenalty != nil {
		c.Writeln("    partial = Partial(%[1]s, limit=%s)", memory.ModuleName(), numbers.Float64Str(*memory.MismatchPenalty))
		c.writeSimilarities("partial")
		c.Writeln("")
	}

	procedural := c.model.Procedural

	if procedural.DefaultActionTime != nil {
		c.Writeln("    production_time = %s", numbers.Float64Str(*procedural.DefaultActionTime))
		c.Writeln("")
	}

	if c.model.HasPrintStatement() {
		c.Writeln("    # create a printer helper and register chunks with their slots for lookup")
		c.Writeln("    printer = CCMPrint()")

		for _, chunk := range c.model.Chunks {
			quotedSlotNames := []string{}

			for _, slot := range chunk.SlotNames {
				quotedSlotNames = append(quotedSlotNames, fmt.Sprintf("%q", slot))
			}
			c.Writeln("    printer.register_chunk(%q, [%s])", chunk.TypeName, strings.Join(quotedSlotNames, ", "))
		}

		c.Writeln("")
	}

	if *options.LogLevel == "info" {
		// this turns on some logging at the high level
		c.Writeln("    def __init__(self):")
		c.Writeln("        super().__init__(log=True)")
		c.Writeln("")
	}

	c.writeInitializers(goal)

	c.Writeln("")

	// Add user-set goal if any
	if goal != nil {
		c.Write("        goal.set(")
		c.outputPattern(goal)
		c.Write(")\n\n")
	}

	c.writeProductions()

	c.Writeln("")

	c.writeMain(options)

	code = c.GetContents()
	return
}

func (c CCMPyACTR) writeHeader() {
	c.Writeln("\"\"\"")

	if c.model.Description != "" {
		c.Write("%s\n\n", c.model.Description)
	}

	c.writeAuthors()

	c.Writeln("Generated by gactar %s", framework.GactarVersion)
	c.Writeln("          https://github.com/asmaloney/gactar")
	c.Writeln("          on %s", framework.TimeNow().Format("2006-01-02 @ 15:04:05"))
	c.Writeln("")
	c.Writeln("NOTE: This is a generated file. Any changes may be overwritten.")

	c.Writeln("\"\"\"\n")
}

func (c CCMPyACTR) writeAuthors() {
	if len(c.model.Authors) == 0 {
		return
	}

	c.Writeln("Authors:")

	for _, author := range c.model.Authors {
		c.Writeln("   %s", author)
	}

	c.Writeln("")
}

func (c CCMPyACTR) writeImports(runOptions *runoptions.Options) {
	if runOptions.RandomSeed != nil {
		c.Writeln("import random")
	}

	memory := c.model.Memory

	imports := []string{"ACTR", "Buffer", "Memory"}

	c.Write("from python_actr import %s\n", strings.Join(imports, ", "))

	additionalImports := []string{}

	if memory.IsUsingBaseLevelLearning() {
		additionalImports = append(additionalImports, "DMBaseLevel")
	}

	if memory.IsUsingSpreadingActivation() {
		additionalImports = append(additionalImports, "DMSpreading")
	}

	if memory.InstantaneousNoise != nil {
		additionalImports = append(additionalImports, "DMNoise")
	}

	if memory.MismatchPenalty != nil {
		additionalImports = append(additionalImports, "Partial")
	}

	if len(additionalImports) > 0 {
		c.Write("from python_actr import %s\n", strings.Join(additionalImports, ", "))
	}

	if *runOptions.LogLevel == "detail" {
		c.Writeln("from python_actr import log, log_everything")
	}

	if c.model.HasPrintStatement() {
		c.Writeln("")
		c.Writeln(fmt.Sprintf("from %s import CCMPrint", ccmPrintImportName))
	}

	if *runOptions.TraceActivations {
		c.Writeln("")
		c.Writeln(fmt.Sprintf("from %s import ActivateTrace", gactarActivateTraceImportName))
	}
}

// If spreading activation is on, write its parameters
func (c CCMPyACTR) writeSpreadingActivation() {
	memory := c.model.Memory

	if !memory.IsUsingSpreadingActivation() {
		return
	}

	c.Writeln("    spread = DMSpreading(%s, %s)", memory.ModuleName(), strings.Join(c.model.BufferNames(), ", "))
	c.Writeln("    spread.strength = %s", numbers.Float64Str(*memory.MaxSpreadStrength))

	for _, buffer := range c.model.Buffers() {
		bufferName := buffer.Name()

		c.Writeln("    spread.weight[%s] = %s", bufferName, numbers.Float64Str(buffer.SpreadingActivation()))
	}

	c.Writeln("")
}

func (c CCMPyACTR) writeInitializers(goal *actr.Pattern) {
	if len(c.model.Initializers) == 0 {
		return
	}

	c.Writeln("    def init():")

	for _, init := range c.model.Initializers {
		module := init.Module

		// allow the user-set goal to override the initializer
		if module.ModuleName() == "goal" && (goal != nil) {
			continue
		}

		c.Write("        # amod line %d", init.AMODLineNumber)
		if init.ChunkName != nil {
			c.Write(" %q", *init.ChunkName)
		}
		c.Writeln("")

		buffer := init.Buffer
		if module.AllowsMultipleInit() {
			c.Write("        %s.add(", module.ModuleName())
		} else {
			c.Write("        %s.set(", buffer.Name())
		}

		c.outputPattern(init.Pattern)
		c.Writeln(")")
	}
}

func (c CCMPyACTR) writeSimilarities(partialName string) {
	if len(c.model.Similarities) == 0 {
		return
	}

	for _, similar := range c.model.Similarities {
		c.Writeln("    # amod line %d", similar.AMODLineNumber)
		c.Writeln("    %s.similarity('%s', '%s', %s)", partialName, similar.ChunkOne, similar.ChunkTwo, numbers.Float64Str(similar.Value))
	}
}

func (c CCMPyACTR) writeProductions() {
	for _, production := range c.model.Productions {
		if production.Description != nil {
			c.Writeln("    # %s", *production.Description)
		}

		c.Writeln("    # amod line %d", production.AMODLineNumber)

		c.Write("    def %s(", production.Name)

		numMatches := len(production.Matches)
		for i, match := range production.Matches {
			c.outputMatch(match)

			if i != numMatches-1 {
				c.Write(", ")
			}
		}

		c.Writeln("):")

		if production.DoStatements != nil {
			for _, statement := range production.DoStatements {
				c.outputStatement(statement)
			}
		}

		c.Write("\n")
	}
}

func (c CCMPyACTR) writeMain(runOptions *runoptions.Options) {
	c.Writeln("if __name__ == \"__main__\":")
	c.Writeln(fmt.Sprintf("    model = %s()", c.className))

	if *runOptions.LogLevel == "detail" {
		c.Writeln("    log(summary=1)")
		c.Writeln("    log_everything(model)")
	}

	c.Writeln("    model.run()")
}

func (c CCMPyACTR) outputPattern(pattern *actr.Pattern) {
	str := fmt.Sprintf("'%s ", pattern.Chunk.TypeName)

	for i, slot := range pattern.Slots {
		str += patternSlotString(slot)

		if i != len(pattern.Slots)-1 {
			str += " "
		}
	}

	str += "'"

	c.Write(str)
}

func (c CCMPyACTR) outputMatch(match *actr.Match) {
	switch {
	case match.BufferPattern != nil:
		bufferName := match.BufferPattern.Buffer.Name()

		c.Write("%s=", bufferName)
		if match.BufferPattern.Pattern.AnyChunk {
			c.Write("'?'")
		} else {
			c.outputPattern(match.BufferPattern.Pattern)
		}

	case match.BufferState != nil:
		bufferName := match.BufferState.Buffer.Name()

		c.Write("%s='%s:True'", bufferName, match.BufferState.State)

	case match.ModuleState != nil:
		moduleName := match.ModuleState.Module.ModuleName()

		c.Write("%s='%s:True'", moduleName, match.ModuleState.State)
	}
}

func patternSlotString(slot *actr.PatternSlot) string {
	var str string

	if slot.Negated {
		str += "!"
	}

	switch {
	case slot.Wildcard:
		str += "?"

	case slot.Nil:
		str += "None"

	case slot.ID != nil:
		str += *slot.ID

	case slot.Str != nil:
		str += strings.ReplaceAll(*slot.Str, " ", "_")

	case slot.Var != nil:
		str += *slot.Var.Name

	case slot.Num != nil:
		str += *slot.Num
	}

	// Check for constraints on a var and output them
	if slot.Var != nil {
		if len(slot.Var.Constraints) > 0 {
			for _, constraint := range slot.Var.Constraints {
				if constraint.Comparison == actr.NotEqual {
					str += "!"
				}

				str += convertValue(constraint.RHS)
			}
		}
	}

	return str
}

func (c CCMPyACTR) outputStatement(s *actr.Statement) {
	switch {
	case s.Set != nil:
		if s.Set.Slots != nil {
			slotAssignments := []string{}
			for _, slot := range *s.Set.Slots {
				value := convertValue(slot.Value)
				slotAssignments = append(slotAssignments, fmt.Sprintf("_%d=%s", slot.SlotIndex, value))
			}
			c.Writeln("        %s.modify(%s)", s.Set.Buffer.Name(), strings.Join(slotAssignments, ", "))
		} else {
			c.Write("        %s.set(", s.Set.Buffer.Name())
			c.outputPattern(s.Set.Pattern)
			c.Writeln(")")
		}

	case s.Recall != nil:
		c.Write("        %s.request(", s.Recall.MemoryModuleName)
		c.outputPattern(s.Recall.Pattern)
		c.Writeln(")")

	case s.Clear != nil:
		for _, name := range s.Clear.BufferNames {
			c.Writeln("        %s.clear()", name)
		}

	case s.Print != nil:
		if s.Print.IsBufferOutput() {
			id := *((*s.Print.Values)[0].ID)
			ids := strings.Split(id, ".")

			if len(ids) == 1 {
				c.Writeln("        printer.print_chunk(%s, %q)", id, id)
			} else {
				c.Writeln("        printer.print_chunk_slot(%s, %q, %q)", ids[0], ids[0], ids[1])
			}
		} else {
			values := pythonValuesToStrings(s.Print.Values, true)
			c.Writeln("        print(%s, sep='')", strings.Join(values, ", "))
		}

	case s.Stop != nil:
		c.Writeln("        self.stop()")
	}
}

func convertValue(s *actr.Value) string {
	switch {
	case s.Nil != nil:
		return "None"

	case s.Var != nil:
		return *s.Var

	case s.ID != nil:
		return "'" + *s.ID + "'"

	case s.Number != nil:
		return *s.Number

	case s.Str != nil:
		return "'" + strings.ReplaceAll(*s.Str, " ", "_") + "'"
	}

	return ""
}

func pythonValuesToStrings(values *[]*actr.Value, quoteStrings bool) []string {
	str := make([]string, len(*values))
	for i, v := range *values {
		switch {
		case v.Var != nil:
			str[i] = strings.TrimPrefix(*v.Var, "?")

		case v.Str != nil:
			if quoteStrings {
				str[i] = fmt.Sprintf("'%s'", *v.Str)
			} else {
				str[i] = *v.Str
			}

		case v.Number != nil:
			str[i] = *v.Number
		}
		// v.ID && v.Nil should not be possible because of validation
	}

	return str
}
