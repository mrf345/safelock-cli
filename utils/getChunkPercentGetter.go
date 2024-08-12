package utils

import (
	"fmt"
	"io/fs"
	"os"
)

// helps calculate file chunks percentage
type ChunkPercentCalculator struct {
	File      *os.File
	ChunkSize int
	Start     float64
	Portion   float64
	counter   int
}

// calculate the percent of file chunks
func (c *ChunkPercentCalculator) GetPercent() (percent string, err error) {
	var stat fs.FileInfo

	if stat, err = c.File.Stat(); err != nil {
		return
	}

	c.counter += 1
	total := float64(int(stat.Size()) / c.ChunkSize)
	total = float64(c.counter) / total * c.Portion
	percent = fmt.Sprintf("%.2f%%", c.Start+total)

	return
}
