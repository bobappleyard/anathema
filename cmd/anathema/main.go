package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/bobappleyard/anathema/component"
)

const mainFile = `
package main

import "github.com/bobappleyard/anathema/server"

func main() {
	server.Run()
}
`

func main() {
	modpath, err := findModulePath()
	if err != nil {
		panic(err)
	}

	tmp, err := ioutil.TempDir(os.TempDir(), "*")
	if err != nil {
		panic(err)
	}
	err = prepareTempDir(tmp, modpath)
	if err != nil {
		panic(err)
	}

	fmt.Println(tmp)

	types, err := component.Scan(".")
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	generateTypeList(types, &out)
	bs, err := format.Source(out.Bytes())
	if err != nil {
		panic(err)
	}

	gen := filepath.Join(tmp, ".anathema")
	err = os.Mkdir(gen, 0777)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(filepath.Join(gen, "init.go"), bs, 0777)
	ioutil.WriteFile(filepath.Join(gen, "main.go"), []byte(mainFile), 0777)

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = gen
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func findModulePath() (string, error) {
	cmd := exec.Command("go", "env")
	env, err := cmd.Output()
	if err != nil {
		return "", err
	}
	prefix := []byte("GOMOD=")
	for _, line := range bytes.Split(env, []byte("\n")) {
		if !bytes.HasPrefix(line, prefix) {
			continue
		}
		modpath, err := strconv.Unquote(string(line[len(prefix):]))
		if err != nil {
			return "", err
		}
		return filepath.Dir(modpath), nil
	}
	return "", errors.New("no module found")
}

func prepareTempDir(tmp, modpath string) error {
	dir, err := os.Open(modpath)
	if err != nil {
		return err
	}

	fs, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, f := range fs {
		err = os.Symlink(
			filepath.Join(modpath, f.Name()),
			filepath.Join(tmp, f.Name()),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
