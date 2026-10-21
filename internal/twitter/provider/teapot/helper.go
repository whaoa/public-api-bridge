package teapot

import (
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/whaoa/public-api-bridge/internal/twitter"
)

func normalizeAttachmentUrl(src string) string {
	if src == "" {
		return src
	}
	if !strings.HasPrefix(src, "https://nitter.cf/") && !strings.HasPrefix(src, "/") {
		return src
	}
	if u, err := url.PathUnescape(path.Base(src)); err == nil {
		return u
	}
	return src
}

func normalizeNumberString(num string) int64 {
	num = strings.ReplaceAll(strings.TrimSpace(num), ",", "")
	n, _ := strconv.ParseInt(num, 10, 64)
	return n
}

func linkToBasicTweet(link string) *twitter.BasicTweet {
	re := regexp.MustCompile(`/([a-zA-Z\d-_]+)/status/(\d+)`)
	if matches := re.FindStringSubmatch(link); len(matches) > 2 {
		return &twitter.BasicTweet{
			ScreenName: matches[1],
			Tweet:      matches[2],
		}
	}
	return nil
}
