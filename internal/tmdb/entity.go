package tmdb

const (
	TypeMovie = "movie"
	TypeTV    = "tv"

	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"

	ImageSizeRaw = "original"

	PosterSizeLarge  = "w780"
	PosterSizeMedium = "w342"
	PosterSizeSmall  = "w154"

	BackdropSizeLarge  = "w1280"
	BackdropSizeMedium = "w780"
	BackdropSizeSmall  = "w300"

	StillSizeLarge  = "w500"
	StillSizeMedium = "w300"
	StillSizeSmall  = "w185"

	LogoSizeLarge  = "w92"
	LogoSizeMedium = "w154"
	LogoSizeSmall  = "w500"

	ProfileSizeLarge  = "h632"
	ProfileSizeMedium = "w185"
	ProfileSizeSmall  = "w45"
)

type Media struct {
	ID               int64   `json:"id"`
	Kind             string  `json:"type"`
	Name             string  `json:"name"`
	OriginalName     string  `json:"original_name"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	Poster           string  `json:"poster"`
	ReleaseDate      string  `json:"release_date"`
	Rating           float64 `json:"rating"`
	Url              string  `json:"url"`
}

type MediaSlice struct {
	Total int      `json:"total"`
	Items []*Media `json:"items"`
}

type Language struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
}

type Country struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Company struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	Country string `json:"country"`
	Url     string `json:"url"`
}

type Movie struct {
	ID                  int64       `json:"id"`
	IMDbID              string      `json:"IMDb"`
	Name                string      `json:"name"`
	OriginalName        string      `json:"original_name"`
	OriginalLanguage    string      `json:"original_language"`
	Poster              string      `json:"poster"`
	Backdrop            string      `json:"backdrop"`
	Genres              []string    `json:"genres"`
	Overview            string      `json:"overview"`
	ReleaseDate         string      `json:"release_date"`
	Rating              float64     `json:"rating"`
	Runtime             int         `json:"runtime"`
	SpokenLanguages     []*Language `json:"spoken_languages"`
	ProductionCompanies []*Company  `json:"production_companies"`
	ProductionCountries []*Country  `json:"production_countries"`
	Url                 string      `json:"url"`
}

type TV struct {
	ID                  int64          `json:"id"`
	IMDbID              string         `json:"IMDb"`
	Name                string         `json:"name"`
	OriginalName        string         `json:"original_name"`
	OriginalLanguage    string         `json:"original_language"`
	Poster              string         `json:"poster"`
	Backdrop            string         `json:"backdrop"`
	Genres              []string       `json:"genres"`
	Overview            string         `json:"overview"`
	ReleaseDate         string         `json:"release_date"`
	Rating              float64        `json:"rating"`
	SpokenLanguages     []*Language    `json:"spoken_languages"`
	ProductionCompanies []*Company     `json:"production_companies"`
	ProductionCountries []*Country     `json:"production_countries"`
	Seasons             []*EmbedSeason `json:"seasons"`
	Url                 string         `json:"url"`
}

type EmbedSeason struct {
	Season   int    `json:"season"`
	Name     string `json:"name"`
	Overview string `json:"overview"`
	Poster   string `json:"poster"`
	Episodes int    `json:"episodes"`
}

type Season struct {
	ID          int64      `json:"id"`
	Season      int        `json:"season"`
	Name        string     `json:"name"`
	Poster      string     `json:"poster"`
	Overview    string     `json:"overview"`
	ReleaseDate string     `json:"release_date"`
	Rating      float64    `json:"rating"`
	Episodes    []*Episode `json:"episodes"`
}

type Episode struct {
	Season      int     `json:"season"`
	Episode     int     `json:"episode"`
	Name        string  `json:"name"`
	Still       string  `json:"still"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	Rating      float64 `json:"rating"`
}

type Image struct {
	Url      string `json:"url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Country  string `json:"country"`
	Language string `json:"language"`
}

type MediaImages struct {
	Logos     []*Image `json:"logos"`
	Posters   []*Image `json:"posters"`
	Backdrops []*Image `json:"backdrops"`
}

type MediaCast struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Profile      string `json:"profile"`
	Gender       string `json:"gender"`
	Role         string `json:"role"`
	Character    string `json:"character"`
	Url          string `json:"url"`
}

type MediaCasts struct {
	Actors []*MediaCast `json:"actors"`
	Staffs []*MediaCast `json:"staffs"`
}

type EmbedPerson struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Profile      string `json:"profile"`
	Gender       string `json:"gender"`
	Url          string `json:"url"`
}

type EmbedPersonSlice struct {
	Total int            `json:"total"`
	Items []*EmbedPerson `json:"items"`
}

type Person struct {
	ID           int64    `json:"id"`
	IMDbID       string   `json:"IMDb"`
	Name         string   `json:"name"`
	Profile      string   `json:"profile"`
	Gender       string   `json:"gender"`
	Birthday     string   `json:"birthday"`
	PlaceOfBirth string   `json:"place_of_birth"`
	OtherNames   []string `json:"other_names"`
	Biography    string   `json:"biography"`
	Url          string   `json:"url"`
}

type MediaCredit struct {
	*Media
	Role      string `json:"role"`
	Character string `json:"character"`
}

type PersonCredits struct {
	AsActor []*MediaCredit `json:"as_actor"`
	AsStaff []*MediaCredit `json:"as_staff"`
}
