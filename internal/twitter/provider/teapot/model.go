package teapot

import (
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/whaoa/public-api-bridge/internal/twitter"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/jq"
)

func buildUser(reader io.Reader) (*twitter.User, error) {
	doc, err := jq.New(reader)
	if err != nil {
		return nil, err
	}

	id := ""
	banner := normalizeAttachmentUrl(doc.Find(".profile-banner img").AttrOr("src", ""))
	avatar := normalizeAttachmentUrl(doc.Find(".profile-card-avatar > img").AttrOr("src", ""))
	name := doc.Find(".profile-card-fullname").AttrOr("title", "")
	screenName := doc.Find(".profile-card-username").AttrOr("title", "")
	description := doc.Find(".profile-bio").Text()
	location := doc.Find(".profile-location").Text()

	matches := regexp.MustCompile(`/profile_banners/(\d+)/`).FindStringSubmatch(banner)
	if len(matches) > 1 {
		id = matches[1]
	}

	joinedAt := time.UnixMilli(0)
	if text := doc.Find(".profile-joindate > span").AttrOr("title", ""); text != "" {
		if t, err := time.Parse("15:04 PM - 02 Jan 2006", text); err == nil {
			joinedAt = t
		}
	}

	following := doc.Find(".profile-statlist > .following > .profile-stat-num").AttrOr("title", "")
	followers := doc.Find(".profile-statlist > .followers > .profile-stat-num").AttrOr("title", "")

	verified := doc.Find(".profile-card-identity .verified-icon").Length() > 0
	protected := doc.Find(".profile-card-identity .icon-lock").Length() > 0

	return &twitter.User{
		ID:          id,
		Name:        name,
		ScreenName:  screenName,
		Avatar:      avatar,
		Banner:      banner,
		Location:    location,
		Description: description,
		JoinedAt:    joinedAt,
		Verified:    verified,
		Protected:   protected,
		Following:   normalizeNumberString(following),
		Followers:   normalizeNumberString(followers),
	}, nil
}

func buildTweet(id string, reader io.Reader) (*twitter.Tweet, error) {
	doc, err := jq.New(reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tweet := doc.Find(".main-tweet")

	name := tweet.Find(".fullname-and-username > .fullname").Text()
	screenName := strings.TrimSuffix(tweet.Find(".fullname-and-username > .username").Text(), "@")
	avatar := normalizeAttachmentUrl(tweet.Find(".tweet-avatar > img").AttrOr("src", ""))
	content := tweet.Find(".tweet-content").Text()

	media := []*twitter.Attachment{}
	tweet.Find(".attachment img").Each(func(i int, s *goquery.Selection) {
		media = append(media, &twitter.Attachment{
			Kind: twitter.MediaTypePhoto,
			Url:  normalizeAttachmentUrl(s.AttrOr("src", "")),
		})
	})
	tweet.Find(".attachment video > source").Each(func(i int, s *goquery.Selection) {
		media = append(media, &twitter.Attachment{
			Kind:      twitter.MediaTypeVideo,
			Url:       normalizeAttachmentUrl(s.AttrOr("src", "")),
			Thumbnail: normalizeAttachmentUrl(s.Parent().AttrOr("poster", "")),
		})
	})

	replies := normalizeNumberString(tweet.Find(".tweet-stat .icon-comment").Parent().Text())
	retweets := normalizeNumberString(tweet.Find(".tweet-stat .icon-retweet").Parent().Text())
	views := normalizeNumberString(tweet.Find(".tweet-stat .icon-heart").Parent().Text())
	likes := normalizeNumberString(tweet.Find(".tweet-stat .icon-views").Parent().Text())

	createAt := time.Now()
	if text := tweet.Find(".tweet-published").Text(); text != "" {
		if time, err := time.Parse("Jan 01, 2006 · 15:04 PM UTC", text); err == nil {
			createAt = time
		}
	}

	quote := linkToBasicTweet(tweet.Find(".quote-link").First().AttrOr("href", ""))
	to := linkToBasicTweet(doc.Find(".before-tweet .tweet-link").Last().AttrOr("href", ""))

	return &twitter.Tweet{
		ID:                id,
		Content:           content,
		Attachments:       media,
		Replies:           replies,
		Retweets:          retweets,
		Likes:             likes,
		Views:             views,
		Quote:             quote,
		ReplyingTo:        to,
		CreatedAt:         createAt,
		PossiblySensitive: false,
		Author: &twitter.BasicUser{
			ScreenName: screenName,
			Name:       name,
			Avatar:     avatar,
		},
	}, nil
}
