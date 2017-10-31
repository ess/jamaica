package jamaica

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command
var commandOutput string
var lastCommandRanErr error

func iRun(fullCommand string) error {
	args := strings.Split(fullCommand, " ")[1:]

	cmd.RootCmd.SetArgs(args)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lastCommandRanErr = cmd.RootCmd.Execute()
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

func StepUp(s *godog.Suite) {
	rootCmd = cmd

	s.Step(`^I run "([^"]*)"$`, iRun)
	s.Step(`the command succeeds`, theCommandSucceeds)
	s.Step(`the command fails`, theCommandFails)

	s.BeforeScenario(func(interface{}) {
		commandOutput = ""
		lastCommandRanErr = nil
	})
}

func SetRootCmd(c *cobra.Command) {
	rootCmd = c
}
