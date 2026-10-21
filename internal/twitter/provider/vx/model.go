package vx

import (
	"github.com/tidwall/gjson"
	"github.com/whaoa/public-api-bridge/internal/twitter"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

func buildTweet(handle, id, payload string) (*twitter.Tweet, error) {
	root := gjson.Parse(payload)

	name := root.Get("user_name").String()
	if name == "" {
		return nil, errors.Errorf("can not get %s's tweet %s", handle, id)
	}

	var to *twitter.BasicTweet
	if v := root.Get("replyingToID").String(); v != "" {
		to = &twitter.BasicTweet{
			ScreenName: root.Get("replyingTo").String(),
			Tweet:      v,
		}
	}

	var quote *twitter.BasicTweet
	if v := root.Get("qrt.tweetID").String(); v != "" {
		quote = &twitter.BasicTweet{
			ScreenName: root.Get("qrt.user_screen_name").String(),
			Tweet:      v,
		}
	}

	return &twitter.Tweet{
		ID:                root.Get("tweetID").String(),
		Content:           root.Get("text").String(),
		Attachments:       toAttachmentList(root.Get("media_extended").Raw),
		Replies:           root.Get("replies").Int(),
		Retweets:          root.Get("retweets").Int(),
		Likes:             root.Get("likes").Int(),
		Views:             0,
		ReplyingTo:        to,
		Quote:             quote,
		CreatedAt:         twitter.NormalizeDateTime(root.Get("date").String()),
		PossiblySensitive: root.Get("possibly_sensitive").Bool(),
		Author: &twitter.BasicUser{
			ScreenName: root.Get("user_screen_name").String(),
			Name:       name,
			Avatar:     root.Get("user_profile_image_url").String(),
		},
	}, nil
}

func toAttachmentList(payload string) []*twitter.Attachment {
	attachments := []*twitter.Attachment{}
	if v := gjson.Parse(payload); v.IsArray() {
		for _, m := range v.Array() {
			kind := m.Get("type").String()
			thumbnail := m.Get("thumbnail_url").String()
			if kind == "image" {
				kind = twitter.MediaTypePhoto
				thumbnail = ""
			}
			attachments = append(attachments, &twitter.Attachment{
				Kind:      kind,
				Url:       m.Get("url").String(),
				Thumbnail: thumbnail,
			})
		}
	}
	return attachments
}
