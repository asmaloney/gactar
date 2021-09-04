package shell

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli/v2"

	"gitlab.com/asmaloney/gactar/actr"
	"gitlab.com/asmaloney/gactar/amod"
	"gitlab.com/asmaloney/gactar/framework"
)

type command struct {
	description string
	method      func(string) error
}

type Shell struct {
	context       *cli.Context
	history       []string
	currentModel  *actr.Model
	actrFramework framework.Framework
	cmds          map[string]command
}

func Initialize(cli *cli.Context, framework framework.Framework) (s *Shell, err error) {
	s = &Shell{
		context:       cli,
		actrFramework: framework,
	}

	s.preamble()

	err = framework.Initialize()
	if err != nil {
		return nil, err
	}

	s.cmds = map[string]command{
		"history": {"outputs your command history", s.cmdHistory},
		"load":    {"loads a model: load [FILENAME]", s.cmdLoad},
		"run":     {"runs the current model: run [INITIAL STATE]", s.cmdRun},
		"reset":   {"resets the current model", s.cmdReset},
		"version": {"outputs version info", s.cmdVersion},

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

	if command, ok := s.cmds[cmd]; ok {
		err = command.method(args)
		return
	}

	err = fmt.Errorf("unrecognized command: '%s'", cmd)

	return
}

func (s *Shell) cmdExit(string) (err error) {
	os.Exit(0)
	return
}

func (s *Shell) cmdReset(string) (err error) {
	s.currentModel = nil
	fmt.Println(" model reset")
	return
}

func (s *Shell) cmdHelp(string) (err error) {
	// sort keys so commands may be output alphabetically
	keys := make([]string, 0, len(s.cmds))
	for k := range s.cmds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)
	for _, k := range keys {
		fmt.Fprintf(w, "  %v:\t%v\n", k, s.cmds[k].description)
	}
	w.Flush()

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

	s.currentModel, err = amod.GenerateModelFromFile(fileName)
	if err != nil {
		return err
	}

	fmt.Println(" model loaded")

	if s.currentModel.Examples != nil {
		fmt.Println(" examples:")

		for _, example := range s.currentModel.Examples {
			fmt.Printf("       run %s\n", example)
		}

	}

	return
}

func (s *Shell) cmdRun(initialGoal string) (err error) {
	if s.currentModel == nil {
		err = fmt.Errorf("no model loaded")
		return
	}

	err = s.actrFramework.SetModel(s.currentModel)
	if err != nil {
		return err
	}

	_, output, err := s.actrFramework.Run(initialGoal)
	if err != nil {
		return err
	}

	fmt.Print(string(output))

	if output[len(output)-1] != '\n' {
		fmt.Println()
	}

	return
}

func (s *Shell) cmdVersion(string) (err error) {
	cli.ShowVersion(s.context)
	return
}
