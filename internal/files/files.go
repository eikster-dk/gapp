package files

import (
	"io/ioutil"
	"os"
)

type Reader struct {}

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

func (f Reader) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}