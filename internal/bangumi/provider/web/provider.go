package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/whaoa/public-api-bridge/internal/bangumi"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/request"
)

type provider struct {
	client *request.Client
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
		Post("https://next.bgm.tv/p1/search/subjects")

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
		Get("https://next.bgm.tv/p1/subjects/{id}/home")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the subject is not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get subject %s", id)
	}

	return buildSubject(gjson.Get(resp.String(), "subject").Raw)
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
		SetPathParam("id", id).
		SetQueryParam("offset", strconv.Itoa(offset)).
		SetQueryParam("limit", strconv.Itoa(count)).
		Get("https://next.bgm.tv/p1/subjects/{id}/episodes")

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

func (p *provider) listCharacters(ctx context.Context, id string, offset, count int) (int, []*bangumi.Character, error) {
	if offset < 0 {
		offset = 0
	}
	if count <= 0 {
		count = 10
	}
	if count <= 100 {
		count = 100
	}
	resp, err := p.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		SetQueryParam("offset", strconv.Itoa(offset)).
		SetQueryParam("limit", strconv.Itoa(count)).
		Get("https://next.bgm.tv/p1/subjects/{id}/characters")

	if err != nil {
		return 0, nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return 0, nil, errors.NotFound("the subject is not exist")
	}

	if resp.IsStatusFailure() {
		return 0, nil, errors.Errorf("can not get characters for %s", id)
	}

	payload := resp.String()
	items, err := buildCharacterList(payload)
	if err != nil {
		return 0, nil, err
	}
	length := offset + len(items)

	if total := int(gjson.Get(payload, "total").Int()); total > length {
		_, result, err := p.listCharacters(ctx, id, length, 100)
		if err != nil {
			return 0, nil, err
		}
		items = append(items, result...)
	}

	return len(items), items, nil
}

func (p *provider) GetCharacters(ctx context.Context, id string) ([]*bangumi.Character, error) {
	_, result, err := p.listCharacters(ctx, id, 0, 100)
	return result, err
}

func (p *provider) listStaffs(ctx context.Context, id string, offset, count int) (int, []*bangumi.Staff, error) {
	resp, err := p.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		SetQueryParam("offset", strconv.Itoa(offset)).
		SetQueryParam("limit", strconv.Itoa(count)).
		Get("https://next.bgm.tv/p1/subjects/{id}/staffs/persons")

	if err != nil {
		return 0, nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return 0, nil, errors.NotFound("the subject is not exist")
	}

	if resp.IsStatusFailure() {
		return 0, nil, errors.Errorf("can not get staffs for %s", id)
	}

	payload := resp.String()
	items, err := buildStaffList(payload)
	if err != nil {
		return 0, nil, err
	}
	length := offset + len(items)

	if total := int(gjson.Get(payload, "total").Int()); total > length {
		_, result, err := p.listStaffs(ctx, id, length, 100)
		if err != nil {
			return 0, nil, err
		}
		items = append(items, result...)
	}

	return len(items), items, nil
}

func (p *provider) GetStaffs(ctx context.Context, id string) ([]*bangumi.Staff, error) {
	_, result, err := p.listStaffs(ctx, id, 0, 100)
	return result, err
}

func (p *provider) GetPerson(ctx context.Context, id string) (*bangumi.Person, error) {
	resp, err := p.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		Get("https://next.bgm.tv/p1/persons/{id}")

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

func (p *provider) listPersonCredits(ctx context.Context, id string, offset, count int) (int, []*bangumi.PersonCredit, error) {
	resp, err := p.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		SetQueryParam("offset", strconv.Itoa(offset)).
		SetQueryParam("limit", strconv.Itoa(count)).
		Get("https://next.bgm.tv/p1/persons/{id}/casts")

	if err != nil {
		return 0, nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return 0, nil, errors.NotFound("the person is not exist")
	}

	if resp.IsStatusFailure() {
		return 0, nil, errors.Errorf("can not get credits for %s", id)
	}

	payload := resp.String()
	items, err := buildPersonCreditList(payload)
	if err != nil {
		return 0, nil, err
	}
	length := offset + len(items)

	if total := int(gjson.Get(payload, "total").Int()); total > length {
		_, result, err := p.listPersonCredits(ctx, id, length, 100)
		if err != nil {
			return 0, nil, err
		}
		items = append(items, result...)
	}

	return len(items), items, nil
}

func (p *provider) GetPersonCredits(ctx context.Context, id string) ([]*bangumi.PersonCredit, error) {
	_, result, err := p.listPersonCredits(ctx, id, 0, 100)
	return result, err
}

func NewProvider(client *request.Client) *provider {
	return &provider{client: client}
}
