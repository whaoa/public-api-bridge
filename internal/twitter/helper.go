package twitter

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/samber/lo"
)

func NormalizeDateTime(dt string) time.Time {
	if t, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", dt); err == nil {
		return t
	}
	return time.UnixMilli(0)
}

func buildTweetMosaicImageUrl(tweet string, media []*Attachment) string {
	if len(media) == 0 {
		return ""
	}

	names := lo.FilterMap(media, func(m *Attachment, _ int) (string, bool) {
		if m.Kind != MediaTypePhoto {
			return "", false
		}

		u, err := url.Parse(m.Url)
		if err != nil {
			return "", false
		}

		return strings.TrimSuffix(path.Base(u.Path), path.Ext(u.Path)), true
	})

	if len(names) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"https://mosaic.fxtwitter.com/jpeg/%s/%s",
		tweet,
		strings.Join(names, "/"),
	)
}

func fillUserEntity(user *User) *User {
	if user != nil {
		user.Url = fmt.Sprintf("https://x.com/%s", user.ScreenName)
	}
	return user
}

func fillTweetEntity(tweet *Tweet) *Tweet {
	if tweet != nil {
		tweet.Author.Url = fmt.Sprintf("https://x.com/%s", tweet.Author.ScreenName)
		tweet.Url = fmt.Sprintf("https://x.com/%s/status/%s", tweet.Author.ScreenName, tweet.ID)
	}
	if tweet != nil && tweet.Quote != nil {
		tweet.Quote.Url = fmt.Sprintf("https://x.com/%s/status/%s", tweet.Quote.ScreenName, tweet.Quote.Tweet)
	}
	if tweet != nil && tweet.ReplyingTo != nil {
		tweet.ReplyingTo.Url = fmt.Sprintf("https://x.com/%s/status/%s", tweet.ReplyingTo.ScreenName, tweet.ReplyingTo.Tweet)
	}

	if tweet != nil && len(tweet.Attachments) > 0 {
		lo.ForEach(tweet.Attachments, func(m *Attachment, _ int) {
			if m.Thumbnail != "" || m.Kind != MediaTypePhoto {
				return
			}
			if u, err := url.Parse(m.Url); err == nil {
				query := u.Query()
				query.Set("name", "small")
				u.RawQuery = query.Encode()
				m.Thumbnail = u.String()
			} else {
				m.Thumbnail = m.Url
			}
		})

		tweet.MosaicPhoto = buildTweetMosaicImageUrl(tweet.ID, tweet.Attachments)
	}

	return tweet
}
