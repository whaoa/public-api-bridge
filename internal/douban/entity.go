package douban

import (
	"time"
)

const (
	SubjectTypeMovie = "movie"
	SubjectTypeTv    = "tv"

	ImageSizePoster = "s_ratio_poster"
	ImageSizeSmall  = "s"
	ImageSizeMedium = "m"
	ImageSizeLarge  = "l"

	VideoTypeVideo      = "video"
	VideoTypeTrailer    = "trailer"
	VideoTypeClip       = "clip"
	VideoTypeFeaturette = "featurette"

	PersonGenderMale   = "male"
	PersonGenderFemale = "female"
	PersonGenderOther  = "other"
)

type Subject struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Year        int     `json:"year"`
	Poster      string  `json:"poster"`
	Kind        string  `json:"type"`
	Rating      *Rating `json:"rating"`
	Url         string  `json:"url"`
}

type Rating struct {
	Start float64 `json:"start"`
	Value float64 `json:"value"`
	Max   float64 `json:"max"`
	Count int64   `json:"count"`
}

type Movie struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Poster      string   `json:"poster"`
	Thumbnail   string   `json:"thumbnail"`
	Description string   `json:"description"`
	Intro       string   `json:"intro"`
	Episodes    int      `json:"episodes"`
	Genres      []string `json:"genres"`
	Directors   []string `json:"directors"`
	Actors      []string `json:"actors"`
	Countries   []string `json:"countries"`
	Languages   []string `json:"languages"`
	ReleaseAt   []string `json:"release_at"`
	Rating      *Rating  `json:"rating"`
	Year        int      `json:"year"`
	Kind        string   `json:"type"`
	Url         string   `json:"url"`
}

type Photo struct {
	ID         string    `json:"id"`
	Large      string    `json:"large"`
	Normal     string    `json:"normal"`
	Small      string    `json:"small"`
	CreateTime time.Time `json:"create_time"`
	Position   int       `json:"position"`
	Likes      int64     `json:"likes"`
	Url        string    `json:"url"`
}

type PhotoSlice struct {
	Items []*Photo `json:"items"`
	Total int      `json:"total"`
}

type Video struct {
	ID         string    `json:"id"`
	Kind       string    `json:"type"`
	Title      string    `json:"title"`
	Source     string    `json:"source"`
	Thumbnail  string    `json:"thumbnail"`
	Duration   int64     `json:"duration"`
	CreateTime time.Time `json:"create_time"`
	Url        string    `json:"url"`
}

type MovieCasts struct {
	Directors []*MovieCast `json:"directors"`
	Actors    []*MovieCast `json:"actors"`
}

type MovieCast struct {
	ID          string   `json:"id"`
	Avatar      string   `json:"avatar"`
	Name        string   `json:"name"`
	LatinName   string   `json:"latin_name"`
	Description string   `json:"description"`
	Character   string   `json:"character"`
	Roles       []string `json:"roles"`
	Url         string   `json:"url"`
}

type MovieRecommendation struct {
	ID          string  `json:"id"`
	Kind        string  `json:"type"`
	Title       string  `json:"title"`
	Poster      string  `json:"poster"`
	Description string  `json:"description"`
	Rating      *Rating `json:"rating"`
	Url         string  `json:"url"`
}

type Person struct {
	ID            string   `json:"id"`
	Photo         string   `json:"photo"`
	Name          string   `json:"name"`
	LatinName     string   `json:"latin_name"`
	Description   string   `json:"description"`
	OtherNames    []string `json:"other_names"`
	Gender        string   `json:"gender"`
	Birthday      string   `json:"birthday"`
	PlaceOfBirth  string   `json:"place_of_birth"`
	FamilyMembers []string `json:"family_members"`
	IMDb          string   `json:"IMDb"`
	Url           string   `json:"url"`
}

type MovieCredit struct {
	ID          string   `json:"id"`
	Kind        string   `json:"type"`
	Title       string   `json:"title"`
	Poster      string   `json:"poster"`
	Thumbnail   string   `json:"thumbnail"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	Rating      *Rating  `json:"rating"`
	Directors   []string `json:"directors"`
	Actors      []string `json:"actors"`
	ReleaseAt   []string `json:"release_at"`
	Year        int      `json:"year"`
	Roles       []string `json:"roles"`
	Url         string   `json:"url"`
}

type MovieCreditSlice struct {
	Items []*MovieCredit `json:"items"`
	Total int            `json:"total"`
}
