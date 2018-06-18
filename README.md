The jamaica pacakge provides some handy step definitions and functions for in-process testing of CLI applications via the [godog](https://github.com/DATA-DOG/godog) Cucumber port.

In short, this is an attempt to implement a very minimal [aruba](https://github.com/cucumber/aruba) for Go.

## Installation ##

It would be best to use [dep](https://github.com/golang/dep) to include jamaica in your dependency tree, but there is absolutely no reason that you shouldn't be
able to simply do the standard `go get`:

```
go get github.com/ess/jamaica
```

## Configuration ##

The same steps are necessary regardless of the way that your `godog` suite is set up:

1. Use `jamaica.SetRootCmd()` to declare the command that you want to test. The object passed must adhere to the [`jamaica.Command`]() interface.
2. Use `jamaica.StepUp()` to inject the jamaica step definitions into your `godog` suite.

### Basic Godog Configuration ###

If you're not using `TestMain` to execute your `godog` suite, you would set up jamaica within the `FeatureContext` in one of your test files:

```go
package something

import (
  "github.com/DATA-DOG/godog"
  "github.com/ess/jamaica"
)

func FeatureContext(s *godog.Suite) {
  jamaica.SetRootCmd(myCommand)
  jamaica.StepUp(s)
  
  // Register the rest of your steps
}
```

### TestMain Configuration ###

If you're using `TestMain` to execute your `godog` suite, you would set up jamaica like so:

```go
package main

import (
  "os"
  "testing"

  "github.com/DATA-DOG/godog"
  "github.com/ess/jamaica"
)

func TestMain(m *testing.M) {
  jamaica.SetRootCmd(cmd.RootCmd)

  status := godog.RunWithOptions(
    "godog",

    func(s *godog.Suite) {
      jamaica.StepUp(s)
      
      // Register the rest of your steps
    },

    godog.Options{
      Format: "pretty",
      Paths: []string{"features"},
    },
  )

  if st := m.Run(); st > status {
    status = st
  }

  os.Exit(status)
}
```

## Provided Steps ##

A fairly minimal set of steps are provided by jamaica, because it's not exactly aruba.

### Running the Command ###

After you've set up all of your givens and such, you can use the following
step to run the command under test in-process:

```gherkin
When I run `mycommand and its arguments`
```

### Checking the Output ###

To assert that the command output contains a given string, use the following step:

```gherkin
Then stdout contains "a given string"
```

To assert that the command output exactly matches a given string, use the following step:

```gherkin
Then stdout is "a given string"
```

If these steps don't suit your needs (or if you'd prefer to use your own terminology), you can access the command's stdout via `jamaica.LastCommandStdout()` for your comparisons.

### Checking the Exit Status ###

To assert the expectation that the application ran successfully (without error), use the following step:

```gherkin
Then it exits successfully
```

On the other hand, to assert that you expect for the command to have failed, use this:

```gherkin
Then it exits with an error
```

## Usage ##

Simply run your test suite via either `godog` or by using godog with the "test main" pattern.

## History ##

* v1.0.2 - Added missing steps
* v1.0.1 - Fixing some typos
* v1.0.0 - First stable release
* v0.0.8 - No hard deps
* v0.0.7 - Minor tweaks
* v0.0.6 - And now you can get the command output
* v0.0.5 - Now maybe even working
* v0.0.4 - Root the command
* v0.0.3 - Always be returning
* v0.0.2 - Or maybe I can cargo cult my own code properly
* v0.0.1 - Initial attempt
