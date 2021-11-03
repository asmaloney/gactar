package shell

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli/v2"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"
)

type command struct {
	description string
	method      func(string) error
}

type Shell struct {
	context          *cli.Context
	history          []string
	currentModel     *actr.Model
	actrFrameworks   *framework.List
	activeFrameworks map[string]bool
	commands         map[string]command
}

func Initialize(cli *cli.Context, frameworks *framework.List) (s *Shell, err error) {
	s = &Shell{
		context:          cli,
		actrFrameworks:   frameworks,
		activeFrameworks: map[string]bool{},
	}

	s.preamble()

	for name, framework := range *frameworks {
		err = framework.Initialize()
		if err != nil {
			return nil, err
		}

		s.activeFrameworks[name] = true
	}

	s.commands = map[string]command{
		"frameworks": {`choose frameworks to run (e.g. "ccm pyactr", "all")`, s.cmdFramework},
		"history":    {"outputs your command history", s.cmdHistory},
		"load":       {"loads a model: load [FILENAME]", s.cmdLoad},
		"reset":      {"resets the current model", s.cmdReset},
		"run":        {"runs the current model: run [INITIAL STATE]", s.cmdRun},
		"version":    {"outputs version info", s.cmdVersion},

		"help": {"exits the program", s.cmdHelp},
		"exit": {"exits the program", s.cmdExit},
		"quit": {"exits the program", s.cmdExit},
	}

	return
}

func (s *Shell) Start() (err error) {
	if err != nil {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		cmd = strings.TrimSpace(cmd)

		if cmd == "" {
			continue
		}

		s.history = append(s.history, cmd)

		err = s.runCommand(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, " error: %s\n", err)
		}
	}
}

func (s *Shell) preamble() {
	cli.ShowVersion(s.context)
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

	err = fmt.Errorf("unrecognized command: '%s'", cmd)

	return
}

func (s *Shell) cmdFramework(fNames string) (err error) {
	names := strings.Split(fNames, " ")
	sort.Strings(names)

	if names[0] == "all" {
		names = s.actrFrameworks.Names()
		sort.Strings(names)
	}

	s.activeFrameworks = map[string]bool{}

	for _, name := range names {
		if !s.actrFrameworks.Exists(name) {
			err = fmt.Errorf("'%s' is not a valid framework. Valid values: %v", name, s.actrFrameworks.Names())
			return
		}

		s.activeFrameworks[name] = true
	}

	if len(s.activeFrameworks) == 0 {
		err = fmt.Errorf("no frameworks selected. Valid values: %v", framework.ValidFrameworks)
		return
	}

	fmt.Print("active frameworks: ")
	for name := range s.activeFrameworks {
		fmt.Printf("%s ", name)
	}
	fmt.Println()

	return
}

func (s *Shell) cmdHistory(string) (err error) {
	fmt.Println(strings.Join(s.history, "\n"))
	return
}

func (s *Shell) cmdLoad(fileName string) (err error) {
	if fileName == "" {
		err = fmt.Errorf("'load' requires a file name")
		return
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

	return
}

func (s *Shell) cmdReset(string) (err error) {
	s.currentModel = nil
	fmt.Println(" model reset")
	return
}

func (s *Shell) cmdRun(initialGoal string) (err error) {
	if s.currentModel == nil {
		err = fmt.Errorf("no model loaded")
		return
	}

	for name, f := range *s.actrFrameworks {
		if !s.activeFrameworks[name] {
			continue
		}

		err = f.SetModel(s.currentModel)
		if err != nil {
			return err
		}

		initialBuffers := framework.InitialBuffers{
			"goal": strings.TrimSpace(initialGoal),
		}

		_, output, err := f.Run(initialBuffers)
		if err != nil {
			return err
		}

		fmt.Print(string(output))

		if output[len(output)-1] != '\n' {
			fmt.Println()
		}
	}

	return
}

func (s *Shell) cmdVersion(string) (err error) {
	cli.ShowVersion(s.context)
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
