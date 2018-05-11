package jamaica

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Command interface {
	SetArgs([]string)
	Execute() error
}

type Suite interface {
	Step(interface{}, interface{})
	BeforeScenario(func(interface{}))
}

var rootCmd Command
var commandOutput string
var lastCommandRanErr error

func iRun(fullCommand string) error {
	if rootCmd == nil {
		return fmt.Errorf("You must set the root command via jamaica.SetRootCmd before running Jamaica steps.")
	}

	args := strings.Split(fullCommand, " ")[1:]

	rootCmd.SetArgs(args)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lastCommandRanErr = rootCmd.Execute()
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	commandOutput = <-outC

	return nil
}

func theCommandSucceeds() error {
	if lastCommandRanErr != nil {
		return fmt.Errorf(
			"Expected a good exit status, got '%s'",
			lastCommandRanErr.Error(),
		)
	}

	return nil
}

func theCommandFails() error {
	if lastCommandRanErr == nil {
		return fmt.Errorf(
			"Expected a bad exit status, got nil",
		)
	}

	return nil
}

func StepUp(s Suite) {
	irunregex := fmt.Sprintf(`^I run %s([^%s]*)%s$`, "`", "`", "`")
	s.Step(irunregex, iRun)
	s.Step(`the command succeeds`, theCommandSucceeds)
	s.Step(`the command fails`, theCommandFails)

	s.BeforeScenario(func(interface{}) {
		commandOutput = ""
		lastCommandRanErr = nil
	})
}

func SetRootCmd(cmd Command) {
	rootCmd = cmd
}

func LastCommandOutput() string {
	return commandOutput
}
