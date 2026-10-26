package vx

import (
	"context"

	"github.com/whaoa/public-api-bridge/internal/twitter"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/request"
)

type provider struct {
	client *request.Client
}

func (p *provider) Name() string {
	return "twitter-vx"
}

func (p *provider) GetUser(ctx context.Context, handle string) (*twitter.User, error) {
	return nil, errors.New("not implemented")
}

func (p provider) ListTweets(ctx context.Context, handle string, cursor string) (*twitter.TweetSlice, error) {
	return nil, errors.New("not implemented")
}

func (p *provider) GetTweet(ctx context.Context, handle string, id string) (*twitter.Tweet, error) {
	resp, err := p.client.R().
		SetHeader("User-Agent", request.UserAgentTelegramBot).
		SetPathParam("handle", handle).
		SetPathParam("tweet", id).
		SetContext(ctx).
		Get("https://api.vxtwitter.com/{handle}/status/{tweet}")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.Header().Get("Content-Type") != "application/json" {
		return nil, errors.NotFound("the tweet is unreachable")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get tweet %s", id)
	}

	return buildTweet(handle, id, resp.String())
}

func NewProvider(client *request.Client) twitter.Provider {
	return &provider{client: client}
}
