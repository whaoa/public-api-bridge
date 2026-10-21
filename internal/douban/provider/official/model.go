package official

import (
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"github.com/whaoa/public-api-bridge/internal/douban"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

var errorInvalidResponse = errors.Internal("invalid douban api response")

func buildSubjectList(payload string) ([]*douban.Subject, error) {
	root := gjson.Get(payload, "subjects.items")
	if !root.IsArray() {
		return nil, errorInvalidResponse
	}

	return lo.MapErr(root.Array(), func(item gjson.Result, _ int) (*douban.Subject, error) {
		id := item.Get("target.id")
		kind := item.Get("target_type")
		if !id.Exists() || !kind.Exists() {
			return nil, errorInvalidResponse
		}

		url := ""
		switch kind.String() {
		case douban.SubjectTypeTv:
			fallthrough
		case douban.SubjectTypeMovie:
			url = fmt.Sprintf("https://movie.douban.com/subject/%s/", id.String())
		}

		return &douban.Subject{
			ID:          id.String(),
			Kind:        kind.String(),
			Title:       item.Get("target.title").String(),
			Description: item.Get("target.card_subtitle").String(),
			Year:        int(item.Get("target.year").Int()),
			Poster:      item.Get("target.cover_url").String(),
			Rating:      toRating(item.Get("target.rating").Raw),
			Url:         url,
		}, nil
	})
}

func toRating(payload string) *douban.Rating {
	root := gjson.Parse(payload)
	return &douban.Rating{
		Start: root.Get("star_count").Float(),
		Value: root.Get("value").Float(),
		Max:   root.Get("max").Float(),
		Count: root.Get("count").Int(),
	}
}

func buildMovie(payload string) (*douban.Movie, error) {
	root := gjson.Parse(payload)

	id := root.Get("id")
	if !id.Exists() {
		return nil, errorInvalidResponse
	}

	return &douban.Movie{
		ID:          id.String(),
		Kind:        root.Get("type").String(),
		Title:       root.Get("title").String(),
		Poster:      root.Get("cover_url").String(),
		Description: root.Get("card_subtitle").String(),
		Intro:       root.Get("intro").String(),
		Episodes:    int(root.Get("episodes_count").Int()),
		Directors:   toPersonNameList(root.Get("directors").Raw),
		Actors:      toPersonNameList(root.Get("actors").Raw),
		Genres:      toStringList(root.Get("genres").Raw),
		Countries:   toStringList(root.Get("countries").Raw),
		Languages:   toStringList(root.Get("languages").Raw),
		ReleaseAt:   toStringList(root.Get("pubdate").Raw),
		Rating:      toRating(root.Get("rating").Raw),
		Year:        int(root.Get("year").Int()),
		Url:         root.Get("url").String(),
	}, nil
}

func toStringList(payload string) []string {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) string {
		return item.String()
	})
}

func toPersonNameList(payload string) []string {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) string {
		return item.Get("name").String()
	})
}

func buildPhotoSlice(payload string) (*douban.PhotoSlice, error) {
	root := gjson.Parse(payload)

	photos := root.Get("photos")
	if !photos.IsArray() {
		return nil, errorInvalidResponse
	}

	return &douban.PhotoSlice{
		Total: int(root.Get("total").Int()),
		Items: lo.Map(photos.Array(), func(item gjson.Result, _ int) *douban.Photo {
			return &douban.Photo{
				ID:         item.Get("id").String(),
				Large:      item.Get("image.large.url").String(),
				Normal:     item.Get("image.normal.url").String(),
				Small:      item.Get("image.small.url").String(),
				CreateTime: normalizeDateTime(item.Get("create_time").String()),
				Position:   int(item.Get("position").Int()),
				Likes:      item.Get("likers_count").Int(),
				Url:        item.Get("url").String(),
			}
		}),
	}, nil
}

func buildVideoList(payload string) ([]*douban.Video, error) {
	root := gjson.Get(payload, "trailers")
	if !root.IsArray() {
		return nil, errorInvalidResponse
	}

	return lo.Map(root.Array(), func(item gjson.Result, _ int) *douban.Video {
		id := item.Get("id").String()

		kind := douban.VideoTypeVideo
		switch item.Get("type").String() {
		case "A":
			kind = douban.VideoTypeTrailer
		case "B":
			kind = douban.VideoTypeClip
		case "C":
			kind = douban.VideoTypeFeaturette
		}

		return &douban.Video{
			ID:         id,
			Title:      item.Get("title").String(),
			Source:     item.Get("video_url").String(),
			Thumbnail:  item.Get("cover_url").String(),
			Duration:   normalDuration(item.Get("runtime").String()),
			CreateTime: normalizeDateTime(item.Get("create_time").String()),
			Kind:       kind,
			Url:        fmt.Sprintf("https://movie.douban.com/trailer/%s/", id),
		}
	}), nil
}

func buildMovieCasts(payload string) (*douban.MovieCasts, error) {
	root := gjson.Parse(payload)

	directors := root.Get("directors")
	if !directors.IsArray() {
		return nil, errorInvalidResponse
	}
	actors := root.Get("actors")
	if !actors.IsArray() {
		return nil, errorInvalidResponse
	}

	return &douban.MovieCasts{
		Directors: toMovieCastList(directors.Raw),
		Actors:    toMovieCastList(actors.Raw),
	}, nil
}

func toMovieCastList(payload string) []*douban.MovieCast {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) *douban.MovieCast {
		id := item.Get("id").String()
		return &douban.MovieCast{
			ID:          id,
			Avatar:      item.Get("avatar.normal").String(),
			Name:        item.Get("name").String(),
			LatinName:   item.Get("latin_name").String(),
			Description: item.Get("title").String(),
			Character:   item.Get("character").String(),
			Roles:       toStringList(item.Get("roles").Raw),
			Url:         fmt.Sprintf("https://movie.douban.com/celebrity/%s/", id),
		}
	})
}

func buildMovieRecommendations(payload string) ([]*douban.MovieRecommendation, error) {
	root := gjson.Parse(payload)
	if !root.IsArray() {
		return nil, errorInvalidResponse
	}

	return lo.Map(root.Array(), func(item gjson.Result, _ int) *douban.MovieRecommendation {
		return &douban.MovieRecommendation{
			ID:          item.Get("id").String(),
			Kind:        item.Get("type").String(),
			Title:       item.Get("title").String(),
			Poster:      item.Get("pic.normal").String(),
			Description: item.Get("card_subtitle").String(),
			Rating:      toRating(item.Get("rating").Raw),
			Url:         item.Get("url").String(),
		}
	}), nil
}

func buildPerson(id, payload string) (*douban.Person, error) {
	root := gjson.Parse(payload)

	if id := root.Get("id").String(); id == "" {
		return nil, errorInvalidResponse
	}

	gender := ""
	birthday := ""
	place := ""
	names := []string{}
	members := []string{}
	IMDb := ""

	if info := root.Get("extra.info"); info.IsArray() {
		info.ForEach(func(_, item gjson.Result) bool {
			key := item.Get("0").String()
			value := item.Get("1").String()

			switch true {
			case strings.Contains(key, "IMDb"):
				IMDb = value
			case strings.Contains(key, "性别"):
				if value == "女" {
					gender = douban.PersonGenderFemale
				} else if value == "男" {
					gender = douban.PersonGenderMale
				}
			case strings.Contains(key, "出生日期"):
				if t, err := time.Parse("2006年1月2日", value); err == nil {
					birthday = t.Format("2006-01-02")
				}
			case strings.Contains(key, "出生地"):
				place = value
			case strings.Contains(key, "外文名"):
				fallthrough
			case strings.Contains(key, "中文名"):
				if value != "" {
					names = append(names, strings.Split(value, " / ")...)
				}
			case strings.Contains(key, "家庭成员"):
				if value != "" {
					members = append(members, strings.Split(value, " / ")...)
				}
			}
			return true
		})
	}

	return &douban.Person{
		ID:            id,
		Photo:         root.Get("cover.large.url").String(),
		Name:          root.Get("title").String(),
		LatinName:     root.Get("latin_title").String(),
		Description:   root.Get("extra.short_info").String(),
		OtherNames:    names,
		Gender:        gender,
		Birthday:      birthday,
		PlaceOfBirth:  place,
		FamilyMembers: members,
		IMDb:          IMDb,
		Url:           root.Get("url").String(),
	}, nil
}

func buildMovieCreditSlice(payload string) (*douban.MovieCreditSlice, error) {
	root := gjson.Parse(payload)

	works := root.Get("works")
	if !works.IsArray() {
		return nil, errorInvalidResponse
	}

	return &douban.MovieCreditSlice{
		Total: int(root.Get("total").Int()),
		Items: lo.Map(works.Array(), func(item gjson.Result, _ int) *douban.MovieCredit {
			return &douban.MovieCredit{
				ID:          item.Get("work.id").String(),
				Kind:        item.Get("work.type").String(),
				Title:       item.Get("work.title").String(),
				Poster:      item.Get("work.cover_url").String(),
				Description: item.Get("work.card_subtitle").String(),
				Genres:      toStringList(item.Get("work.genres").Raw),
				Rating:      toRating(item.Get("work.rating").Raw),
				Directors:   toPersonNameList(item.Get("work.directors").Raw),
				Actors:      toPersonNameList(item.Get("work.actors").Raw),
				ReleaseAt:   toStringList(item.Get("work.pubdate").Raw),
				Year:        int(item.Get("work.year").Int()),
				Roles:       toStringList(item.Get("roles").Raw),
				Url:         item.Get("work.url").String(),
			}
		}),
	}, nil
}
