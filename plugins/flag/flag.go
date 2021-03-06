// Package flag provides flags support for uconfig
package flag

import (
	"flag"
	"os"
	"strings"

	"github.com/omeid/uconfig/flat"
)

const tag = "flag"

// Flags is it
type Flags interface {
	Visit(f flat.Fields) error

	Parse() error

	SetUsage(fn func())
}

// ErrorHandling defines how FlagSet.Parse behaves if the parse fails.
type ErrorHandling flag.ErrorHandling

//  These constants cause FlagSet.Parse to behave as described if the parse fails.
const (
	ContinueOnError = ErrorHandling(flag.ContinueOnError)
	ExitOnError     = ErrorHandling(flag.ExitOnError)
	PanicOnError    = ErrorHandling(flag.PanicOnError)
)

// New returns a new Flags
func New(name string, errorHandling ErrorHandling, args []string) Flags {
	return &visitor{
		fs:   flag.NewFlagSet(name, flag.ErrorHandling(errorHandling)),
		args: args,
	}
}

// Standard returns a set of flags configured in the common way.
// It is same as: `New(os.Args[0], ContinueOnError, os.Args[1:])`
func Standard() Flags {
	return New(os.Args[0], ContinueOnError, os.Args[1:])
}

var _ Flags = (*visitor)(nil)

type visitor struct {
	fs   *flag.FlagSet
	args []string
}

func (v *visitor) Parse() error {
	return v.fs.Parse(v.args)
}

func (v *visitor) Visit(fields flat.Fields) error {

	for _, f := range fields {
		usage, _ := f.Tag("usage")

		name, ok := f.Tag(tag)
		if name == "-" {
			continue
		}

		if !ok || name == "" {
			name = f.Name()
			name = strings.Replace(name, ".", "-", -1)
			name = strings.ToLower(name)
		}

		f.Meta()[tag] = "-" + name
		v.fs.Var(f, name, usage)
	}

	return nil

}

func (v *visitor) SetUsage(usage func()) {
	v.fs.Usage = usage
}
