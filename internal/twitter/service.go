package twitter

import (
	"context"
	"log/slog"

	"github.com/samber/lo"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/log"
)

type Provider interface {
	Name() string
	GetUser(ctx context.Context, handle string) (*User, error)
	GetTweet(ctx context.Context, handle string, id string) (*Tweet, error)
	ListTweets(ctx context.Context, handle string, cursor string) (*TweetSlice, error)
}

type Service struct {
	logger    *slog.Logger
	providers []Provider
}

func (s *Service) GetUser(ctx context.Context, handle string) (*User, error) {
	for _, provider := range s.providers {
		user, err := provider.GetUser(ctx, handle)
		fillUserEntity(user)
		if err == nil || errors.IsReason(err) {
			return user, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.GetUser failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
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
		s.logger.ErrorContext(
			ctx,
			"provider.GetTweet failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
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
		s.logger.ErrorContext(
			ctx,
			"provider.ListTweets failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.ServiceUnavailable("no twitter provider available")
}

func NewService(logger *slog.Logger, providers ...Provider) *Service {
	return &Service{
		logger:    logger,
		providers: providers,
	}
}
