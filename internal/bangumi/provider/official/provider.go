package official

import (
	"context"
	"net/http"
	"strconv"

	"github.com/whaoa/public-api-bridge/internal/bangumi"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/request"
)

type provider struct {
	client *request.Client
}

func (p *provider) Name() string {
	return "bangumi-official"
}

func (p *provider) SearchAnime(ctx context.Context, keyword string, offset, count int) (*bangumi.SubjectSlice, error) {
	if offset < 0 {
		offset = 0
	}
	if count <= 0 {
		count = 10
	}
	resp, err := p.client.R().
		SetContext(ctx).
		SetBody(map[string]any{
			"keyword": keyword,
			"sort":    "match",
			"filter":  map[string]any{"type": []int{2}},
			"limit":   count,
			"offset":  offset,
		}).
		Post("https://api.bgm.tv/v0/search/subjects")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not search for %s", keyword)
	}

	return buildSubjectSlice(resp.String())
}

func (p *provider) GetAnime(ctx context.Context, id string) (*bangumi.Subject, error) {
	resp, err := p.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		Get("https://api.bgm.tv/v0/subjects/{id}")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusBadRequest {
		return nil, normalizeBadRequestError(resp.String())
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get subject %s", id)
	}

	return buildSubject(resp.String())
}

func (p *provider) ListEpisodes(ctx context.Context, id string, offset, count int) (*bangumi.EpisodeSlice, error) {
	if offset < 0 {
		offset = 0
	}
	if count <= 0 {
		count = 10
	}
	resp, err := p.client.R().
		SetContext(ctx).
		SetQueryParam("subject_id", id).
		SetQueryParam("offset", strconv.Itoa(offset)).
		SetQueryParam("limit", strconv.Itoa(count)).
		Get("https://api.bgm.tv/v0/episodes")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the subject is not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get episodes for %s", id)
	}

	return buildEpisodeSlice(resp.String())
}

func (p *provider) GetCharacters(ctx context.Context, id string) ([]*bangumi.Character, error) {
	resp, err := p.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		Get("https://api.bgm.tv/v0/subjects/{id}/characters")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusBadRequest {
		return nil, normalizeBadRequestError(resp.String())
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get characters for %s", id)
	}

	return buildCharacterList(resp.String())
}

func (p *provider) GetStaffs(ctx context.Context, id string) ([]*bangumi.Staff, error) {
	resp, err := p.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		Get("https://api.bgm.tv/v0/subjects/{id}/persons")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the subject is not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get staffs for %s", id)
	}

	return buildStaffList(resp.String())
}

func (p *provider) GetPerson(ctx context.Context, id string) (*bangumi.Person, error) {
	resp, err := p.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		Get("https://api.bgm.tv/v0/persons/{id}")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the person is not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get person %s", id)
	}

	return buildPerson(resp.String())
}

func (p *provider) GetPersonCredits(ctx context.Context, id string) ([]*bangumi.PersonCredit, error) {
	resp, err := p.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		Get("https://api.bgm.tv/v0/persons/{id}/characters")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the person is not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get credits for %s", id)
	}

	return buildPersonCreditList(resp.String())
}

func NewProvider(client *request.Client) *provider {
	return &provider{client: client}
}
