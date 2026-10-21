package official

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/whaoa/public-api-bridge/internal/tmdb"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/request"
)

var (
	errorInvalidMediaType = errors.BadRequest("invalid media type")
)

type provider struct {
	key    string
	client *request.Client
}

func (p *provider) SearchMedia(ctx context.Context, keyword, lang string, page int) (*tmdb.MediaSlice, error) {
	if page <= 0 {
		page = 1
	}

	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("query", keyword).
		SetQueryParam("language", lang).
		SetQueryParam("page", strconv.Itoa(page)).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/search/multi")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buildMediaSlice(resp.String())
}

func (p *provider) GetMovie(ctx context.Context, id, lang string) (*tmdb.Movie, error) {
	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("language", lang).
		SetPathParam("id", id).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/movie/{id}")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the movie does not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get movie %s", id)
	}

	return buildMovie(resp.String())
}

func (p *provider) GetTV(ctx context.Context, id, lang string) (*tmdb.TV, error) {
	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("language", lang).
		SetQueryParam("append_to_response", "external_ids").
		SetPathParam("id", id).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/tv/{id}")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the tv does not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get tv %s", id)
	}

	return buildTV(resp.String())
}

func (p *provider) GetSeason(ctx context.Context, id, season, lang string) (*tmdb.Season, error) {
	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("language", lang).
		SetPathParam("id", id).
		SetPathParam("season", season).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/tv/{id}/season/{season}")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the season does not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get season %s", season)
	}

	return buildSeason(resp.String())
}

func (p *provider) GetMediaCasts(ctx context.Context, kind, id, lang string) (*tmdb.MediaCasts, error) {
	if !tmdb.ValidateMediaType(kind) {
		return nil, errorInvalidMediaType
	}

	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("language", lang).
		SetPathParam("type", kind).
		SetPathParam("id", id).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/{type}/{id}/credits")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound(fmt.Sprintf("the media %s not exist", kind))
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get casts for %s %s", kind, id)
	}

	return buildMediaCasts(resp.String())
}

func (p *provider) GetMediaImages(ctx context.Context, kind, id, lang string) (*tmdb.MediaImages, error) {
	if !tmdb.ValidateMediaType(kind) {
		return nil, errorInvalidMediaType
	}

	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("language", lang).
		SetPathParam("type", kind).
		SetPathParam("id", id).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/{type}/{id}/images")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound(fmt.Sprintf("the media %s not exist", kind))
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get images for %s %s", kind, id)
	}

	return buildMediaImages(resp.String())
}

func (p *provider) SearchPerson(ctx context.Context, keyword, lang string, page int) (*tmdb.EmbedPersonSlice, error) {
	if page <= 0 {
		page = 1
	}

	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("language", lang).
		SetQueryParam("query", keyword).
		SetQueryParam("page", strconv.Itoa(page)).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/search/person")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buildEmbedPersonSlice(resp.String())
}

func (p *provider) GetPerson(ctx context.Context, id, lang string) (*tmdb.Person, error) {
	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("language", lang).
		SetPathParam("id", id).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/person/{id}")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the person does not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get person %s", id)
	}

	return buildPerson(resp.String())
}

func (p *provider) GetPersonProfiles(ctx context.Context, id, lang string) ([]*tmdb.Image, error) {
	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("language", lang).
		SetPathParam("id", id).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/person/{id}/images")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the person does not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get images for person %s", id)
	}

	return buildPersonProfiles(resp.String())
}

func (p *provider) GetPersonCredits(ctx context.Context, id, lang string) (*tmdb.PersonCredits, error) {
	resp, err := p.client.R().
		SetQueryParam("api_key", p.key).
		SetQueryParam("language", lang).
		SetPathParam("id", id).
		SetContext(ctx).
		Get("https://api.themoviedb.org/3/person/{id}/combined_credits")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound("the person does not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Errorf("can not get credits for person %s", id)
	}

	return buildPersonCredits(resp.String())
}

func NewProvider(client *request.Client, key string) *provider {
	return &provider{client: client, key: key}
}
