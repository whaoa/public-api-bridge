package official

import (
	"strings"
	"time"
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

func withProxy(url string) string {
	if strings.HasPrefix(url, requestUrl) {
		return strings.Replace(url, requestUrl, requestProxyUrl, 1)
	}
	if strings.HasPrefix(url, "/") {
		return requestProxyUrl + strings.TrimPrefix(url, "/")
	}
	return url
}
