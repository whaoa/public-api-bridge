package bangumi

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
	router.GET("/bangumi/subjects", h.SearchAnime)
	router.GET("/bangumi/subjects/:id", h.GetAnime)
	router.GET("/bangumi/subjects/:id/episodes", h.ListEpisodes)
	router.GET("/bangumi/subjects/:id/characters", h.GetCharacters)
	router.GET("/bangumi/subjects/:id/staffs", h.GetStaffs)
	router.GET("/bangumi/person/:id", h.GetPerson)
	router.GET("/bangumi/person/:id/credits", h.GetPersonCredits)
}

// SearchAnime searches anime by keyword.
// @Summary Search anime
// @Tags Bangumi
// @Param keyword query string true "Search keyword"
// @Param offset query int false "Result offset"
// @Param count query int false "Result count"
// @Success 200 {object} SubjectSlice
// @Failure 500 {object} errors.ReasonError
// @Router /bangumi/subjects [get]
func (h *handler) SearchAnime(c *echo.Context) error {
	keyword := strings.TrimSpace(c.QueryParam("keyword"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	count, _ := strconv.Atoi(c.QueryParam("count"))
	result, err := h.service.SearchAnime(c.Request().Context(), keyword, offset, count)
	return server.JSON(c, result, err)
}

// GetAnime returns a bangumi subject.
// @Summary Get a anime subject
// @Tags Bangumi
// @Param id path string true "Subject ID"
// @Success 200 {object} Subject
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /bangumi/subjects/{id} [get]
func (h *handler) GetAnime(c *echo.Context) error {
	anime, err := h.service.GetAnime(c.Request().Context(), c.Param("id"))
	return server.JSON(c, anime, err)
}

// ListEpisodes returns episodes for a bangumi subject.
// @Summary get anime episodes
// @Tags Bangumi
// @Param id path string true "Subject ID"
// @Param offset query int false "Result offset"
// @Param count query int false "Result count"
// @Success 200 {object} EpisodeSlice
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /bangumi/subjects/{id}/episodes [get]
func (h *handler) ListEpisodes(c *echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	count, _ := strconv.Atoi(c.QueryParam("count"))
	result, err := h.service.ListEpisodes(c.Request().Context(), c.Param("id"), offset, count)
	return server.JSON(c, result, err)
}

// GetCharacters returns characters for a bangumi subject.
// @Summary get anime characters
// @Tags Bangumi
// @Param id path string true "Subject ID"
// @Success 200 {array} Character
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /bangumi/subjects/{id}/characters [get]
func (h *handler) GetCharacters(c *echo.Context) error {
	result, err := h.service.GetCharacters(c.Request().Context(), c.Param("id"))
	return server.JSON(c, result, err)
}

// GetStaffs returns staffs for a bangumi subject.
// @Summary get anime staffs
// @Tags Bangumi
// @Param id path string true "Subject ID"
// @Success 200 {array} Staff
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /bangumi/subjects/{id}/staffs [get]
func (h *handler) GetStaffs(c *echo.Context) error {
	result, err := h.service.GetStaffs(c.Request().Context(), c.Param("id"))
	return server.JSON(c, result, err)
}

// GetPerson returns a bangumi person.
// @Summary get a bangumi person
// @Tags Bangumi
// @Param id path string true "Person ID"
// @Success 200 {object} Person
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /bangumi/{id}/person [get]
func (h *handler) GetPerson(c *echo.Context) error {
	result, err := h.service.GetPerson(c.Request().Context(), c.Param("id"))
	return server.JSON(c, result, err)
}

// GetPersonCredits returns credits for a bangumi person.
// @Summary list person credits
// @Tags Bangumi
// @Param id path string true "Person ID"
// @Success 200 {array} PersonCredit
// @Failure 400 {object} errors.ReasonError
// @Failure 500 {object} errors.ReasonError
// @Router /bangumi/{id}/credits [get]
func (h *handler) GetPersonCredits(c *echo.Context) error {
	result, err := h.service.GetPersonCredits(c.Request().Context(), c.Param("id"))
	return server.JSON(c, result, err)
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}
