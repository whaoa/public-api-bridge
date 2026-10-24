package douban

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/whaoa/public-api-bridge/internal/server"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

type handler struct {
	service *Service
}

func (h *handler) RegisterRoutes(router server.Router) {
	t := mediaTypeGuardMiddleware
	id := idGuardMiddleware

	router.GET("/douban/:type/subjects", h.SearchMovie, t)

	router.GET("/douban/:type/subjects/:id", h.GetMovie, t, id)
	router.GET("/douban/:type/subjects/:id/photos", h.ListMoviePhotos, t, id)
	router.GET("/douban/:type/subjects/:id/videos", h.GetMovieVideos, t, id)
	router.GET("/douban/:type/subjects/:id/casts", h.GetMovieCasts, t, id)
	router.GET("/douban/:type/subjects/:id/recommendations", h.GetMovieRecommendations, t, id)

	router.GET("/douban/person/:id", h.GetPerson, id)
	router.GET("/douban/person/:id/photos", h.ListPersonPhotos, id)
	router.GET("/douban/person/:id/movies", h.ListPersonMovieCredits, id)

	router.GET("/douban/raw/:type/videos/:id", h.GetVideoSource, t, id)
}

// SearchMovie searches movie and TV subjects.
// @Summary Search Douban subjects
// @Tags Douban
// @Param type path string true "Subject type" Enums(movie,tv)
// @Param keyword query string true "Search keyword"
// @Success 200 {array} Subject
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/{type}/subjects [get]
func (h *handler) SearchMovie(c *echo.Context) error {
	keyword := strings.TrimSpace(c.QueryParam("keyword"))
	result, err := h.service.SearchMovie(c.Request().Context(), keyword)
	return server.JSON(c, result, err)
}

// GetMovie returns a Douban movie or TV subject.
// @Summary Get a Douban subject
// @Tags Douban
// @Param type path string true "Subject type" Enums(movie,tv)
// @Param id path string true "Subject ID"
// @Success 200 {object} Movie
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/{type}/subjects/{id} [get]
func (h *handler) GetMovie(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetMovie(c.Request().Context(), c.Param("type"), id)
	return server.JSON(c, result, err)
}

// ListMoviePhotos returns photos for a Douban subject.
// @Summary List subject photos
// @Tags Douban
// @Param type path string true "Subject type" Enums(movie,tv)
// @Param id path string true "Subject ID"
// @Param offset query int false "Result offset"
// @Param count query int false "Result count"
// @Success 200 {object} PhotoSlice
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/{type}/subjects/{id}/photos [get]
func (h *handler) ListMoviePhotos(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	count, _ := strconv.Atoi(c.QueryParam("count"))
	result, err := h.service.ListMoviePhotos(c.Request().Context(), c.Param("type"), id, offset, count)
	return server.JSON(c, result, err)
}

// GetMovieVideos returns videos for a Douban subject.
// @Summary List subject videos
// @Tags Douban
// @Param type path string true "Subject type" Enums(movie,tv)
// @Param id path string true "Subject ID"
// @Success 200 {array} Video
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/{type}/subjects/{id}/videos [get]
func (h *handler) GetMovieVideos(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetMovieVideos(c.Request().Context(), c.Param("type"), id)
	return server.JSON(c, result, err)
}

// GetMovieCasts returns cast and crew for a Douban subject.
// @Summary Get subject cast
// @Tags Douban
// @Param type path string true "Subject type" Enums(movie,tv)
// @Param id path string true "Subject ID"
// @Success 200 {object} MovieCasts
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/{type}/subjects/{id}/casts [get]
func (h *handler) GetMovieCasts(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetMovieCasts(c.Request().Context(), c.Param("type"), id)
	return server.JSON(c, result, err)
}

// GetMovieRecommendations returns recommendations for a Douban subject.
// @Summary List subject recommendations
// @Tags Douban
// @Param type path string true "Subject type" Enums(movie,tv)
// @Param id path string true "Subject ID"
// @Success 200 {array} MovieRecommendation
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/{type}/subjects/{id}/recommendations [get]
func (h *handler) GetMovieRecommendations(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetMovieRecommendations(c.Request().Context(), c.Param("type"), id)
	return server.JSON(c, result, err)
}

// GetPerson returns a Douban person.
// @Summary Get a Douban person
// @Tags Douban
// @Param id path string true "Person ID"
// @Success 200 {object} Person
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/person/{id} [get]
func (h *handler) GetPerson(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetPerson(c.Request().Context(), id)
	return server.JSON(c, result, err)
}

// ListPersonPhotos returns photos for a Douban person.
// @Summary List person photos
// @Tags Douban
// @Param id path string true "Person ID"
// @Param offset query int false "Result offset"
// @Param count query int false "Result count"
// @Success 200 {object} PhotoSlice
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/person/{id}/photos [get]
func (h *handler) ListPersonPhotos(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	count, _ := strconv.Atoi(c.QueryParam("count"))
	result, err := h.service.ListPersonPhotos(c.Request().Context(), id, offset, count)
	return server.JSON(c, result, err)
}

// ListPersonMovieCredits returns movie credits for a Douban person.
// @Summary List person movie credits
// @Tags Douban
// @Param id path string true "Person ID"
// @Param offset query int false "Result offset"
// @Param count query int false "Result count"
// @Success 200 {object} MovieCreditSlice
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/person/{id}/movies [get]
func (h *handler) ListPersonMovieCredits(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	count, _ := strconv.Atoi(c.QueryParam("count"))
	result, err := h.service.ListPersonMovieCredits(c.Request().Context(), id, offset, count)
	return server.JSON(c, result, err)
}

// GetVideoSource redirects to the source URL for a video.
// @Summary Get video source
// @Tags Douban
// @Produce text/html
// @Param type path string true "Subject type" Enums(movie,tv)
// @Param id path string true "Video ID"
// @Success 307 {string} string
// @Header 307 {string} Location "Redirect target"
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /douban/raw/{type}/videos/{id} [get]
func (h *handler) GetVideoSource(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	source, err := h.service.GetVideoSource(c.Request().Context(), id)
	if err != nil {
		return server.JSON(c, nil, err)
	}
	return c.Redirect(http.StatusTemporaryRedirect, source)
}

func NewHandler(service *Service) *handler {
	return &handler{
		service: service,
	}
}

func mediaTypeGuardMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		kind := c.Param("type")
		if kind != SubjectTypeMovie && kind != SubjectTypeTv {
			return server.JSON(c, nil, errors.BadRequest("invalid subject type"))
		}
		return next(c)
	}
}

func idGuardMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if strings.TrimSpace(c.Param("id")) == "" {
			return server.JSON(c, nil, errors.BadRequest("invalid subject id"))
		}
		return next(c)
	}
}
