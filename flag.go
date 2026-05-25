package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Flag is the interface that all flag types must implement.
type Flag interface {
	// GetName returns the name of the flag.
	GetName() string
	// GetUsage returns the usage string for the flag.
	GetUsage() string
	// IsRequired returns whether the flag is required.
	IsRequired() bool
	// Apply applies the flag to the given FlagSet.
	Apply(*flag.FlagSet) error
}

// StringFlag defines a flag with a string value.
type StringFlag struct {
	Name     string
	Aliases  []string
	Usage    string
	Value    string
	Required bool
	EnvVars  []string
	Destination *string
}

func (f *StringFlag) GetName() string  { return f.Name }
func (f *StringFlag) GetUsage() string { return f.Usage }
func (f *StringFlag) IsRequired() bool { return f.Required }

// Apply applies the StringFlag to the given FlagSet.
func (f *StringFlag) Apply(set *flag.FlagSet) error {
	val := f.Value
	for _, envVar := range f.EnvVars {
		if v, ok := lookupEnv(envVar); ok {
			val = v
			break
		}
	}
	if f.Destination != nil {
		set.StringVar(f.Destination, f.Name, val, f.Usage)
	} else {
		set.String(f.Name, val, f.Usage)
	}
	for _, alias := range f.Aliases {
		set.String(alias, val, fmt.Sprintf("alias for --%s", f.Name))
	}
	return nil
}

// BoolFlag defines a flag with a boolean value.
type BoolFlag struct {
	Name        string
	Aliases     []string
	Usage       string
	Value       bool
	Required    bool
	EnvVars     []string
	Destination *bool
}

func (f *BoolFlag) GetName() string  { return f.Name }
func (f *BoolFlag) GetUsage() string { return f.Usage }
func (f *BoolFlag) IsRequired() bool { return f.Required }

// Apply applies the BoolFlag to the given FlagSet.
func (f *BoolFlag) Apply(set *flag.FlagSet) error {
	val := f.Value
	for _, envVar := range f.EnvVars {
		if v, ok := lookupEnv(envVar); ok {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("could not parse %q as bool for flag %s: %w", v, f.Name, err)
			}
			val = b
			break
		}
	}
	if f.Destination != nil {
		set.BoolVar(f.Destination, f.Name, val, f.Usage)
	} else {
		set.Bool(f.Name, val, f.Usage)
	}
	for _, alias := range f.Aliases {
		set.Bool(alias, val, fmt.Sprintf("alias for --%s", f.Name))
	}
	return nil
}

// IntFlag defines a flag with an integer value.
type IntFlag struct {
	Name        string
	Aliases     []string
	Usage       string
	Value       int
	Required    bool
	EnvVars     []string
	Destination *int
}

func (f *IntFlag) GetName() string  { return f.Name }
func (f *IntFlag) GetUsage() string { return f.Usage }
func (f *IntFlag) IsRequired() bool { return f.Required }

// Apply applies the IntFlag to the given FlagSet.
func (f *IntFlag) Apply(set *flag.FlagSet) error {
	val := f.Value
	for _, envVar := range f.EnvVars {
		if v, ok := lookupEnv(envVar); ok {
			i, err := strconv.Atoi(strings.TrimSpace(v))
			if err != nil {
				return fmt.Errorf("could not parse %q as int for flag %s: %w", v, f.Name, err)
			}
			val = i
			break
		}
	}
	if f.Destination != nil {
		set.IntVar(f.Destination, f.Name, val, f.Usage)
	} else {
		set.Int(f.Name, val, f.Usage)
	}
	for _, alias := range f.Aliases {
		set.Int(alias, val, fmt.Sprintf("alias for --%s", f.Name))
	}
	return nil
}
