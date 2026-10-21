package tmdb

import "fmt"

const (
	webUrl       = "https://www.themoviedb.org"
	imageBaseUrl = "https://image.tmdb.org/t/p"

	urlTypePerson  = "person"
	urlTypeCompany = "company"
)

func buildImageUrl(path, size string) string {
	if path == "" {
		return ""
	}
	if size == "" {
		size = ImageSizeRaw
	}
	return fmt.Sprintf("https://image.tmdb.org/t/p/%s%s", size, path)
}

func buildRefUrl(kind string, id int64) string {
	return fmt.Sprintf("%s/%s/%d", webUrl, kind, id)
}

func ValidateMediaType(kind string) bool {
	switch kind {
	case TypeTV:
		fallthrough
	case TypeMovie:
		return true
	default:
		return false
	}
}
