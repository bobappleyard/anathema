package services

import (
	"github.com/bobappleyard/anathema/component"
	"github.com/bobappleyard/anathema/server/a"
	"go/types"
	"golang.org/x/tools/go/packages"
	"path"
)

type typeSearcher interface {
	TypesAssignableTo(pkg *packages.Package, marker types.Type) []types.Object
}

type componentFileGenerator interface {
	GenerateComponentFile(name, path string, components []string) error
}

type fileRemover interface {
	Remove(path string) error
}

type componentFileCompiler struct {
	a.Service

	Searcher  typeSearcher
	Generator componentFileGenerator
	Remover   fileRemover
}

func (s *componentFileCompiler) CompileComponentFile(dir string, pkg *packages.Package) error {
	var components []string
	for _, t := range s.Searcher.TypesAssignableTo(pkg, component.CompileType()) {
		components = append(components, t.Name())
	}
	if len(components) == 0 {
		_ = s.Remover.Remove(path.Join(dir, "anathema_components.go"))
		return nil
	}
	return s.Generator.GenerateComponentFile(pkg.Name, dir, components)
}
