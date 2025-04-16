package utils

import (
	"os"
)

// OpenFile() opens a file at the specified path. If the file does not exist, it will be created.
//
// The file is opened with write-only and append-only permissions, and the file permissions are set to 0777. It returns the opened file and any error that occurred during the process.
func OpenFile(path string) (*os.File, error) {
	var file *os.File
	var err error

	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return file, nil
}
