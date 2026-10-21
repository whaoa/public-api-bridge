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

func (h *handler) SearchMedia(c *echo.Context) error {
	keyword := strings.TrimSpace(c.QueryParam("keyword"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	res, err := h.service.SearchMedia(c.Request().Context(), keyword, c.QueryParam("lang"), page)
	return server.JSON(c, res, err)
}

func (h *handler) GetMovie(c *echo.Context) error {
	res, err := h.service.GetMovie(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

func (h *handler) GetTV(c *echo.Context) error {
	res, err := h.service.GetTV(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

func (h *handler) GetSeason(c *echo.Context) error {
	res, err := h.service.GetSeason(c.Request().Context(), c.Param("id"), c.Param("season"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

func (h *handler) NewMediaCastsRoute(kind string) echo.HandlerFunc {
	return func(c *echo.Context) error {
		res, err := h.service.GetMediaCasts(c.Request().Context(), kind, c.Param("id"), c.QueryParam("lang"))
		return server.JSON(c, res, err)
	}
}

func (h *handler) NewMediaImagesRoute(kind string) echo.HandlerFunc {
	return func(c *echo.Context) error {
		res, err := h.service.GetMediaImages(c.Request().Context(), kind, c.Param("id"), c.QueryParam("lang"))
		return server.JSON(c, res, err)
	}
}

func (h *handler) SearchPerson(c *echo.Context) error {
	keyword := strings.TrimSpace(c.QueryParam("keyword"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	res, err := h.service.SearchPerson(c.Request().Context(), keyword, c.QueryParam("lang"), page)
	return server.JSON(c, res, err)
}

func (h *handler) GetPerson(c *echo.Context) error {
	res, err := h.service.GetPerson(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

func (h *handler) GetPersonProfiles(c *echo.Context) error {
	res, err := h.service.GetPersonProfiles(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

func (h *handler) GetPersonCredits(c *echo.Context) error {
	res, err := h.service.GetPersonCredits(c.Request().Context(), c.Param("id"), c.QueryParam("lang"))
	return server.JSON(c, res, err)
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}
