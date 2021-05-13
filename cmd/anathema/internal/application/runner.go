package application

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
)

type Subcommand interface {
	Name() string
	Run(args []string) error
}

var ErrUnknownCommand = errors.New("unknown command")

type Runner struct {
	subs []Subcommand

	profile string
}

func (r *Runner) Inject(subs []Subcommand) {
	r.subs = subs
}

func (r *Runner) Run(args []string) error {
	rest, err := r.parseArgs(args)
	if err != nil {
		return err
	}
	if r.profile != "" {
		f, err := os.Create(r.profile)
		if err != nil {
			return err
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}
		defer pprof.StopCPUProfile()
	}
	return r.runCommand(rest)
}

func (r *Runner) runCommand(args []string) error {
	cmd := args[0]
	for _, sub := range r.subs {
		if cmd == sub.Name() {
			return sub.Run(args[1:])
		}
	}
	return fmt.Errorf("%w %s", ErrUnknownCommand, cmd)
}

func (r *Runner) parseArgs(args []string) ([]string, error) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringVar(&r.profile, "cpuprofile", "", "filename for cpu profile")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return fs.Args(), nil
}
