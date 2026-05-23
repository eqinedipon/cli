// Package cli provides a framework for building command line applications.
// It is a fork of urfave/cli with additional features and improvements.
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// ActionFunc is the type for action handlers.
type ActionFunc func(ctx context.Context, cmd *Command) error

// BeforeFunc is the type for before-action handlers.
type BeforeFunc func(ctx context.Context, cmd *Command) (context.Context, error)

// AfterFunc is the type for after-action handlers.
type AfterFunc func(ctx context.Context, cmd *Command) error

// Command represents a CLI command or subcommand.
type Command struct {
	// Name is the name of the command.
	Name string

	// Usage is a short description of the command.
	Usage string

	// UsageText overrides the auto-generated usage text.
	UsageText string

	// Description is a longer description of the command.
	Description string

	// Version is the version of the command (typically used on root command).
	Version string

	// Aliases are alternative names for the command.
	Aliases []string

	// Flags are the flags associated with this command.
	Flags []Flag

	// Commands are the subcommands.
	Commands []*Command

	// Action is the function to call when the command is invoked.
	Action ActionFunc

	// Before is called before Action.
	Before BeforeFunc

	// After is called after Action.
	After AfterFunc

	// Writer is where output is written (defaults to os.Stdout).
	Writer io.Writer

	// ErrWriter is where error output is written (defaults to os.Stderr).
	ErrWriter io.Writer

	// HideHelp hides the help command and flag.
	HideHelp bool

	// HideVersion hides the version flag.
	HideVersion bool

	// parent is the parent command, if any.
	parent *Command

	// flagSet is the underlying flag set.
	flagSet *flag.FlagSet
}

// Run executes the command with the given arguments.
func (c *Command) Run(ctx context.Context, args []string) error {
	if err := c.setup(); err != nil {
		return err
	}

	if len(args) == 0 {
		args = os.Args
	}

	// Parse the flags
	if err := c.flagSet.Parse(args[1:]); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	remaining := c.flagSet.Args()

	// Check for subcommand
	if len(remaining) > 0 {
		subName := remaining[0]
		for _, sub := range c.Commands {
			if sub.Name == subName || containsString(sub.Aliases, subName) {
				sub.parent = c
				return sub.Run(ctx, append([]string{subName}, remaining[1:]...))
			}
		}
	}

	// Run Before hook
	if c.Before != nil {
		var err error
		ctx, err = c.Before(ctx, c)
		if err != nil {
			return err
		}
	}

	// Run the action
	if c.Action != nil {
		if err := c.Action(ctx, c); err != nil {
			return err
		}
	}

	// Run After hook
	if c.After != nil {
		if err := c.After(ctx, c); err != nil {
			return err
		}
	}

	return nil
}

// setup initializes the command's flag set and registers flags.
func (c *Command) setup() error {
	if c.Name == "" {
		c.Name = os.Args[0]
	}

	if c.Writer == nil {
		c.Writer = os.Stdout
	}

	if c.ErrWriter == nil {
		c.ErrWriter = os.Stderr
	}

	c.flagSet = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	c.flagSet.SetOutput(c.ErrWriter)

	for _, f := range c.Flags {
		if err := f.Apply(c.flagSet); err != nil {
			return fmt.Errorf("error applying flag %q: %w", f.GetName(), err)
		}
	}

	return nil
}

// String returns the full command name including parent names.
func (c *Command) String() string {
	if c.parent != nil {
		return strings.Join([]string{c.parent.String(), c.Name}, " ")
	}
	return c.Name
}

// containsString checks if a slice contains a given string.
func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
