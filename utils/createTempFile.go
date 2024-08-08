package utils

import (
	"fmt"
	"os"
)

func CreateTempFile(path string) (tempPath string, tempFile *os.File, err error) {
	tempPath = fmt.Sprintf("%s.temp", path)

	if tempFile, err = os.Create(tempPath); err != nil {
		return
	}

	return
}
