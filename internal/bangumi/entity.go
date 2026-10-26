package bangumi

const (
	TypeAnime = "anime"

	PersonGenderMale   = "male"
	PersonGenderFemale = "female"
	PersonGenderOther  = "other"
)

type Subject struct {
	ID               int64    `json:"id"`
	Kind             string   `json:"type"`
	Name             string   `json:"name"`
	OriginalName     string   `json:"original_name"`
	Overview         string   `json:"overview"`
	OriginalOverview string   `json:"original_overview"`
	Poster           string   `json:"poster"`
	Thumbnail        string   `json:"thumbnail"`
	Tags             []string `json:"tags"`
	Episodes         int      `json:"episodes"`
	ReleaseDate      string   `json:"release_date"`
	Rating           *Rating  `json:"rating"`
	Url              string   `json:"url"`
}

type Rating struct {
	Value float64 `json:"value"`
	Count int64   `json:"count"`
}

type SubjectSlice struct {
	Items []*Subject `json:"items"`
	Total int64      `json:"total"`
}

type Episode struct {
	ID           int64  `json:"id"`
	SubjectID    int64  `json:"subject_id"`
	Episode      int    `json:"episode"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Overview     string `json:"overview"`
	Duration     int64  `json:"duration"`
	ReleaseDate  string `json:"release_date"`
	Url          string `json:"url"`
}

type EpisodeSlice struct {
	Items []*Episode `json:"items"`
	Total int64      `json:"total"`
}

type Character struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Profile  string         `json:"profile"`
	Overview string         `json:"overview"`
	Actors   []*EmbedPerson `json:"actors"`
	Url      string         `json:"url"`
}

type EmbedPerson struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Profile  string `json:"profile"`
	Overview string `json:"overview"`
	Url      string `json:"url"`
}

type Staff struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
	Role    string `json:"role"`
	Url     string `json:"url"`
}

type Person struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	OriginalName string   `json:"original_name"`
	Profile      string   `json:"profile"`
	Gender       string   `json:"gender"`
	Birthday     string   `json:"birthday"`
	PlaceOfBirth string   `json:"place_of_birth"`
	OtherNames   []string `json:"other_names"`
	Biography    string   `json:"biography"`
	Url          string   `json:"url"`
}

type PersonCredit struct {
	ID               int64  `json:"id"`
	Kind             string `json:"type"`
	Name             string `json:"name"`
	OriginalName     string `json:"original_name"`
	CharacterID      int64  `json:"character_id"`
	CharacterName    string `json:"characterName"`
	CharacterProfile string `json:"character_profile"`
	Url              string `json:"url"`
}
