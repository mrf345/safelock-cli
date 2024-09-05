package utils

import (
	"io"

	"github.com/mholt/archiver/v4"
)

// helps calculating file completion percentage
type PercentCalculator struct {
	OutputSize int
	InputSize  int
	start      float64
	end        float64
}

// create new instance of [utils.PercentCalculator] for input paths
func NewPathsCalculator(start float64, files []archiver.File) *PercentCalculator {
	pc := &PercentCalculator{
		start: start,
		end:   100.0,
	}
	pc.setInputPathsSize(files)
	return pc
}

func (pc *PercentCalculator) setInputPathsSize(files []archiver.File) {
	for _, file := range files {
		if !file.FileInfo.IsDir() {
			pc.InputSize += int(file.FileInfo.Size())
		}
	}
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
