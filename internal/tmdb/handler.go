package tmdb

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/whaoa/public-api-bridge/internal/server"
)

type handler struct {
	service *Service
}

func (h *handler) RegisterRoutes(router server.Router) {
	router.GET("/tmdb/media", h.SearchMedia)

	router.GET("/tmdb/movie/:id", h.GetMovie)
	router.GET("/tmdb/movie/:id/casts", h.NewMediaCastsRoute(TypeMovie))
	router.GET("/tmdb/movie/:id/images", h.NewMediaImagesRoute(TypeMovie))

	router.GET("/tmdb/tv/:id", h.GetTV)
	router.GET("/tmdb/tv/:id/seasons/:season", h.GetSeason)
	router.GET("/tmdb/tv/:id/casts", h.NewMediaCastsRoute(TypeTV))
	router.GET("/tmdb/tv/:id/images", h.NewMediaImagesRoute(TypeTV))

	router.GET("/tmdb/person", h.SearchPerson)
	router.GET("/tmdb/person/:id", h.GetPerson)
	router.GET("/tmdb/person/:id/profiles", h.GetPersonProfiles)
	router.GET("/tmdb/person/:id/credits", h.GetPersonCredits)
}

// SearchMedia searches movies and TV shows by keyword.
// @Summary Search media
// @Tags TMDB
// @Param keyword query string true "Search keyword"
// @Param lang query string false "Language"
// @Param page query int false "Page number"
// @Success 200 {object} MediaSlice
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/media [get]
func (h *handler) SearchMedia(c *echo.Context) error {
	keyword := strings.TrimSpace(c.QueryParam("keyword"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	res, err := h.service.SearchMedia(c.Request().Context(), keyword, c.QueryParam("lang"), page)
	return server.JSON(c, res, err)
}

// GetMovie returns a movie by TMDB ID.
// @Summary Get a movie
// @Tags TMDB
// @Param id path string true "TMDB movie ID"
// @Param lang query string false "Language"
// @Success 200 {object} Movie
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/movie/{id} [get]
func (h *handler) GetMovie(c *echo.Context) error {
	res, err := h.service.GetMovie(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

// GetTV returns a TV show by TMDB ID.
// @Summary Get a TV show
// @Tags TMDB
// @Param id path string true "TMDB TV show ID"
// @Param lang query string false "Language"
// @Success 200 {object} TV
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/tv/{id} [get]
func (h *handler) GetTV(c *echo.Context) error {
	res, err := h.service.GetTV(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

// GetSeason returns a TV season by show ID and season number.
// @Summary Get a TV season
// @Tags TMDB
// @Param id path string true "TMDB TV show ID"
// @Param season path string true "Season number"
// @Param lang query string false "Language"
// @Success 200 {object} Season
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/tv/{id}/seasons/{season} [get]
func (h *handler) GetSeason(c *echo.Context) error {
	res, err := h.service.GetSeason(c.Request().Context(), c.Param("id"), c.Param("season"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

// NewMediaCastsRoute returns the cast and crew for a movie or TV show.
// @Summary Get media cast
// @Tags TMDB
// @Param id path string true "TMDB media ID"
// @Param lang query string false "Language"
// @Success 200 {object} MediaCasts
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/movie/{id}/casts [get]
// @Router /tmdb/tv/{id}/casts [get]
func (h *handler) NewMediaCastsRoute(kind string) echo.HandlerFunc {
	return func(c *echo.Context) error {
		res, err := h.service.GetMediaCasts(c.Request().Context(), kind, c.Param("id"), c.QueryParam("lang"))
		return server.JSON(c, res, err)
	}
}

// NewMediaImagesRoute returns images for a movie or TV show.
// @Summary Get media images
// @Tags TMDB
// @Param id path string true "TMDB media ID"
// @Param lang query string false "Language"
// @Success 200 {object} MediaImages
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/movie/{id}/images [get]
// @Router /tmdb/tv/{id}/images [get]
func (h *handler) NewMediaImagesRoute(kind string) echo.HandlerFunc {
	return func(c *echo.Context) error {
		res, err := h.service.GetMediaImages(c.Request().Context(), kind, c.Param("id"), c.QueryParam("lang"))
		return server.JSON(c, res, err)
	}
}

// SearchPerson searches people by keyword.
// @Summary Search people
// @Tags TMDB
// @Param keyword query string true "Search keyword"
// @Param lang query string false "Language"
// @Param page query int false "Page number"
// @Success 200 {object} EmbedPersonSlice
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/person [get]
func (h *handler) SearchPerson(c *echo.Context) error {
	keyword := strings.TrimSpace(c.QueryParam("keyword"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	res, err := h.service.SearchPerson(c.Request().Context(), keyword, c.QueryParam("lang"), page)
	return server.JSON(c, res, err)
}

// GetPerson returns a person by TMDB ID.
// @Summary Get a person
// @Tags TMDB
// @Param id path string true "TMDB person ID"
// @Param lang query string false "Language"
// @Success 200 {object} Person
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/person/{id} [get]
func (h *handler) GetPerson(c *echo.Context) error {
	res, err := h.service.GetPerson(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

// GetPersonProfiles returns profile images for a person.
// @Summary Get person profiles
// @Tags TMDB
// @Param id path string true "TMDB person ID"
// @Param lang query string false "Language"
// @Success 200 {array} Image
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/person/{id}/profiles [get]
func (h *handler) GetPersonProfiles(c *echo.Context) error {
	res, err := h.service.GetPersonProfiles(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

// GetPersonCredits returns a person's movie and TV credits.
// @Summary Get person credits
// @Tags TMDB
// @Param id path string true "TMDB person ID"
// @Param lang query string false "Language"
// @Success 200 {object} PersonCredits
// @Failure 500 {object} errors.ReasonError
// @Router /tmdb/person/{id}/credits [get]
func (h *handler) GetPersonCredits(c *echo.Context) error {
	res, err := h.service.GetPersonCredits(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}
