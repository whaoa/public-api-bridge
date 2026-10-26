package web

import (
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"github.com/whaoa/public-api-bridge/internal/bangumi"
	"github.com/whaoa/public-api-bridge/pkg/errors"
	"github.com/whaoa/public-api-bridge/pkg/util"
)

var errorInvalidResponse = errors.New("invalid bangumi web api response")

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

	date := ""
	episodes := 0
	overview := root.Get("info").String()
	lo.ForEach(strings.Split(overview, " / "), func(value string, _ int) {
		switch true {
		case strings.HasSuffix(value, "话"):
			episodes, _ = strconv.Atoi(strings.TrimSuffix(value, "话"))
		case strings.Contains(value, "月"):
			if t, err := time.Parse("2006年1月2日", value); err == nil {
				date = t.Format("2006-01-02")
			}
		}
	})

	return &bangumi.Subject{
		ID:               id,
		Kind:             toSubjectType(int(root.Get("type").Int())),
		Name:             root.Get("nameCN").String(),
		OriginalName:     root.Get("name").String(),
		Overview:         overview,
		OriginalOverview: "",
		Poster:           root.Get("images.large").String(),
		Thumbnail:        root.Get("images.common").String(),
		Tags:             toStringList(root.Get("metaTags").Raw),
		Episodes:         episodes,
		ReleaseDate:      date,
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
				SubjectID:    item.Get("subjectID").Int(),
				Episode:      int(item.Get("sort").Int()),
				Name:         item.Get("nameCN").String(),
				OriginalName: item.Get("name").String(),
				Overview:     item.Get("desc").String(),
				Duration:     util.NormalDuration(item.Get("duration").String()),
				ReleaseDate:  item.Get("airdate").String(),
			}
		}),
		Total: root.Get("total").Int(),
	}, nil
}

func buildCharacterList(payload string) ([]*bangumi.Character, error) {
	items := gjson.Get(payload, "data")
	if !items.IsArray() {
		return nil, errorInvalidResponse
	}

	return lo.Map(items.Array(), func(item gjson.Result, _ int) *bangumi.Character {
		return &bangumi.Character{
			ID:       item.Get("character.id").Int(),
			Name:     item.Get("character.name").String(),
			Profile:  item.Get("character.images.medium").String(),
			Overview: item.Get("info").String(),
			Actors: lo.Map(item.Get("casts").Array(), func(item gjson.Result, _ int) *bangumi.EmbedPerson {
				return &bangumi.EmbedPerson{
					ID:       item.Get("person.id").Int(),
					Name:     item.Get("person.name").String(),
					Profile:  item.Get("person.images.medium").String(),
					Overview: item.Get("person.info").String(),
				}
			}),
		}
	}), nil
}

func buildStaffList(payload string) ([]*bangumi.Staff, error) {
	items := gjson.Get(payload, "data")
	if !items.IsArray() {
		return nil, errorInvalidResponse
	}

	return lo.Map(items.Array(), func(item gjson.Result, _ int) *bangumi.Staff {
		return &bangumi.Staff{
			ID:      item.Get("staff.id").Int(),
			Name:    item.Get("staff.name").String(),
			Profile: item.Get("staff.images.medium").String(),
			Role:    item.Get("positions.0.type.cn").String(),
		}
	}), nil
}

func buildPerson(payload string) (*bangumi.Person, error) {
	root := gjson.Parse(payload)

	id := root.Get("id").Int()
	if id == 0 {
		return nil, errorInvalidResponse
	}

	gender := bangumi.PersonGenderOther
	birthday := ""
	place := ""
	names := []string{}

	if values := root.Get("infobox"); values.IsArray() {
		values.ForEach(func(key, value gjson.Result) bool {
			switch value.Get("key").String() {
			case "生日":
				birthday = value.Get("values.0.v").String()
			case "出生地":
				place = value.Get("values.0.v").String()
			case "性别":
				switch value.Get("values.0.v").String() {
				case "男":
					gender = bangumi.PersonGenderMale
				case "女":
					gender = bangumi.PersonGenderFemale
				}
			case "别名":
				if v := value.Get("values"); v.IsArray() {
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
		Name:         root.Get("nameCN").String(),
		OriginalName: root.Get("name").String(),
		Profile:      root.Get("images.medium").String(),
		Gender:       gender,
		Birthday:     birthday,
		PlaceOfBirth: place,
		OtherNames:   names,
		Biography:    root.Get("summary").String(),
	}, nil
}

func buildPersonCreditList(payload string) ([]*bangumi.PersonCredit, error) {
	items := gjson.Get(payload, "data")
	if !items.IsArray() {
		return nil, errorInvalidResponse
	}

	return lo.FilterMap(items.Array(), func(item gjson.Result, _ int) (*bangumi.PersonCredit, bool) {
		subject, exist := lo.Find(item.Get("relations").Array(), func(item gjson.Result) bool {
			return toSubjectType(int(item.Get("subject.type").Int())) == bangumi.TypeAnime
		})
		if !exist {
			return nil, false
		}
		return &bangumi.PersonCredit{
			ID:               subject.Get("id").Int(),
			Kind:             toSubjectType(int(subject.Get("type").Int())),
			Name:             subject.Get("name").String(),
			OriginalName:     subject.Get("nameCn").String(),
			CharacterID:      item.Get("character.id").Int(),
			CharacterName:    item.Get("character.name").String(),
			CharacterProfile: item.Get("character.images.medium").String(),
		}, true
	}), nil
}
