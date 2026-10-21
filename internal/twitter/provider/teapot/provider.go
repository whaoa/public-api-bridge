package teapot

import (
	"context"
	"net/http"

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
		Get("https://nitter.cf/{handle}")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the user does not exist")
	}

	if resp.StatusCode() == http.StatusForbidden {
		return nil, errors.Forbidden("the user is suspended")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get profile for user @%s", handle)
	}

	defer resp.Body.Close()
	return buildUser(resp.Body)
}

func (p provider) ListTweets(ctx context.Context, handle string, cursor string) (*twitter.TweetSlice, error) {
	return nil, errors.New("not implemented")
}

func (p *provider) GetTweet(ctx context.Context, handle string, id string) (*twitter.Tweet, error) {
	resp, err := p.client.R().
		SetPathParam("handle", handle).
		SetPathParam("tweet", id).
		SetContext(ctx).
		Get("https://nitter.cf/{handle}/status/{tweet}")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.Forbidden("the tweet is unreachable")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get tweet %s", id)
	}

	defer resp.Body.Close()
	return buildTweet(id, resp.Body)
}

func NewProvider(client *request.Client) *provider {
	return &provider{client: client}
}
