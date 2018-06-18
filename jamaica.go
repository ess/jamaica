// Copyright © 2018 Dennis Walters
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jamaica provides handy godog steps for testing CLI applications
package jamaica

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// Command is an interface that describes a command that we want to test
// in-process.
//
// As this command will be run in-process (rather than via Exec()), it is
// expected to have a SetArgs() method to set its CLI arguments.
type Command interface {
	SetArgs([]string)
	Execute() error
}

// Suite is an interface that describes the parts of the godog.Suite interface
// used within the package
type Suite interface {
	Step(interface{}, interface{})
	BeforeScenario(func(interface{}))
}

// SetRootCmd takes a Command and registers it as the app that we want to
// test in-process. If this is not set, calls to the "I run ..." step will
// fail with a relevant message.
func SetRootCmd(cmd Command) {
	rootCmd = cmd
}

// StepUp takes a Suite and injects the jamaica step definitions into it:
//  * When I run `somecommand`
//  * Then it exits successfully
//  * Then it exits with an error
//  * Then stdout contains "some string"
//  * Then stdout is "some string"
func StepUp(s Suite) {
	irunregex := fmt.Sprintf(`^I run %s([^%s]*)%s$`, "`", "`", "`")
	s.Step(irunregex, iRun)
	s.Step(`^it exits successfully$`, theCommandSucceeds)
	s.Step(`^it exits with an error$`, theCommandFails)
	s.Step(`^stdout contains "([^"]*)"$`, stdoutContains)
	s.Step(`^stdout is "([^"]*)"$`, stdoutIs)

	s.BeforeScenario(func(interface{}) {
		commandStdout = ""
		lastCommandRanErr = nil
	})
}

// LastCommandStdout returns a string representation of the output (including
// newlines) of a command run with the "I run" step
func LastCommandStdout() string {
	return commandStdout
}

var rootCmd Command
var commandStdout string
var lastCommandRanErr error

func iRun(fullCommand string) error {
	if rootCmd == nil {
		return fmt.Errorf("jamaica.SetRootCmd must be set before running Jamaica steps")
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
	commandStdout = <-outC

	return nil
}

func stdoutContains(s string) error {
	if !strings.Contains(commandStdout, s) {
		return fmt.Errorf(`Expected stdout to contain "%s"`, s)
	}

	return nil
}

func stdoutIs(s string) error {
	if commandStdout != s {
		return fmt.Errorf(`Expected stdout to contain exactly "%s"`, s)
	}

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
