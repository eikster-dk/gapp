package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Reader struct{}

type ReadFile interface {
	ReadFile(path string) ([]byte, error)
}

func (f Reader) IsFile(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	mode := fi.Mode()
	isFile := mode.IsRegular()

	return isFile, nil
}

func (f Reader) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (f Reader) ReadDir(path string) ([]string, error) {
	ff, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, f := range ff {
		paths = append(paths, filepath.Join(path, f.Name()))
	}

	return paths, nil
}
