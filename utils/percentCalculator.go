package utils

import (
	"io"
)

// helps calculating file completion percentage
type PercentCalculator struct {
	OutputSize int
	InputSize  int
	Start      float64
	End        float64
}

// create new instance of [utils.PercentCalculator] for input paths
func NewPathsCalculator(start float64) *PercentCalculator {
	return &PercentCalculator{
		Start: start,
		End:   100.0,
	}
}

// create new instance of [utils.PercentCalculator] for [io.Seeker]
func NewSeekerCalculator(inputSeeker io.Seeker, start float64) (pc *PercentCalculator, err error) {
	pc = &PercentCalculator{
		Start: start,
		End:   100.0,
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
	percent := pc.Start + (float64(pc.OutputSize) / float64(pc.InputSize) * pc.End)

	if pc.End > percent {
		return percent
	}

	return pc.End
}
