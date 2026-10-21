package twitter

import (
	"context"

	"github.com/samber/lo"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

type Provider interface {
	GetUser(ctx context.Context, handle string) (*User, error)
	GetTweet(ctx context.Context, handle string, id string) (*Tweet, error)
	ListTweets(ctx context.Context, handle string, cursor string) (*TweetSlice, error)
}

type Service struct {
	providers []Provider
}

func (s *Service) GetUser(ctx context.Context, handle string) (*User, error) {
	for _, provider := range s.providers {
		user, err := provider.GetUser(ctx, handle)
		fillUserEntity(user)
		if err == nil || errors.IsReason(err) {
			return user, err
		}
	}
	return nil, errors.ServiceUnavailable("no twitter provider available")
}

func (s *Service) GetTweet(ctx context.Context, handle string, id string) (*Tweet, error) {
	for _, provider := range s.providers {
		tweet, err := provider.GetTweet(ctx, handle, id)
		fillTweetEntity(tweet)
		if err == nil || errors.IsReason(err) {
			return tweet, err
		}
	}
	return nil, errors.ServiceUnavailable("no twitter provider available")
}

func (s *Service) ListTweets(ctx context.Context, handle string, cursor string) (*TweetSlice, error) {
	for _, provider := range s.providers {
		result, err := provider.ListTweets(ctx, handle, cursor)
		if result != nil && len(result.Items) > 0 {
			lo.ForEach(result.Items, func(tweet *Tweet, _ int) {
				fillTweetEntity(tweet)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
	}
	return nil, errors.ServiceUnavailable("no twitter provider available")
}

func NewService(logger *slog.Logger, providers ...Provider) *Service {
	return &Service{
		logger:    logger,
		providers: providers,
	}
}
