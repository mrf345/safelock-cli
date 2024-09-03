package utils

import (
	"io"
	"io/fs"
	"os"
)

// helps calculating file completion percentage
type PercentCalculator struct {
	OutputSize int
	InputSize  int
	start      float64
	end        float64
}

// create new instance of [utils.PercentCalculator] for input paths
func NewPathsCalculator(inputPaths []string, start float64) (pc *PercentCalculator, err error) {
	pc = &PercentCalculator{
		start: start,
		end:   100.0,
	}
	err = pc.setInputPathsSize(inputPaths)
	return
}

func (pc *PercentCalculator) setInputPathsSize(paths []string) (err error) {
	for _, path := range paths {
		var info fs.FileInfo

		if info, err = os.Stat(path); err != nil {
			return
		}

		if info.IsDir() {
			if pc.InputSize, err = GetDirSize(path); err != nil {
				return
			}
		} else {
			pc.InputSize += int(info.Size())
		}
	}

	return
}

// create new instance of [utils.PercentCalculator] for [io.Seeker]
func NewSeekerCalculator(inputSeeker io.Seeker, start float64) (pc *PercentCalculator, err error) {
	pc = &PercentCalculator{
		start: start,
		end:   100.0,
	}
	err = pc.setInputSeekerSize(inputSeeker)
	return
}

func (pc *PercentCalculator) setInputSeekerSize(seeker io.Seeker) (err error) {
	size, err := seeker.Seek(0, io.SeekEnd)
	pc.InputSize = int(size)
	return
}

// get current completion percentage
func (pc *PercentCalculator) GetPercent() float64 {
	return pc.start + (float64(pc.OutputSize) / float64(pc.InputSize) * pc.end)
}
