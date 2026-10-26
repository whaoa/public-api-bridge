package bangumi

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

type Provider interface {
	SearchAnime(ctx context.Context, keyword string, offset, count int) (*SubjectSlice, error)
	GetAnime(ctx context.Context, id string) (*Subject, error)
	ListEpisodes(ctx context.Context, id string, offset, count int) (*EpisodeSlice, error)
	GetCharacters(ctx context.Context, id string) ([]*Character, error)
	GetStaffs(ctx context.Context, id string) ([]*Staff, error)
	GetPerson(ctx context.Context, id string) (*Person, error)
	GetPersonCredits(ctx context.Context, id string) ([]*PersonCredit, error)
}

type Service struct {
	providers []Provider
}

func (s *Service) SearchAnime(ctx context.Context, keyword string, offset, count int) (*SubjectSlice, error) {
	for _, provider := range s.providers {
		result, err := provider.SearchAnime(ctx, keyword, offset, count)
		if result != nil {
			lo.ForEach(result.Items, func(item *Subject, _ int) {
				item.Url = fmt.Sprintf("https://bangumi.tv/subject/%d", item.ID)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
	}
	return nil, errors.ServiceUnavailable("no bangumi provider available")
}

func (s *Service) GetAnime(ctx context.Context, id string) (*Subject, error) {
	for _, provider := range s.providers {
		result, err := provider.GetAnime(ctx, id)
		if result != nil {
			result.Url = fmt.Sprintf("https://bangumi.tv/subject/%d", result.ID)
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
	}
	return nil, errors.ServiceUnavailable("no bangumi provider available")
}

func (s *Service) ListEpisodes(ctx context.Context, id string, offset, count int) (*EpisodeSlice, error) {
	for _, provider := range s.providers {
		result, err := provider.ListEpisodes(ctx, id, offset, count)
		if result != nil && len(result.Items) > 0 {
			lo.ForEach(result.Items, func(item *Episode, _ int) {
				item.Url = fmt.Sprintf("https://bangumi.tv/ep/%d", item.ID)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
	}
	return nil, errors.ServiceUnavailable("no bangumi provider available")
}

func (s *Service) GetCharacters(ctx context.Context, id string) ([]*Character, error) {
	for _, provider := range s.providers {
		result, err := provider.GetCharacters(ctx, id)
		if result != nil && len(result) > 0 {
			lo.ForEach(result, func(item *Character, _ int) {
				item.Url = fmt.Sprintf("https://bangumi.tv/character/%d", item.ID)
				lo.ForEach(item.Actors, func(actor *EmbedPerson, _ int) {
					actor.Url = fmt.Sprintf("https://bangumi.tv/person/%d", actor.ID)
				})
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
	}
	return nil, errors.ServiceUnavailable("no bangumi provider available")
}

func (s *Service) GetStaffs(ctx context.Context, id string) ([]*Staff, error) {
	for _, provider := range s.providers {
		result, err := provider.GetStaffs(ctx, id)
		if result != nil && len(result) > 0 {
			lo.ForEach(result, func(item *Staff, _ int) {
				item.Url = fmt.Sprintf("https://bangumi.tv/person/%d", item.ID)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
	}
	return nil, errors.ServiceUnavailable("no bangumi provider available")
}

func (s *Service) GetPerson(ctx context.Context, id string) (*Person, error) {
	for _, provider := range s.providers {
		result, err := provider.GetPerson(ctx, id)
		if result != nil {
			result.Url = fmt.Sprintf("https://bangumi.tv/person/%d", result.ID)
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
	}
	return nil, errors.ServiceUnavailable("no bangumi provider available")
}

func (s *Service) GetPersonCredits(ctx context.Context, id string) ([]*PersonCredit, error) {
	for _, provider := range s.providers {
		result, err := provider.GetPersonCredits(ctx, id)
		if result != nil && len(result) > 0 {
			lo.ForEach(result, func(item *PersonCredit, _ int) {
				item.Url = fmt.Sprintf("https://bangumi.tv/subject/%d", item.ID)
			})
		}
		if err == nil || errors.IsReason(err) {
			return result, err
		}
	}
	return nil, errors.ServiceUnavailable("no bangumi provider available")
}

func NewService(providers ...Provider) *Service {
	return &Service{providers: providers}
}
