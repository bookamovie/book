package utils

import (
	"os"
)

func OpenLog(path string) (*os.File, error) {
	var file *os.File
	var err error

	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return file, nil
}
