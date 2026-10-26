package util

import (
	"math"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func NormalDuration(duration string) int64 {
	if duration != "" && duration != "0:0" {
		parts := strings.Split(duration, ":")
		length := len(parts)
		return lo.Reduce(parts, func(r int64, s string, i int) int64 {
			v, _ := strconv.ParseInt(s, 10, 64)
			return r + v*int64(math.Pow(60, float64(length-i-1)))
		}, 0)
	}
	return 0
}
