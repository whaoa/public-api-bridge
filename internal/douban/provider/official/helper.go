package official

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

const (
	requestUrl      = "https://m.douban.com/"
	requestProxyUrl = "https://m.douban.cmliussss.net/"
)

var (
	timeZoneChina = time.FixedZone("CST", 60*60*8)
)

func normalizeDateTime(dt string) time.Time {
	layout := time.DateTime
	if !strings.Contains(dt, " ") {
		layout = time.DateOnly
	}
	if t, err := time.ParseInLocation(layout, dt, timeZoneChina); err == nil {
		return t.UTC()
	}
	return time.UnixMilli(0)
}

func normalDuration(duration string) int64 {
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

func withProxy(url string) string {
	if strings.HasPrefix(url, requestUrl) {
		return strings.Replace(url, requestUrl, requestProxyUrl, 1)
	}
	if strings.HasPrefix(url, "/") {
		return requestProxyUrl + strings.TrimPrefix(url, "/")
	}
	return url
}
