package cli

import (
	"context"
	"fmt"
	"io"
	"os"
)

// App is the main structure of a CLI application. It is recommended that
// an app be created with the cli.NewApp() function.
type App struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string
	// Full name of command for help, defaults to Name
	HelpName string
	// Description of the program.
	Usage string
	// Text to override the USAGE section of help
	UsageText string
	// Description of the program argument format.
	ArgsUsage string
	// Version of the program
	Version string
	// Description of the program
	Description string
	// List of commands to execute
	Commands []*Command
	// List of flags to parse
	Flags []Flag
	// Boolean to enable bash completion commands
	EnableBashCompletion bool
	// Boolean to hide built-in help command
	HideHelp bool
	// Boolean to hide built-in help flag
	HideHelpFlag bool
	// Boolean to hide built-in version flag and the VERSION section of help
	HideVersion bool
	// Writer writer to write output to
	Writer io.Writer
	// ErrWriter writes error output
	ErrWriter io.Writer
	// Execute this function if the proper command cannot be found
	CommandNotFound CommandNotFoundFunc
	// Execute this function if a usage error occurs
	OnUsageError OnUsageErrorFunc
	// Execute this function when an action is not explicitly defined
	Action ActionFunc
	// Execute this function before any subcommands are run, but after the context is ready
	Before BeforeFunc
	// Execute this function after any subcommands are run, but after the subcommand has finished
	After AfterFunc
	// rootCommand is the root level command for the app
	rootCommand *Command
}

// CommandNotFoundFunc is executed if the proper command cannot be found.
type CommandNotFoundFunc func(ctx context.Context, cmd *Command, command string)

// OnUsageErrorFunc is executed if a usage error occurs. This is useful for displaying
// customized usage error messages.
type OnUsageErrorFunc func(ctx context.Context, cmd *Command, err error, isSubcommand bool) error

// BeforeFunc is an action to execute before any subcommands are run, but after
// the context is ready. If a non-nil error is returned, no subcommands are run.
type BeforeFunc func(ctx context.Context, cmd *Command) error

// AfterFunc is an action to execute after any subcommands are run, but after
// the subcommand has finished.
type AfterFunc func(ctx context.Context, cmd *Command) error

// ActionFunc is the action to execute when no subcommands are specified.
type ActionFunc func(ctx context.Context, cmd *Command) error

// NewApp creates a new cli Application with some reasonable defaults for Name
// and Usage.
func NewApp() *App {
	return &App{
		Name: os.Args[0],
		HelpName: os.Args[0],
		Usage: "A new cli application",
		Version: "0.0.0",
		Writer: os.Stdout,
		ErrWriter: os.Stderr,
	}
}

// Run is the entry point to the cli app. Parses the arguments slice and routes
// to the proper flag/args combination.
func (a *App) Run(ctx context.Context, arguments []string) error {
	a.rootCommand = a.newRootCommand()
	return a.rootCommand.Run(ctx, arguments)
}

// newRootCommand creates a root Command from the App definition.
func (a *App) newRootCommand() *Command {
	return &Command{
		Name:                   a.Name,
		HelpName:               a.HelpName,
		Usage:                  a.Usage,
		UsageText:              a.UsageText,
		ArgsUsage:              a.ArgsUsage,
		Version:                a.Version,
		Description:            a.Description,
		Commands:               a.Commands,
		Flags:                  a.Flags,
		EnableBashCompletion:   a.EnableBashCompletion,
		HideHelp:               a.HideHelp,
		HideHelpFlag:           a.HideHelpFlag,
		HideVersion:            a.HideVersion,
		Writer:                 a.Writer,
		ErrWriter:              a.ErrWriter,
		CommandNotFound:        a.CommandNotFound,
		OnUsageError:           a.OnUsageError,
		Action:                 a.Action,
		Before:                 a.Before,
		After:                  a.After,
		isRoot:                 true,
	}
}

// ToMarkdown generates a markdown string for the `*App`
// The function errors if either parsing or writing of the string fails.
func (a *App) ToMarkdown() (string, error) {
	a.rootCommand = a.newRootCommand()
	if a.rootCommand == nil {
		return "", fmt.Errorf("root command not found")
	}
	return a.rootCommand.ToMarkdown()
}

// ToMan generates a man page string for the `*App`
// The function errors if either parsing or writing of the string fails.
func (a *App) ToMan() (string, error) {
	a.rootCommand = a.newRootCommand()
	if a.rootCommand == nil {
		return "", fmt.Errorf("root command not found")
	}
	return a.rootCommand.ToMan()
}
