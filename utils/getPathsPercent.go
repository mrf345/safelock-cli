package utils

import (
	"io/fs"
	"os"
)

// compare the input paths with the output path total file size and calculates a percentage
func GetPathsPercent(
	inputPaths []string,
	outputPath string,
	start float64,
	portion float64,
) (percent float64, err error) {
	var inputInfo, outputInfo fs.FileInfo
	var inputSize, outputSize = int64(0), int64(0)

	if outputInfo, err = os.Stat(outputPath); err != nil {
		return
	}

	for _, inputPath := range inputPaths {
		if inputInfo, err = os.Stat(inputPath); err != nil {
			return
		}

		if inputInfo.IsDir() {
			if inputSize, err = GetDirSize(inputPath); err != nil {
				return
			}
		} else {
			inputSize += inputInfo.Size()
		}

	}

	if outputInfo.IsDir() {
		if outputSize, err = GetDirSize(outputPath); err != nil {
			return
		}
	} else {
		outputSize += outputInfo.Size()
	}

	percent = start + (float64(outputSize) / float64(inputSize) * portion)

	return
}
