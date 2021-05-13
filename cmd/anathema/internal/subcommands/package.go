package subcommands

import (
	"errors"
	"github.com/bobappleyard/anathema/server/a"
	"golang.org/x/tools/go/packages"
)

var ErrUnableToLoadPackage = errors.New("unable to load package")

type packageLoader interface {
	LoadPackages(patterns ...string) ([]*packages.Package, error)
}

type componentFileCompiler interface {
	CompileComponentFile(dir string, pkg *packages.Package) error
}

type packageCommand struct {
	a.Service

	Compiler componentFileCompiler
	Loader   packageLoader
}

func (c *packageCommand) Name() string {
	return "package"
}

func (c *packageCommand) Run(args []string) error {
	pkgs, err := c.Loader.LoadPackages(".")
	if err != nil {
		return err
	}
	if len(pkgs) != 1 {
		return ErrUnableToLoadPackage
	}
	return c.Compiler.CompileComponentFile(".", pkgs[0])
}
