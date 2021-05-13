package util

import (
	"github.com/bobappleyard/anathema/server/a"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
)

type packageAnalyzer struct {
	a.Service
}

func (a *packageAnalyzer) LoadPackages(patterns ...string) ([]*packages.Package, error) {
	mode := packages.NeedImports | packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo
	config := &packages.Config{
		Mode: mode,
		Fset: token.NewFileSet(),
	}
	return packages.Load(config, patterns...)
}

func (a *packageAnalyzer) TypesAssignableTo(pkg *packages.Package, marker types.Type) []types.Object {
	var res []types.Object

	if pkg.PkgPath == marker.(*types.Named).Obj().Pkg().Path() {
		return nil
	}
	s := pkg.Types.Scope()
	for _, n := range s.Names() {
		object, ok := s.Lookup(n).(*types.TypeName)
		if !ok {
			continue
		}
		if object.IsAlias() {
			continue
		}
		if !types.AssignableTo(object.Type(), marker) {
			continue
		}
		res = append(res, object)
	}

	return res
}
