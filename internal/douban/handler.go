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

func (h *handler) SearchMovie(c *echo.Context) error {
	keyword := strings.TrimSpace(c.QueryParam("keyword"))
	result, err := h.service.SearchMovie(c.Request().Context(), keyword)
	return server.JSON(c, result, err)
}

func (h *handler) GetMovie(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetMovie(c.Request().Context(), c.Param("type"), id)
	return server.JSON(c, result, err)
}

func (h *handler) ListMoviePhotos(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	count, _ := strconv.Atoi(c.QueryParam("count"))
	result, err := h.service.ListMoviePhotos(c.Request().Context(), c.Param("type"), id, offset, count)
	return server.JSON(c, result, err)
}

func (h *handler) GetMovieVideos(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetMovieVideos(c.Request().Context(), c.Param("type"), id)
	return server.JSON(c, result, err)
}

func (h *handler) GetMovieCasts(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetMovieCasts(c.Request().Context(), c.Param("type"), id)
	return server.JSON(c, result, err)
}

func (h *handler) GetMovieRecommendations(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetMovieRecommendations(c.Request().Context(), c.Param("type"), id)
	return server.JSON(c, result, err)
}

func (h *handler) GetPerson(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.service.GetPerson(c.Request().Context(), id)
	return server.JSON(c, result, err)
}

func (h *handler) ListPersonPhotos(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	count, _ := strconv.Atoi(c.QueryParam("count"))
	result, err := h.service.ListPersonPhotos(c.Request().Context(), id, offset, count)
	return server.JSON(c, result, err)
}

func (h *handler) ListPersonMovieCredits(c *echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	count, _ := strconv.Atoi(c.QueryParam("count"))
	result, err := h.service.ListPersonMovieCredits(c.Request().Context(), id, offset, count)
	return server.JSON(c, result, err)
}

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
