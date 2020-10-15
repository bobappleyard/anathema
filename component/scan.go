package component

import (
	"go/token"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/packages"
)

// Marker interface that represents something to be scanned for.
//
// To have anathema detect your types, embed them with this marker type. You can
// embed them in other interface types to turn them into component markers too.
// This is how, for example, the Resource marker works.
type Marker interface {
	marker()
}

var markerType = reflect.TypeOf(new(Marker)).Elem()

// Scan searches for types that embed the marker (directly or not) in the
// provided package, along with every package that it imports transitively.
func Scan(pkg string) ([]types.Object, error) {
	pkgs, err := loadPackages(pkg)
	if err != nil {
		return nil, err
	}

	marker := findMarker(pkgs)

	// If we can't find the marker then it must not have been imported, so there
	// cannot be any marked types.
	if marker == nil {
		return nil, nil
	}

	return searchPackages(pkgs, marker), nil
}

func loadPackages(pkg string) ([]*packages.Package, error) {
	mode := packages.NeedImports | packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo
	config := &packages.Config{
		Mode: mode,
		Fset: token.NewFileSet(),
	}
	return packages.Load(config, pkg)
}

func searchPackages(pkgs []*packages.Package, marker types.Type) []types.Object {
	var res []types.Object

	packages.Visit(pkgs, func(p *packages.Package) bool {
		s := p.Types.Scope()
		for _, n := range s.Names() {
			object, ok := s.Lookup(n).(*types.TypeName)
			if !ok {
				continue
			}
			if !types.AssignableTo(object.Type(), marker) {
				continue
			}
			res = append(res, object)
		}
		return true
	}, nil)

	return res
}

func findMarker(pkgs []*packages.Package) types.Type {
	var res types.Type

	packages.Visit(pkgs, func(p *packages.Package) bool {
		if p.PkgPath != markerType.PkgPath() {
			return true
		}

		object := p.Types.Scope().Lookup(markerType.Name())
		if object == nil {
			return false
		}

		tbind, ok := object.(*types.TypeName)
		if !ok {
			return false
		}

		res = tbind.Type()
		return false
	}, nil)

	return res
}
