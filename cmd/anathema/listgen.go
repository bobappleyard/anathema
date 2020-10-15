package main

import (
	"fmt"
	"go/types"
	"io"
)

func generateTypeList(marked []types.Object, dest io.Writer) error {
	fmt.Fprintln(dest, "package main")

	packages := map[string]bool{}

	for _, t := range marked {
		packages[t.Pkg().Path()] = true
	}

	fmt.Fprintln(dest, "import (")
	fmt.Fprintln(dest, "\"reflect\"")
	fmt.Fprintln(dest, "\"github.com/bobappleyard/anathema/typereg\"")
	for p := range packages {
		fmt.Fprintf(dest, "%q\n", p)
	}
	fmt.Fprintln(dest, ")")

	fmt.Fprintln(dest, "func init() {")
	for _, t := range marked {
		fmt.Fprintf(dest, "typereg.RegisterType(reflect.TypeOf(new(%s.%s)).Elem())\n", t.Pkg().Name(), t.Name())
	}
	fmt.Fprintln(dest, "}")

	return nil
}
