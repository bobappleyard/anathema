package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/bobappleyard/anathema/server/a"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type packageLoader interface {
	LoadPackages(patterns ...string) ([]*packages.Package, error)
}

type moduleAnalyzer struct {
	a.Service

	Loader packageLoader
}

func (a *moduleAnalyzer) OwningModule(dir string) (*packages.Module, error) {
	adir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("go", "list", "-m", "-json")
	cmd.Dir = adir
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var m packages.Module
	err = json.Unmarshal(data, &m)
	return &m, nil
}

func (a *moduleAnalyzer) AllAvailablePackages(mod *packages.Module) ([]*packages.Package, error) {
	mods, err := a.ParseModuleDependencies(mod)
	if err != nil {
		return nil, err
	}

	modules := make([]string, len(mods)+1)
	for i, m := range mods {
		modules[i] = m + "/..."
	}
	modules[len(mods)] = mod.Path + "/..."

	return a.Loader.LoadPackages(modules...)
}

func (p *moduleAnalyzer) ParseModuleDependencies(mod *packages.Module) ([]string, error) {
	f, err := os.Open(mod.GoMod)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return p.scanFile(f)
}

func (p *moduleAnalyzer) scanFile(f io.Reader) ([]string, error) {
	var deps []string

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		deps = p.scanLine(deps, bytes.NewReader(scanner.Bytes()))
	}

	return deps, scanner.Err()
}

func (p *moduleAnalyzer) scanLine(deps []string, f io.Reader) []string {
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return deps
	}

	if scanner.Text() != "require" {
		return deps
	}

	if !scanner.Scan() {
		return deps
	}

	return append(deps, scanner.Text())
}
