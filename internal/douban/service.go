package douban

import (
	"context"
	"log/slog"

	"github.com/samber/lo"
	"github.com/whaoa/public-api-bridge/pkg/log"
)

type Provider interface {
	Name() string

	SearchMovie(ctx context.Context, keyword string) ([]*Subject, error)

	GetMovie(ctx context.Context, kind, id string) (*Movie, error)
	ListMoviePhotos(ctx context.Context, kind, id string, offset, count int) (*PhotoSlice, error)
	GetMovieVideos(ctx context.Context, kind, id string) ([]*Video, error)
	GetMovieCasts(ctx context.Context, kind, id string) (*MovieCasts, error)
	GetMovieRecommendations(ctx context.Context, kind, id string) ([]*MovieRecommendation, error)

	GetPerson(ctx context.Context, id string) (*Person, error)
	ListPersonPhotos(ctx context.Context, id string, offset, count int) (*PhotoSlice, error)
	ListPersonMovieCredits(ctx context.Context, id string, offset, count int) (*MovieCreditSlice, error)

	GetVideoSource(ctx context.Context, id string) (string, error)
}

type Service struct {
	logger   *slog.Logger
	provider Provider
}

func (s *Service) SearchMovie(ctx context.Context, keyword string) ([]*Subject, error) {
	result, err := s.provider.SearchMovie(ctx, keyword)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.SearchMovie failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	if result != nil && len(result) > 0 {
		lo.ForEach(result, func(item *Subject, _ int) {
			item.Poster = withImageProxy(withImageResize(item.Poster, ImageSizePoster))
		})
	}
	return result, err
}

func (s *Service) GetMovie(ctx context.Context, kind, id string) (*Movie, error) {
	movie, err := s.provider.GetMovie(ctx, kind, id)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.GetMovie failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	if movie != nil {
		movie.Poster = withImageProxy(withImageResize(movie.Poster, ImageSizeLarge))
		movie.Thumbnail = withImageResize(movie.Poster, ImageSizePoster)
	}
	return movie, err
}

func (s *Service) ListMoviePhotos(ctx context.Context, kind, id string, offset, count int) (*PhotoSlice, error) {
	result, err := s.provider.ListMoviePhotos(ctx, kind, id, offset, count)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.ListMoviePhotos failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	if result != nil && len(result.Items) > 0 {
		lo.ForEach(result.Items, func(item *Photo, _ int) {
			item.Large = withImageProxy(item.Large)
			item.Normal = withImageProxy(item.Normal)
			item.Small = withImageProxy(item.Small)
		})
	}
	return result, err
}

func (s *Service) GetMovieVideos(ctx context.Context, kind, id string) ([]*Video, error) {
	trailers, err := s.provider.GetMovieVideos(ctx, kind, id)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.GetMovieVideos failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	if len(trailers) > 0 {
		lo.ForEach(trailers, func(item *Video, _ int) {
			item.Thumbnail = withImageProxy(item.Thumbnail)
		})
	}
	return trailers, err
}

func (s *Service) GetMovieCasts(ctx context.Context, kind, id string) (*MovieCasts, error) {
	casts, err := s.provider.GetMovieCasts(ctx, kind, id)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.GetMovieCasts failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	if casts != nil {
		lo.ForEach(casts.Directors, func(item *MovieCast, _ int) {
			item.Avatar = withImageProxy(item.Avatar)
		})
		lo.ForEach(casts.Actors, func(item *MovieCast, _ int) {
			item.Avatar = withImageProxy(item.Avatar)
		})
	}
	return casts, err
}

func (s *Service) GetMovieRecommendations(ctx context.Context, kind, id string) ([]*MovieRecommendation, error) {
	recommendations, err := s.provider.GetMovieRecommendations(ctx, kind, id)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.GetMovieRecommendations failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	if recommendations != nil && len(recommendations) > 0 {
		lo.ForEach(recommendations, func(item *MovieRecommendation, _ int) {
			item.Poster = withImageProxy(item.Poster)
		})
	}
	return recommendations, err
}

func (s *Service) GetPerson(ctx context.Context, id string) (*Person, error) {
	person, err := s.provider.GetPerson(ctx, id)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.GetPerson failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	if person != nil {
		person.Photo = withImageProxy(person.Photo)
	}
	return person, err
}

func (s *Service) ListPersonPhotos(ctx context.Context, id string, offset, count int) (*PhotoSlice, error) {
	photos, err := s.provider.ListPersonPhotos(ctx, id, offset, count)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.ListPersonPhotos failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	if photos != nil && len(photos.Items) > 0 {
		lo.ForEach(photos.Items, func(item *Photo, _ int) {
			item.Large = withImageProxy(item.Large)
			item.Normal = withImageProxy(item.Normal)
			item.Small = withImageProxy(item.Small)
		})
	}
	return photos, err
}

func (s *Service) ListPersonMovieCredits(ctx context.Context, id string, offset, count int) (*MovieCreditSlice, error) {
	result, err := s.provider.ListPersonMovieCredits(ctx, id, offset, count)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.ListPersonMovieCredits failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	if result != nil && len(result.Items) > 0 {
		lo.ForEach(result.Items, func(item *MovieCredit, _ int) {
			item.Poster = withImageProxy(withImageResize(item.Poster, ImageSizeLarge))
			item.Thumbnail = withImageResize(item.Poster, ImageSizePoster)
		})
	}
	return result, err
}

func (s *Service) GetVideoSource(ctx context.Context, id string) (string, error) {
	result, err := s.provider.GetVideoSource(ctx, id)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"provider.GetVideoSource failed",
			slog.String("provider", s.provider.Name()),
			log.Err(err),
		)
	}
	return result, err
}

func NewService(logger *slog.Logger, provider Provider) *Service {
	return &Service{
		logger:   logger,
		provider: provider,
	}
}
