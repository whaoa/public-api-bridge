package fx

import (
	"context"
	"time"

	"github.com/whaoa/public-api-bridge/internal/twitter"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/request"
)

type provider struct {
	client *request.Client
}

func (p *provider) GetUser(ctx context.Context, handle string) (*twitter.User, error) {
	resp, err := p.client.R().
		SetPathParam("handle", handle).
		SetContext(ctx).
		Get("https://api.fxtwitter.com/2/profile/{handle}")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buildUser(handle, resp.String())
}

func (p *provider) ListTweets(ctx context.Context, handle string, cursor string) (*twitter.TweetSlice, error) {
	resp, err := p.client.R().
		SetPathParam("handle", handle).
		SetQueryParam("count", "20").
		SetQueryParam("cursor", cursor).
		SetTimeout(time.Second * 10).
		SetContext(ctx).
		Get("https://api.fxtwitter.com/2/profile/{handle}/statuses")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buildTweetSlice(handle, resp.String())
}

func (p *provider) GetTweet(ctx context.Context, handle string, id string) (*twitter.Tweet, error) {
	resp, err := p.client.R().
		SetPathParam("tweet", id).
		SetContext(ctx).
		Get("https://api.fxtwitter.com/2/status/{tweet}")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buildTweet(handle, id, resp.String())
}

func NewProvider(client *request.Client) *provider {
	return &provider{client: client}
}
