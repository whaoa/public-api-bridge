package douban

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	matherSupportImageUrl   = regexp.MustCompile(`^https://(qnmob\d-sign|img\d)\.doubanio\.com/`)
	matherResizablePhotoUrl = regexp.MustCompile("/view/photo/([a-z_]+)/public/")

	placeholderPhotoUrlSize = "/view/photo/%s/public/"
)

func withImageProxy(src string) string {
	if matherSupportImageUrl.MatchString(src) {
		if strings.Contains(src, "-sign.doubanio.com/") {
			if values := strings.Split(src, "?"); len(values) > 0 {
				src = values[0]
			}
		}
		return matherSupportImageUrl.ReplaceAllString(src, "https://img.doubanio.cmliussss.net/")
	}
	return src
}

func withImageResize(src string, size string) string {
	if matches := matherResizablePhotoUrl.FindStringSubmatch(src); len(matches) > 0 {
		return strings.Replace(src, matches[0], fmt.Sprintf(placeholderPhotoUrlSize, size), 1)
	}
	return src
}
