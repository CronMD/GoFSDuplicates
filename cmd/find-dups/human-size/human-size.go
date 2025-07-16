package humansize

import (
	"fmt"
	"math"
)

func SizeToString(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%d b", size)
	} else if size < 1_000_000 {
		return fmt.Sprintf("%d Kb", int64(math.Round(float64(size)/1_000)))
	} else if size < 1_000_000_000 {
		return fmt.Sprintf("%d Mb", int64(math.Round(float64(size)/1_000_000)))
	}

	return fmt.Sprintf("%.1f Gb", float64(size)/1_000_000_000)
}
