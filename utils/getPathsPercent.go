package utils

import (
	"fmt"
	"io/fs"
	"os"
)

// compare the input path with the output path total file size and calculates a percentage
func GetPathsPercent(inputPath, outputPath string, start float64, portion float64) (percent string, err error) {
	var inputInfo, outputInfo fs.FileInfo
	var inputSize, outputSize int64

	if inputInfo, err = os.Stat(inputPath); err != nil {
		return
	}

	if outputInfo, err = os.Stat(outputPath); err != nil {
		return
	}

	if inputInfo.IsDir() {
		if inputSize, err = GetDirSize(inputPath); err != nil {
			return
		}
	} else {
		inputSize = inputInfo.Size()
	}

	if outputInfo.IsDir() {
		if outputSize, err = GetDirSize(outputPath); err != nil {
			return
		}
	} else {
		outputSize = outputInfo.Size()
	}

	total := start + (float64(outputSize) / float64(inputSize) * portion)
	percent = fmt.Sprintf("%.2f%%", total)

	return
}
