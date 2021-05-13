package subcommands

import (
	"errors"
	"github.com/bobappleyard/anathema/server/a"
	"golang.org/x/tools/go/packages"
	"path"
	"path/filepath"
)

var ErrNoPackage = errors.New("no package")

type moduleFinder interface {
	OwningModule(dir string) (*packages.Module, error)
}

type importFileCompiler interface {
	CompileImportFile(importingPkg string, mod *packages.Module) error
}

type applicationCommand struct {
	a.Service

	Loader     packageLoader
	Finder     moduleFinder
	Components componentFileCompiler
	Imports    importFileCompiler
}

func (c *applicationCommand) Name() string {
	return "application"
}

func (c *applicationCommand) Run(args []string) error {
	mod, err := c.Finder.OwningModule(".")
	if err != nil {
		return err
	}

	err = c.createPackageInitFiles(mod)
	if err != nil {
		return err
	}

	thisDir, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	thisPkg := filepath.Join(mod.Path, thisDir[len(mod.Dir):])

	err = c.Imports.CompileImportFile(thisPkg, mod)
	if err != nil {
		return err
	}

	return nil
}

func (c *applicationCommand) createPackageInitFiles(mod *packages.Module) error {
	pkgs, err := c.Loader.LoadPackages(mod.Path + "/...")
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		dir := path.Join(mod.Dir, pkg.PkgPath[len(mod.Path):])
		err := c.Components.CompileComponentFile(dir, pkg)
		if err != nil {
			return err
		}
	}

	return nil
}
