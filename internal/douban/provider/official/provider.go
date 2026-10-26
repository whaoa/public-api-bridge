package official

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/whaoa/public-api-bridge/internal/douban"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/request"
)

const (
	microMessengerKey     = "0ac44ae016490db2204ce0a042db2916"
	microMessengerReferer = "https://servicewechat.com/wx2f9b06c1de1ccfca/99/page-frame.html"
)

type provider struct {
	client *request.Client
}

func (p *provider) Name() string {
	return "douban-official"
}

func (p *provider) SearchMovie(ctx context.Context, keyword string) ([]*douban.Subject, error) {
	resp, err := p.client.R().
		SetHeader("Referer", "https://m.douban.com/search/").
		SetQueryParam("type", douban.SubjectTypeMovie).
		SetQueryParam("q", keyword).
		SetContext(ctx).
		Get(withProxy("https://m.douban.com/rexxar/api/v2/search/subjects"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.IsStatusFailure() {
		return nil, errors.Internal("unexpected upstream response")
	}

	return buildSubjectList(resp.String())
}

func (p *provider) GetMovie(ctx context.Context, kind, id string) (*douban.Movie, error) {
	resp, err := p.client.R().
		SetHeader("Referer", fmt.Sprintf("https://m.douban.com/%s/subject/%s/", kind, id)).
		SetPathParam("type", kind).
		SetPathParam("id", id).
		SetQueryParam("ck", "").
		SetQueryParam("for_mobile", "1").
		SetContext(ctx).
		Get(withProxy("https://m.douban.com/rexxar/api/v2/{type}/{id}"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() >= 400 && resp.StatusCode() < 500 {
		return nil, errors.NotFound("the movie does not exist")
	}

	if resp.IsStatusFailure() {
		return nil, errors.Internal(fmt.Sprintf("can not get %s detail for %s", kind, id))
	}

	return buildMovie(resp.String())
}

func (p *provider) ListMoviePhotos(ctx context.Context, kind, id string, offset, count int) (*douban.PhotoSlice, error) {
	if count == 0 {
		count = 20
	}

	resp, err := p.client.R().
		SetHeader("Referer", fmt.Sprintf("https://m.douban.com/%s/subject/%s/", kind, id)).
		SetPathParam("type", kind).
		SetPathParam("id", id).
		SetQueryParam("start", strconv.Itoa(offset)).
		SetQueryParam("count", strconv.Itoa(count)).
		SetContext(ctx).
		Get(withProxy("https://m.douban.com/rexxar/api/v2/{type}/{id}/photos"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound(fmt.Sprintf("the %s does not exist", kind))
	}

	if resp.IsStatusFailure() {
		return nil, errors.Internal(fmt.Sprintf("can not get %s photos for %s", kind, id))
	}

	return buildPhotoSlice(resp.String())
}

func (p *provider) GetMovieVideos(ctx context.Context, kind, id string) ([]*douban.Video, error) {
	resp, err := p.client.R().
		SetHeader("Referer", fmt.Sprintf("https://m.douban.com/%s/subject/%s/", kind, id)).
		SetPathParam("type", kind).
		SetPathParam("id", id).
		SetContext(ctx).
		Get(withProxy("https://m.douban.com/rexxar/api/v2/{type}/{id}/trailers"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound(fmt.Sprintf("the %s does not exist", kind))
	}

	if resp.IsStatusFailure() {
		return nil, errors.Internal(fmt.Sprintf("can not get %s trailers for %s", kind, id))
	}

	return buildVideoList(resp.String())
}

func (p *provider) GetMovieCasts(ctx context.Context, kind, id string) (*douban.MovieCasts, error) {
	resp, err := p.client.R().
		SetHeader("Referer", fmt.Sprintf("https://m.douban.com/%s/subject/%s/", kind, id)).
		SetPathParam("type", kind).
		SetPathParam("id", id).
		SetContext(ctx).
		Get(withProxy("https://m.douban.com/rexxar/api/v2/{type}/{id}/celebrities"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound(fmt.Sprintf("the %s does not exist", kind))
	}

	if resp.IsStatusFailure() {
		return nil, errors.Internal(fmt.Sprintf("can not get %s casts for %s", kind, id))
	}

	return buildMovieCasts(resp.String())
}

func (p *provider) GetMovieRecommendations(ctx context.Context, kind, id string) ([]*douban.MovieRecommendation, error) {
	resp, err := p.client.R().
		SetHeader("Referer", fmt.Sprintf("https://m.douban.com/%s/subject/%s/", kind, id)).
		SetPathParam("type", kind).
		SetPathParam("id", id).
		SetContext(ctx).
		Get(withProxy("https://m.douban.com/rexxar/api/v2/{type}/{id}/recommendations"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound(fmt.Sprintf("the %s does not exist", kind))
	}

	if resp.IsStatusFailure() {
		return nil, errors.Internal(fmt.Sprintf("can not get %s recommendations for %s", kind, id))
	}

	return buildMovieRecommendations(resp.String())
}

func (p *provider) GetPerson(ctx context.Context, id string) (*douban.Person, error) {
	resp, err := p.client.R().
		SetHeader("Referer", fmt.Sprintf("https://m.douban.com/movie/celebrity/%s/", id)).
		SetPathParam("id", id).
		SetContext(ctx).
		Get(withProxy("https://m.douban.com/rexxar/api/v2/celebrity/{id}"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound(fmt.Sprintf("the person does not exist"))
	}

	if resp.IsStatusFailure() {
		return nil, errors.Internal(fmt.Sprintf("can not get person for %s", id))
	}

	return buildPerson(id, resp.String())
}

func (p *provider) ListPersonPhotos(ctx context.Context, id string, offset, count int) (*douban.PhotoSlice, error) {
	resp, err := p.client.R().
		SetHeader("Referer", fmt.Sprintf("https://m.douban.com/movie/celebrity/%s/", id)).
		SetPathParam("id", id).
		SetContext(ctx).
		Get(withProxy("https://m.douban.com/rexxar/api/v2/celebrity/{id}/photos"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound(fmt.Sprintf("the person does not exist"))
	}

	if resp.IsStatusFailure() {
		return nil, errors.Internal(fmt.Sprintf("can not get photos for %s", id))
	}

	return buildPhotoSlice(resp.String())
}

func (p *provider) ListPersonMovieCredits(ctx context.Context, id string, offset, count int) (*douban.MovieCreditSlice, error) {
	resp, err := p.client.R().
		SetHeader("Referer", fmt.Sprintf("https://m.douban.com/movie/celebrity/%s/", id)).
		SetPathParam("id", id).
		SetContext(ctx).
		Get(withProxy("https://m.douban.com/rexxar/api/v2/celebrity/{id}/works"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.NotFound(fmt.Sprintf("the person does not exist"))
	}

	if resp.IsStatusFailure() {
		return nil, errors.Internal(fmt.Sprintf("can not get movie credit for %s", id))
	}

	return buildMovieCreditSlice(resp.String())
}

func (p *provider) GetVideoSource(ctx context.Context, id string) (string, error) {
	resp, err := p.client.R().
		SetHeader("User-Agent", request.UserAgentMicroMessenger).
		SetHeader("Referer", microMessengerReferer).
		SetQueryParam("apiKey", microMessengerKey).
		SetPathParam("id", id).
		SetContext(ctx).
		Get(withProxy("https://frodo.douban.com/api/v2/trailer/{id}"))
	if err != nil {
		return "", errors.WithStack(err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return "", errors.NotFound(fmt.Sprintf("the video does not exist"))
	}

	if url := gjson.Get(resp.String(), "video_url").String(); url != "" {
		return url, nil
	}

	return "", errors.Internal(fmt.Sprintf("can not get video source for %s", id))
}

func NewProvider(client *request.Client) *provider {
	return &provider{
		client: client,
	}
}
