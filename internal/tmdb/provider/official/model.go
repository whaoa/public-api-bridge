package official

import (
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"github.com/whaoa/public-api-bridge/internal/tmdb"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

var errorInvalidResponse = errors.New("invalid tmdb api response")

func buildMediaSlice(payload string) (*tmdb.MediaSlice, error) {
	root := gjson.Parse(payload)

	items := root.Get("results")
	if !items.IsArray() {
		return nil, errorInvalidResponse
	}

	return &tmdb.MediaSlice{
		Total: int(root.Get("total_results").Int()),
		Items: lo.Map(items.Array(), func(item gjson.Result, _ int) *tmdb.Media {
			name := item.Get("name").String()
			originalName := item.Get("original_name").String()
			releaseDate := item.Get("first_air_date").String()
			if name == "" {
				name = item.Get("title").String()
			}
			if originalName == "" {
				originalName = item.Get("original_title").String()
			}
			if releaseDate == "" {
				releaseDate = item.Get("first_air_date").String()
			}
			return &tmdb.Media{
				ID:               item.Get("id").Int(),
				Kind:             item.Get("media_type").String(),
				Name:             name,
				OriginalName:     originalName,
				OriginalLanguage: item.Get("original_language").String(),
				Overview:         item.Get("overview").String(),
				Poster:           item.Get("poster_path").String(),
				ReleaseDate:      releaseDate,
				Rating:           item.Get("vote_average").Float(),
			}
		}),
	}, nil
}

func buildMovie(payload string) (*tmdb.Movie, error) {
	root := gjson.Parse(payload)

	id := root.Get("id").Int()
	if id == 0 {
		return nil, errorInvalidResponse
	}

	return &tmdb.Movie{
		ID:                  id,
		IMDbID:              root.Get("imdb_id").String(),
		Name:                root.Get("title").String(),
		OriginalName:        root.Get("original_title").String(),
		OriginalLanguage:    root.Get("original_language").String(),
		Poster:              root.Get("poster_path").String(),
		Backdrop:            root.Get("backdrop_path").String(),
		Genres:              toGenreList(root.Get("genres").Raw),
		Overview:            root.Get("overview").String(),
		ReleaseDate:         root.Get("release_date").String(),
		Rating:              root.Get("vote_average").Float(),
		Runtime:             int(root.Get("runtime").Int()),
		SpokenLanguages:     toLanguageList(root.Get("spoken_languages").Raw),
		ProductionCountries: toCountryList(root.Get("production_countries").Raw),
		ProductionCompanies: toCompanyList(root.Get("production_companies").Raw),
	}, nil
}

func toGenreList(payload string) []string {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) string {
		return item.Get("name").String()
	})
}

func toLanguageList(payload string) []*tmdb.Language {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) *tmdb.Language {
		return &tmdb.Language{
			ID:          item.Get("iso_639_1").String(),
			Name:        item.Get("name").String(),
			EnglishName: item.Get("english_name").String(),
		}
	})
}

func toCountryList(payload string) []*tmdb.Country {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) *tmdb.Country {
		return &tmdb.Country{
			ID:   item.Get("iso_3166_1").String(),
			Name: item.Get("name").String(),
		}
	})
}

func toCompanyList(payload string) []*tmdb.Company {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) *tmdb.Company {
		return &tmdb.Company{
			ID:      item.Get("id").Int(),
			Name:    item.Get("name").String(),
			Logo:    item.Get("logo_path").String(),
			Country: item.Get("origin_country").String(),
		}
	})
}

func buildTV(payload string) (*tmdb.TV, error) {
	root := gjson.Parse(payload)

	id := root.Get("id").Int()
	if id == 0 {
		return nil, errorInvalidResponse
	}

	return &tmdb.TV{
		ID:                  id,
		IMDbID:              root.Get("external_ids.imdb_id").String(),
		Name:                root.Get("name").String(),
		OriginalName:        root.Get("original_name").String(),
		OriginalLanguage:    root.Get("original_language").String(),
		Poster:              root.Get("poster_path").String(),
		Backdrop:            root.Get("backdrop_path").String(),
		Genres:              toGenreList(root.Get("genres").Raw),
		Overview:            root.Get("overview").String(),
		ReleaseDate:         root.Get("first_air_date").String(),
		Rating:              root.Get("vote_average").Float(),
		SpokenLanguages:     toLanguageList(root.Get("spoken_languages").Raw),
		ProductionCountries: toCountryList(root.Get("production_countries").Raw),
		ProductionCompanies: toCompanyList(root.Get("production_companies").Raw),
		Seasons:             toEmbedSeasonList(root.Get("seasons").Raw),
	}, nil
}

func toEmbedSeasonList(payload string) []*tmdb.EmbedSeason {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) *tmdb.EmbedSeason {
		return &tmdb.EmbedSeason{
			Season:   int(item.Get("season_number").Int()),
			Name:     item.Get("name").String(),
			Overview: item.Get("overview").String(),
			Poster:   item.Get("poster_path").String(),
			Episodes: int(item.Get("episode_count").Int()),
		}
	})
}

func buildSeason(payload string) (*tmdb.Season, error) {
	root := gjson.Parse(payload)

	id := root.Get("id").Int()
	if id == 0 {
		return nil, errorInvalidResponse
	}

	return &tmdb.Season{
		ID:          id,
		Season:      int(root.Get("season_number").Int()),
		Name:        root.Get("name").String(),
		Poster:      root.Get("poster_path").String(),
		Overview:    root.Get("overview").String(),
		ReleaseDate: root.Get("air_date").String(),
		Rating:      root.Get("vote_average").Float(),
		Episodes:    toEpisodeList(root.Get("episodes").Raw),
	}, nil
}

func toEpisodeList(payload string) []*tmdb.Episode {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) *tmdb.Episode {
		return &tmdb.Episode{
			Season:      int(item.Get("season_number").Int()),
			Episode:     int(item.Get("episode_number").Int()),
			Name:        item.Get("name").String(),
			Still:       item.Get("still_path").String(),
			Overview:    item.Get("overview").String(),
			ReleaseDate: item.Get("air_date").String(),
			Rating:      item.Get("vote_average").Float(),
		}
	})
}

func buildMediaImages(payload string) (*tmdb.MediaImages, error) {
	root := gjson.Parse(payload)

	return &tmdb.MediaImages{
		Logos:     toImageList(root.Get("logos").Raw),
		Posters:   toImageList(root.Get("posters").Raw),
		Backdrops: toImageList(root.Get("backdrops").Raw),
	}, nil
}

func toImageList(payload string) []*tmdb.Image {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) *tmdb.Image {
		return &tmdb.Image{
			Url:      item.Get("file_path").String(),
			Width:    int(item.Get("width").Int()),
			Height:   int(item.Get("height").Int()),
			Country:  item.Get("iso_3166_1").String(),
			Language: item.Get("iso_639_1").String(),
		}
	})
}

func buildMediaCasts(payload string) (*tmdb.MediaCasts, error) {
	root := gjson.Parse(payload)

	if id := root.Get("id").Int(); id == 0 {
		return nil, errorInvalidResponse
	}

	return &tmdb.MediaCasts{
		Actors: toMediaCastList(root.Get("cast").Raw),
		Staffs: toMediaCastList(root.Get("crew").Raw),
	}, nil
}

func toMediaCastList(payload string) []*tmdb.MediaCast {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) *tmdb.MediaCast {
		role := item.Get("job").String()
		character := item.Get("character").String()

		if role == "" && character != "" {
			role = "Actor"
		}

		return &tmdb.MediaCast{
			ID:           item.Get("id").Int(),
			Name:         item.Get("name").String(),
			OriginalName: item.Get("original_name").String(),
			Profile:      item.Get("profile_path").String(),
			Gender:       toGender(item.Get("gender").Raw),
			Role:         role,
			Character:    character,
		}
	})
}

func toGender(payload string) string {
	switch payload {
	case "1":
		return tmdb.GenderFemale
	case "2":
		return tmdb.GenderMale
	default:
		return tmdb.GenderOther
	}
}

func buildEmbedPersonSlice(payload string) (*tmdb.EmbedPersonSlice, error) {
	root := gjson.Parse(payload)

	if id := root.Get("id").Int(); id == 0 {
		return nil, errorInvalidResponse
	}

	return &tmdb.EmbedPersonSlice{
		Total: int(root.Get("total_results").Int()),
		Items: lo.Map(root.Get("results").Array(), func(item gjson.Result, _ int) *tmdb.EmbedPerson {
			return &tmdb.EmbedPerson{
				ID:           item.Get("id").Int(),
				Name:         item.Get("name").String(),
				OriginalName: item.Get("original_name").String(),
				Profile:      item.Get("profile_path").String(),
				Gender:       toGender(item.Get("gender").Raw),
			}
		}),
	}, nil
}

func buildPerson(payload string) (*tmdb.Person, error) {
	root := gjson.Parse(payload)

	id := root.Get("id").Int()
	if id == 0 {
		return nil, errorInvalidResponse
	}

	return &tmdb.Person{
		ID:           id,
		IMDbID:       root.Get("imdb_id").String(),
		Name:         root.Get("name").String(),
		Profile:      root.Get("profile_path").String(),
		Gender:       toGender(root.Get("gender").Raw),
		Birthday:     root.Get("birthday").String(),
		PlaceOfBirth: root.Get("place_of_birth").String(),
		OtherNames:   toStringList(root.Get("also_known_as").Raw),
		Biography:    root.Get("biography").String(),
	}, nil
}

func toStringList(payload string) []string {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) string {
		return item.String()
	})
}

func buildPersonCredits(payload string) (*tmdb.PersonCredits, error) {
	root := gjson.Parse(payload)

	if id := root.Get("id").Int(); id == 0 {
		return nil, errorInvalidResponse
	}

	return &tmdb.PersonCredits{
		AsActor: toMediaCreditList(root.Get("cast").Raw),
		AsStaff: toMediaCreditList(root.Get("crew").Raw),
	}, nil
}

func toMediaCreditList(payload string) []*tmdb.MediaCredit {
	return lo.Map(gjson.Parse(payload).Array(), func(item gjson.Result, _ int) *tmdb.MediaCredit {
		role := item.Get("job").String()
		character := item.Get("character").String()
		if role == "" {
			role = "Actor"
		}

		return &tmdb.MediaCredit{
			Role:      role,
			Character: character,
			Media: &tmdb.Media{
				ID:               item.Get("id").Int(),
				Kind:             item.Get("media_type").String(),
				Name:             item.Get("title").String(),
				OriginalName:     item.Get("original_title").String(),
				OriginalLanguage: item.Get("original_language").String(),
				Overview:         item.Get("overview").String(),
				Poster:           item.Get("poster_path").String(),
				ReleaseDate:      item.Get("release_date").String(),
				Rating:           item.Get("vote_average").Float(),
			},
		}
	})
}

func buildPersonProfiles(payload string) ([]*tmdb.Image, error) {
	profiles := gjson.Get(payload, "profiles")
	if !profiles.IsArray() {
		return nil, errorInvalidResponse
	}
	return toImageList(profiles.Raw), nil
}
