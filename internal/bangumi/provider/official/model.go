package official

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"github.com/whaoa/public-api-bridge/internal/bangumi"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

var errorInvalidResponse = errors.New("invalid bangumi api response")

func toSubjectType(kind int) string {
	switch kind {
	case 2:
		return bangumi.TypeAnime
	default:
		return ""
	}
}

func buildSubjectSlice(payload string) (*bangumi.SubjectSlice, error) {
	root := gjson.Parse(payload)

	items := root.Get("data")
	if !items.IsArray() {
		return nil, errorInvalidResponse
	}

	return &bangumi.SubjectSlice{
		Items: lo.Map(items.Array(), func(item gjson.Result, _ int) *bangumi.Subject {
			if subject, err := buildSubject(item.Raw); err == nil {
				return subject
			}
			return &bangumi.Subject{}
		}),
		Total: root.Get("total").Int(),
	}, nil
}

func buildSubject(payload string) (*bangumi.Subject, error) {
	root := gjson.Parse(payload)

	id := root.Get("id").Int()
	if id == 0 {
		return nil, errorInvalidResponse
	}

	overview := ""
	originalOverview := root.Get("summary").String()
	if values := strings.Split(originalOverview, "[简介原文]"); len(values) > 1 {
		overview = strings.TrimSpace(values[0])
		originalOverview = strings.TrimSpace(values[1])
	}

	return &bangumi.Subject{
		ID:               id,
		Kind:             toSubjectType(int(root.Get("type").Int())),
		Name:             root.Get("name_cn").String(),
		OriginalName:     root.Get("name").String(),
		Overview:         overview,
		OriginalOverview: originalOverview,
		Poster:           root.Get("images.large").String(),
		Thumbnail:        root.Get("images.common").String(),
		Tags:             toStringList(root.Get("meta_tags").Raw),
		Episodes:         int(root.Get("total_episodes").Int()),
		ReleaseDate:      root.Get("date").String(),
		Rating: &bangumi.Rating{
			Value: root.Get("rating.score").Float(),
			Count: root.Get("rating.total").Int(),
		},
	}, nil
}

func toStringList(payload string) []string {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) string {
		return item.String()
	})
}

func buildEpisodeSlice(payload string) (*bangumi.EpisodeSlice, error) {
	root := gjson.Parse(payload)

	items := root.Get("data")
	if !items.IsArray() {
		return nil, errorInvalidResponse
	}

	return &bangumi.EpisodeSlice{
		Items: lo.Map(items.Array(), func(item gjson.Result, _ int) *bangumi.Episode {
			return &bangumi.Episode{
				ID:           item.Get("id").Int(),
				SubjectID:    item.Get("subject_id").Int(),
				Episode:      int(item.Get("ep").Int()),
				Name:         item.Get("name_cn").String(),
				OriginalName: item.Get("name").String(),
				Overview:     item.Get("desc").String(),
				Duration:     item.Get("duration_seconds").Int(),
				ReleaseDate:  item.Get("airdate").String(),
			}
		}),
		Total: root.Get("total").Int(),
	}, nil
}

func buildCharacterList(payload string) ([]*bangumi.Character, error) {
	items := gjson.Parse(payload)
	if !items.IsArray() {
		return nil, errorInvalidResponse
	}

	return lo.Map(items.Array(), func(item gjson.Result, _ int) *bangumi.Character {
		return &bangumi.Character{
			ID:       item.Get("id").Int(),
			Name:     item.Get("name").String(),
			Profile:  item.Get("images.medium").String(),
			Overview: item.Get("summary").String(),
			Actors: lo.Map(item.Get("actors").Array(), func(item gjson.Result, _ int) *bangumi.EmbedPerson {
				return &bangumi.EmbedPerson{
					ID:       item.Get("id").Int(),
					Name:     item.Get("name").String(),
					Profile:  item.Get("images.medium").String(),
					Overview: item.Get("short_summary").String(),
				}
			}),
		}
	}), nil
}

func buildStaffList(payload string) ([]*bangumi.Staff, error) {
	items := gjson.Parse(payload)
	if !items.IsArray() {
		return nil, errorInvalidResponse
	}

	return lo.Map(items.Array(), func(item gjson.Result, _ int) *bangumi.Staff {
		return &bangumi.Staff{
			ID:      item.Get("id").Int(),
			Name:    item.Get("name").String(),
			Profile: item.Get("images.medium").String(),
			Role:    item.Get("relation").String(),
		}
	}), nil
}

func buildPerson(payload string) (*bangumi.Person, error) {
	root := gjson.Parse(payload)

	id := root.Get("id").Int()
	if id == 0 {
		return nil, errorInvalidResponse
	}

	name := ""
	place := ""
	names := []string{}

	if values := root.Get("infobox"); values.IsArray() {
		values.ForEach(func(key, value gjson.Result) bool {
			switch value.Get("key").String() {
			case "简体中文名":
				name = value.Get("value").String()
			case "出生地":
				place = value.Get("value").String()
			case "别名":
				if v := value.Get("value"); v.IsArray() {
					v.ForEach(func(key, value gjson.Result) bool {
						if val := value.Get("v").String(); val != "" {
							names = append(names, val)
						}
						return false
					})
				}
			}
			return false
		})
	}

	return &bangumi.Person{
		ID:           id,
		Name:         name,
		OriginalName: root.Get("name").String(),
		Profile:      root.Get("images.medium").String(),
		Gender:       root.Get("gender").String(),
		Birthday: fmt.Sprintf(
			"%s-%s-%s",
			root.Get("birth_year").String(),
			root.Get("birth_mon").String(),
			root.Get("birth_day").String(),
		),
		PlaceOfBirth: place,
		OtherNames:   names,
		Biography:    root.Get("summary").String(),
	}, nil
}

func buildPersonCreditList(payload string) ([]*bangumi.PersonCredit, error) {
	items := gjson.Parse(payload)
	if !items.IsArray() {
		return nil, errorInvalidResponse
	}

	return lo.FilterMap(items.Array(), func(item gjson.Result, _ int) (*bangumi.PersonCredit, bool) {
		kind := toSubjectType(int(item.Get("type").Int()))
		if kind == "" {
			return nil, false
		}
		return &bangumi.PersonCredit{
			ID:               item.Get("subject_id").Int(),
			Kind:             kind,
			Name:             item.Get("subject_name_cn").String(),
			OriginalName:     item.Get("subject_name").String(),
			CharacterID:      item.Get("id").Int(),
			CharacterName:    item.Get("name").String(),
			CharacterProfile: item.Get("images.medium").String(),
		}, true
	}), nil
}
