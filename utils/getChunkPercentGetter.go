package utils

import (
	"fmt"
	"io/fs"
	"os"
)

type PercentGetter = func(index int) string

func GetChunkPercentGetter(file *os.File, chunkSize int) (getter PercentGetter, err error) {
	var stat fs.FileInfo

	if stat, err = file.Stat(); err != nil {
		return
	}

	total := float64(int(stat.Size()) / chunkSize)
	getter = func(index int) string {
		return fmt.Sprintf("%.2f%%", 5+(float64(index+1)/total*90))
	}

	return
}
