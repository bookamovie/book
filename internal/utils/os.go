package utils

import (
	"errors"
	"os"
)

func OpenLog(path string) (*os.File, error) {
	var file *os.File
	var err error

	file, err = os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			file, err = os.Create(path)
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	err = file.Chmod(0777)
	if err != nil {
		return nil, err
	}

	return file, nil
}
