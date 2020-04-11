package os

import "os"

type files struct {}

func (f files) IsFile(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	mode := fi.Mode()
	isFile := mode.IsRegular()

	return isFile, nil
}