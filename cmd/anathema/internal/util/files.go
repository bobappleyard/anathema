package util

import (
	"github.com/bobappleyard/anathema/server/a"
	"io"
	"os"
)

type fileSystem struct {
	a.Service
}

func (s *fileSystem) Replace(path string) (io.WriteCloser, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
}

func (s *fileSystem) Remove(path string) error {
	return os.Remove(path)
}