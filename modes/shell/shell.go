// Package shell provides an interactive shell to load & run amod files.
package shell

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"golang.org/x/term"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"

	"github.com/asmaloney/gactar/util/chalk"
	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/issues"
	"github.com/asmaloney/gactar/util/runoptions"
	"github.com/asmaloney/gactar/util/validate"
)

var (
	ErrLoadRequiresName = errors.New("'load' requires a file name")
	ErrNoModel          = errors.New("no model loaded")
)

type ErrUnrecognizedCommand struct {
	Command string
}

func (e ErrUnrecognizedCommand) Error() string {
	return fmt.Sprintf("unrecognized command: %q", e.Command)
}

type ErrInvalidSetCommand struct {
	Command string
}

func (e ErrInvalidSetCommand) Error() string {
	return fmt.Sprintf("invalid set command: %q; expected form `set <option> <value>`", e.Command)
}

type ErrUnrecognizedSetOption struct {
	Option string
}

func (e ErrUnrecognizedSetOption) Error() string {
	return fmt.Sprintf("unrecognized option: %q; run `set` to see list of valid options", e.Option)
}

type ErrInvalidSetValue struct {
	OptionName string
	Value      string

	Expected string
}

func (e ErrInvalidSetValue) Error() string {
	return fmt.Sprintf("invalid value for %q: %q; expected %s", e.OptionName, e.Value, e.Expected)
}

type ErrNoFrameworkSelected struct {
	ValidFrameworks []string
}

func (e ErrNoFrameworkSelected) Error() string {
	valid := strings.Join(e.ValidFrameworks, ", ")
	return fmt.Sprintf("no framework selected; expected one of %q or \"all\"", valid)
}

type command struct {
	description string
	method      func(string) error
}

type Shell struct {
	settings   *cli.Settings
	runOptions runoptions.Options

	history          []string
	currentModel     *actr.Model
	activeFrameworks map[string]bool
	commands         map[string]command
}

func Initialize(settings *cli.Settings) (s *Shell, err error) {
	s = &Shell{
		settings:         settings,
		activeFrameworks: map[string]bool{},
	}

	s.preamble()

	for name := range settings.ActiveFrameworks {
		s.activeFrameworks[name] = true
	}

	s.commands = map[string]command{
		"frameworks": {`choose frameworks to run (e.g. "ccm pyactr", "all") - called without arguments, it will list active frameworks`, s.cmdFramework},
		"history":    {"outputs your command history", s.cmdHistory},
		"load":       {"loads a model: load [FILENAME]", s.cmdLoad},
		"reset":      {"resets the current model", s.cmdReset},
		"run":        {"runs the current model: run [INITIAL STATE]", s.cmdRun},
		"set":        {"set options: set [OPTION] [VALUE] - without arguments, lists options", s.cmdSet},
		"version":    {"outputs version info", s.cmdVersion},

		"help": {"outputs information about all available commands", s.cmdHelp},
		"exit": {"exits the program", s.cmdExit},
		"quit": {"exits the program", s.cmdExit},
	}

	return
}

func (s *Shell) Start() (err error) {
	if err != nil {
		return
	}

	termState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	defer func() {
		err = term.Restore(int(os.Stdin.Fd()), termState)
	}()

	terminal := term.NewTerminal(os.Stdin, "> ")

	for {
		line, err := terminal.ReadLine()
		if err != nil {
			break
		}

		err = term.Restore(int(os.Stdin.Fd()), termState)
		if err != nil {
			break
		}

		cmd := strings.TrimSpace(line)

		s.history = append(s.history, cmd)

		err = s.runCommand(cmd)
		if err != nil {
			chalk.PrintErr(err)
		}

		termState, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			break
		}
	}

	return
}

func (s *Shell) preamble() {
	fmt.Println("Type 'help' for a list of commands.")
	fmt.Println("To exit, type 'exit' or 'quit'.")
}

func (s *Shell) runCommand(c string) (err error) {
	space := strings.Index(c, " ")

	cmd := c
	args := ""

	if space != -1 {
		cmd = c[0:space]
		args = strings.TrimSpace(c[space+1:])
	}

	if command, ok := s.commands[cmd]; ok {
		err = command.method(args)
		return
	}

	return &ErrUnrecognizedCommand{Command: cmd}
}

func (s Shell) printActiveFrameworks() {
	fmt.Print("active frameworks: ")
	for name := range s.activeFrameworks {
		fmt.Printf("%s ", name)
	}
	fmt.Println()
}

func (s *Shell) cmdFramework(fNames string) (err error) {
	// If no frameworks were specified, then return list of active frameworks
	if len(fNames) == 0 {
		s.printActiveFrameworks()
		return
	}

	names := strings.Split(fNames, " ")
	sort.Strings(names)

	if names[0] == "all" {
		names = s.settings.ActiveFrameworks.Names()
		sort.Strings(names)
	}

	s.activeFrameworks = map[string]bool{}

	for _, name := range names {
		if !s.settings.ActiveFrameworks.Exists(name) {
			err = runoptions.ErrInvalidFrameworkName{
				Name:            name,
				ValidFrameworks: s.settings.ActiveFrameworks.Names(),
			}
			return
		}

		s.activeFrameworks[name] = true
	}

	// Should not be possible..
	if len(s.activeFrameworks) == 0 {
		err = ErrNoFrameworkSelected{ValidFrameworks: framework.ValidFrameworks}
		return
	}

	s.printActiveFrameworks()

	return
}

func (s *Shell) cmdHistory(string) (err error) {
	fmt.Println(strings.Join(s.history, "\n"))
	return
}

func (s *Shell) cmdLoad(fileName string) (err error) {
	if fileName == "" {
		return ErrLoadRequiresName
	}

	model, log, err := amod.GenerateModelFromFile(fileName)
	fmt.Print(log)
	if err != nil {
		return err
	}

	s.currentModel = model

	fmt.Println(" model loaded")

	if s.currentModel.Examples != nil {
		fmt.Println(" examples:")

		for _, example := range s.currentModel.Examples {
			fmt.Printf("       run %s\n", example)
		}
	}

	for name, f := range s.settings.ActiveFrameworks {
		if !s.activeFrameworks[name] {
			continue
		}

		log := f.ValidateModel(s.currentModel)
		if log.HasIssues() {
			fmt.Printf("== %s ==\n", f.Info().Name)
			fmt.Print(log)
			if log.HasError() {
				continue
			}
		}
	}

	return
}

func (s *Shell) cmdReset(string) (err error) {
	s.currentModel = nil
	fmt.Println(" model reset")
	return
}

func (s *Shell) cmdRun(initialGoal string) (err error) {
	if s.currentModel == nil {
		return ErrNoModel
	}

	log := issues.New()
	validate.Goal(s.currentModel, initialGoal, log)
	fmt.Print(log)

	for name, f := range s.settings.ActiveFrameworks {
		if !s.activeFrameworks[name] {
			continue
		}

		fmt.Printf("== %s ==\n", f.Info().Name)

		err = f.SetModel(s.currentModel)
		if err != nil {
			return err
		}

		options := s.currentModel.DefaultParams.Override(&s.runOptions)
		options.InitialBuffers = runoptions.InitialBuffers{
			"goal": strings.TrimSpace(initialGoal),
		}

		result, err := f.Run(options)
		if err != nil {
			return err
		}

		fmt.Print(string(result.Output))

		if result.Output[len(result.Output)-1] != '\n' {
			fmt.Println()
		}
	}

	return
}

func (s Shell) printActiveRunOptions() {
	notSet := chalk.Italic("<not set>")

	// logging
	option := notSet
	if s.runOptions.LogLevel != nil {
		option = string(*s.runOptions.LogLevel)
	}
	fmt.Printf("  %s %s (valid values are: %v)\n", chalk.Bold("logging"), option, strings.Join(runoptions.ACTRLoggingLevels, ", "))

	// trace
	option = notSet
	if s.runOptions.TraceActivations != nil {
		option = "off"
		if *s.runOptions.TraceActivations {
			option = "on"
		}
	}
	fmt.Printf("  %s %s (valid values are: on, off)\n", chalk.Bold("trace"), option)

	// random seed
	option = notSet
	if s.runOptions.RandomSeed != nil {
		option = fmt.Sprintf("%v", *s.runOptions.RandomSeed)
	}
	fmt.Printf("  %s %s\n", chalk.Bold("seed"), option)
}

func (s *Shell) cmdSet(args string) (err error) {
	if len(args) == 0 {
		s.printActiveRunOptions()
		return
	}

	options := strings.Split(args, " ")

	// we have at least one option
	optionName := options[0]

	if len(options) != 2 {
		return ErrInvalidSetCommand{
			Command: optionName,
		}
	}

	arg := options[1]

	switch optionName {
	case "logging":
		if !runoptions.ValidLogLevel(arg) {
			return runoptions.ErrInvalidLogLevel{Level: arg}
		}

		level := runoptions.ACTRLogLevel(arg)
		s.runOptions.LogLevel = &level

	case "trace":
		valid := []string{"on", "off"}
		if !slices.Contains(valid, arg) {
			return ErrInvalidSetValue{
				OptionName: optionName,
				Value:      arg,
				Expected:   fmt.Sprintf("one of %q", strings.Join(valid, ", ")),
			}
		}

		value := arg == "on"
		s.runOptions.TraceActivations = &value

	case "seed":
		value, err := strconv.ParseUint(arg, 10, 32)
		if err != nil {
			return ErrInvalidSetValue{
				OptionName: optionName,
				Value:      arg,
				Expected:   "a positive number",
			}
		}

		value32 := uint32(value)
		s.runOptions.RandomSeed = &value32

	default:
		return ErrUnrecognizedSetOption{Option: optionName}
	}

	return
}

func (s *Shell) cmdVersion(string) (err error) {
	fmt.Println(chalk.Bold(s.settings.Version))
	return
}

func (s *Shell) cmdHelp(string) (err error) {
	// sort keys so commands may be output alphabetically
	keys := make([]string, 0, len(s.commands))
	for k := range s.commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)
	for _, k := range keys {
		fmt.Fprintf(w, "  %v:\t%v\n", k, s.commands[k].description)
	}
	w.Flush()

	return
}

func (s *Shell) cmdExit(string) (err error) {
	os.Exit(0)
	return
}
