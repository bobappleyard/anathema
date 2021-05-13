package services

import (
	"github.com/bobappleyard/anathema/component"
	"github.com/bobappleyard/anathema/server/a"
	"golang.org/x/tools/go/packages"
	"strings"
)

type moduleAnalyzer interface {
	AllAvailablePackages(mod *packages.Module) ([]*packages.Package, error)
}

type importFileGenerator interface {
	GenerateImportFile(name, path string, imports []string) error
}

type importFileCompiler struct {
	a.Service

	Package   string `env:"GOPACKAGE"`
	Analyzer  moduleAnalyzer
	Generator importFileGenerator
	Searcher  typeSearcher
}

func (c *importFileCompiler) CompileImportFile(importingPkg string, mod *packages.Module) error {
	paths, err := c.packagesToImport(importingPkg, mod)
	if err != nil {
		return err
	}
	return c.Generator.GenerateImportFile(c.Package, ".", paths)
}

func (c *importFileCompiler) packagesToImport(importingPkg string, mod *packages.Module) ([]string, error) {
	pkgs, err := c.Analyzer.AllAvailablePackages(mod)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, pkg := range pkgs {
		if pkg.PkgPath == importingPkg {
			continue
		}
		if pkg.Name == "main" {
			continue
		}
		if !c.packageVisible(importingPkg, pkg.PkgPath) {
			continue
		}
		if len(c.Searcher.TypesAssignableTo(pkg, component.CompileType())) == 0 {
			continue
		}
		paths = append(paths, pkg.PkgPath)
	}

	return paths, nil
}

func (c *importFileCompiler) packageVisible(importing string, imported string) bool {
	internalPos := strings.Index(imported, "/internal/")
	if internalPos == -1 {
		return true
	}
	return importing[:internalPos] == imported[:internalPos]
}

