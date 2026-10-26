package tmdb

import (
	"context"
	"log/slog"

	"github.com/samber/lo"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/log"
)

type Provider interface {
	Name() string
	SearchMedia(ctx context.Context, keyword, lang string, page int) (*MediaSlice, error)
	GetMovie(ctx context.Context, id string, lang string) (*Movie, error)
	GetTV(ctx context.Context, id string, lang string) (*TV, error)
	GetSeason(ctx context.Context, id, season, lang string) (*Season, error)
	GetMediaCasts(ctx context.Context, kind, id string, lang string) (*MediaCasts, error)
	GetMediaImages(ctx context.Context, kind, id string, lang string) (*MediaImages, error)
	SearchPerson(ctx context.Context, keyword, lang string, page int) (*EmbedPersonSlice, error)
	GetPerson(ctx context.Context, id, lang string) (*Person, error)
	GetPersonProfiles(ctx context.Context, id, lang string) ([]*Image, error)
	GetPersonCredits(ctx context.Context, id, lang string) (*PersonCredits, error)
}

type Service struct {
	logger    *slog.Logger
	providers []Provider
}

func (s *Service) SearchMedia(ctx context.Context, keyword, lang string, page int) (*MediaSlice, error) {
	for _, provider := range s.providers {
		result, err := provider.SearchMedia(ctx, keyword, lang, page)
		if result != nil && len(result.Items) > 0 {
			lo.ForEach(result.Items, func(item *Media, _ int) {
				item.Poster = buildImageUrl(item.Poster, PosterSizeSmall)
				item.Url = buildRefUrl(item.Kind, item.ID)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.SearchMedia failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func (s *Service) GetMovie(ctx context.Context, id, lang string) (*Movie, error) {
	for _, provider := range s.providers {
		result, err := provider.GetMovie(ctx, id, lang)
		if result != nil {
			result.Backdrop = buildImageUrl(result.Backdrop, BackdropSizeLarge)
			result.Poster = buildImageUrl(result.Poster, PosterSizeLarge)
			lo.ForEach(result.ProductionCompanies, func(item *Company, _ int) {
				item.Logo = buildImageUrl(item.Logo, LogoSizeSmall)
				item.Url = buildRefUrl(urlTypeCompany, item.ID)
			})
			result.Url = buildRefUrl(TypeMovie, result.ID)
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.GetMovie failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func (s *Service) GetTV(ctx context.Context, id, lang string) (*TV, error) {
	for _, provider := range s.providers {
		result, err := provider.GetTV(ctx, id, lang)
		if result != nil {
			result.Backdrop = buildImageUrl(result.Backdrop, BackdropSizeLarge)
			result.Poster = buildImageUrl(result.Poster, PosterSizeLarge)
			lo.ForEach(result.ProductionCompanies, func(item *Company, _ int) {
				item.Logo = buildImageUrl(item.Logo, LogoSizeSmall)
				item.Url = buildRefUrl(urlTypeCompany, item.ID)
			})
			lo.ForEach(result.Seasons, func(item *EmbedSeason, _ int) {
				item.Poster = buildImageUrl(item.Poster, PosterSizeMedium)
			})
			result.Url = buildRefUrl(TypeTV, result.ID)
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.GetTV failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func (s *Service) GetSeason(ctx context.Context, id, season, lang string) (*Season, error) {
	for _, provider := range s.providers {
		result, err := provider.GetSeason(ctx, id, season, lang)
		if result != nil {
			lo.ForEach(result.Episodes, func(item *Episode, _ int) {
				item.Still = buildImageUrl(item.Still, StillSizeLarge)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.GetSeason failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func (s *Service) GetMediaCasts(ctx context.Context, kind, id, lang string) (*MediaCasts, error) {
	for _, provider := range s.providers {
		result, err := provider.GetMediaCasts(ctx, kind, id, lang)
		if result != nil {
			lo.ForEach(result.Actors, func(item *MediaCast, _ int) {
				item.Profile = buildImageUrl(item.Profile, ProfileSizeMedium)
				item.Url = buildRefUrl(urlTypePerson, item.ID)
			})
			lo.ForEach(result.Staffs, func(item *MediaCast, _ int) {
				item.Profile = buildImageUrl(item.Profile, ProfileSizeMedium)
				item.Url = buildRefUrl(urlTypePerson, item.ID)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.GetMediaCasts failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func (s *Service) GetMediaImages(ctx context.Context, kind, id, lang string) (*MediaImages, error) {
	for _, provider := range s.providers {
		result, err := provider.GetMediaImages(ctx, kind, id, lang)
		if result != nil {
			lo.ForEach(result.Backdrops, func(item *Image, _ int) {
				item.Url = buildImageUrl(item.Url, ImageSizeRaw)
			})
			lo.ForEach(result.Logos, func(item *Image, _ int) {
				item.Url = buildImageUrl(item.Url, ImageSizeRaw)
			})
			lo.ForEach(result.Posters, func(item *Image, _ int) {
				item.Url = buildImageUrl(item.Url, ImageSizeRaw)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.GetMediaImages failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func (s *Service) SearchPerson(ctx context.Context, keyword, lang string, page int) (*EmbedPersonSlice, error) {
	for _, provider := range s.providers {
		result, err := provider.SearchPerson(ctx, keyword, lang, page)
		if result != nil {
			lo.ForEach(result.Items, func(item *EmbedPerson, _ int) {
				item.Profile = buildImageUrl(item.Profile, ProfileSizeMedium)
				item.Url = buildRefUrl(urlTypePerson, item.ID)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.SearchPerson failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func (s *Service) GetPerson(ctx context.Context, id, lang string) (*Person, error) {
	for _, provider := range s.providers {
		result, err := provider.GetPerson(ctx, id, lang)
		if result != nil {
			result.Profile = buildImageUrl(result.Profile, PosterSizeMedium)
			result.Url = buildRefUrl(urlTypePerson, result.ID)
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.GetPerson failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func (s *Service) GetPersonProfiles(ctx context.Context, id, lang string) ([]*Image, error) {
	for _, provider := range s.providers {
		result, err := provider.GetPersonProfiles(ctx, id, lang)
		if result != nil {
			lo.ForEach(result, func(item *Image, _ int) {
				item.Url = buildImageUrl(item.Url, ImageSizeRaw)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.GetPersonProfiles failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func (s *Service) GetPersonCredits(ctx context.Context, id, lang string) (*PersonCredits, error) {
	for _, provider := range s.providers {
		result, err := provider.GetPersonCredits(ctx, id, lang)
		if result != nil {
			lo.ForEach(result.AsActor, func(item *MediaCredit, _ int) {
				item.Media.Poster = buildImageUrl(item.Media.Poster, PosterSizeMedium)
				item.Url = buildRefUrl(item.Kind, item.ID)
			})
			lo.ForEach(result.AsStaff, func(item *MediaCredit, _ int) {
				item.Media.Poster = buildImageUrl(item.Media.Poster, PosterSizeMedium)
				item.Url = buildRefUrl(item.Kind, item.ID)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
		s.logger.ErrorContext(
			ctx,
			"provider.GetPersonCredits failed",
			slog.String("provider", provider.Name()),
			log.Err(err),
		)
	}
	return nil, errors.Internal("no tmdb provider available")
}

func NewService(logger *slog.Logger, providers ...Provider) *Service {
	return &Service{
		logger:    logger,
		providers: providers,
	}
}
