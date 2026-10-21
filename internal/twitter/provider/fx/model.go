package fx

import (
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"github.com/whaoa/public-api-bridge/internal/twitter"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

var errorInvalidResponse = errors.New("invalid fx-embed api response")

func buildUser(handle, payload string) (*twitter.User, error) {
	root := gjson.Parse(payload)

	code := root.Get("code").Int()

	if code == 404 {
		if root.Get("reason").String() == "suspended" {
			return nil, errors.Forbidden("the user is suspended")
		}
		return nil, errors.NotFound("the user does not exist")
	}

	if code != 200 {
		return nil, errors.Errorf("can not get profile for user @%s", handle)
	}

	user := root.Get("user")
	if !user.Exists() {
		return nil, errorInvalidResponse
	}

	return &twitter.User{
		ID:          user.Get("id").String(),
		Banner:      user.Get("banner_url").String(),
		Avatar:      user.Get("avatar_url").String(),
		Name:        user.Get("name").String(),
		ScreenName:  user.Get("screen_name").String(),
		Location:    user.Get("location").String(),
		Description: user.Get("description").String(),
		JoinedAt:    twitter.NormalizeDateTime(user.Get("joined").String()),
		Verified:    user.Get("verification.verified").Bool(),
		Protected:   user.Get("protected").Bool(),
		Following:   user.Get("following").Int(),
		Followers:   user.Get("followers").Int(),
	}, nil
}

func buildTweetSlice(handle, payload string) (*twitter.TweetSlice, error) {
	root := gjson.Parse(payload)
	code := root.Get("code").Int()

	if code == 404 {
		return nil, errors.NotFound("the user is unreachable")
	}

	if code != 200 {
		return nil, errors.Errorf("can not get tweets for user @%s", handle)
	}

	list := root.Get("results")
	if !list.IsArray() {
		return nil, errorInvalidResponse
	}

	return &twitter.TweetSlice{
		Items: lo.Map(list.Array(), func(t gjson.Result, _ int) *twitter.Tweet {
			return toTweet(t.Raw)
		}),
		PrevCursor: root.Get("cursor.top").String(),
		NextCursor: root.Get("cursor.bottom").String(),
	}, nil
}

func buildTweet(handle, id, payload string) (*twitter.Tweet, error) {
	root := gjson.Parse(payload)

	code := root.Get("code").Int()
	if code == 401 {
		return nil, errors.Forbidden("the tweet is unavailable")
	}
	if code == 404 {
		return nil, errors.NotFound("the tweet is unreachable")
	}
	if code != 200 {
		return nil, errors.Errorf("can not get %s's tweet %s", handle, id)
	}

	status := root.Get("status")
	if !status.Exists() {
		return nil, errorInvalidResponse
	}

	return toTweet(status.Raw), nil
}

func toTweet(payload string) *twitter.Tweet {
	root := gjson.Parse(payload)

	var quote *twitter.BasicTweet
	if v := root.Get(""); v.Exists() {
		quote = &twitter.BasicTweet{
			ScreenName: v.Get("author.screen_name").String(),
			Tweet:      v.Get("id").String(),
		}
	}

	var to *twitter.BasicTweet
	if v := root.Get("replying_to"); v.Exists() && v.IsObject() {
		to = &twitter.BasicTweet{
			ScreenName: v.Get("screen_name").String(),
			Tweet:      v.Get("status").String(),
		}
	}

	return &twitter.Tweet{
		ID: root.Get("id").String(),
		Author: &twitter.BasicUser{
			ScreenName: root.Get("author.screen_name").String(),
			Name:       root.Get("author.name").String(),
			Avatar:     root.Get("author.avatar_url").String(),
		},
		Content:           root.Get("text").String(),
		Attachments:       toAttachmentList(root.Get("media.all").Raw),
		Replies:           root.Get("replies").Int(),
		Retweets:          root.Get("retweets").Int(),
		Likes:             root.Get("likes").Int(),
		Views:             root.Get("views").Int(),
		ReplyingTo:        to,
		Quote:             quote,
		PossiblySensitive: root.Get("possibly_sensitive").Bool(),
		CreatedAt:         twitter.NormalizeDateTime(root.Get("created_at").String()),
	}
}

func toAttachmentList(payload string) []*twitter.Attachment {
	attachments := []*twitter.Attachment{}
	if v := gjson.Parse(payload); v.IsArray() {
		for _, item := range v.Array() {
			attachments = append(attachments, &twitter.Attachment{
				Kind:      item.Get("type").String(),
				Url:       item.Get("url").String(),
				Thumbnail: item.Get("thumbnail_url").String(),
			})
		}
	}
	return attachments
}
