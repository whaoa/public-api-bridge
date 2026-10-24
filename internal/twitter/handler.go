package twitter

import (
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/whaoa/public-api-bridge/internal/server"
)

type handler struct {
	service *Service
}

func (h *handler) RegisterRoutes(router server.Router) {
	router.GET("/twitter/:handle", h.GetUser)
	router.GET("/twitter/:handle/statuses", h.ListTweets)
	router.GET("/twitter/:handle/statuses/:id", h.GetTweet)

	router.GET("/twitter/:handle/status/:id", h.GetTweet)
}

// GetUser returns a public user profile by handle.
// @Summary Get a Twitter user
// @Tags Twitter
// @Param handle path string true "User handle"
// @Success 200 {object} User
// @Failure 500 {object} errors.ReasonError
// @Router /twitter/{handle} [get]
func (h *handler) GetUser(c *echo.Context) error {
	ctx := c.Request().Context()
	user, err := h.service.GetUser(ctx, c.Param("handle"))
	return server.JSON(c, user, err)
}

// GetTweet returns a public tweet by handle and tweet ID.
// @Summary Get a tweet
// @Tags Twitter
// @Param handle path string true "User handle"
// @Param id path string true "Tweet ID"
// @Success 200 {object} Tweet
// @Failure 500 {object} errors.ReasonError
// @Router /twitter/{handle}/statuses/{id} [get]
// @Router /twitter/{handle}/status/{id} [get]
func (h *handler) GetTweet(c *echo.Context) error {
	ctx := c.Request().Context()
	tweet, err := h.service.GetTweet(ctx, c.Param("handle"), c.Param("id"))
	return server.JSON(c, tweet, err)
}

// ListTweets returns a page of tweets for a user.
// @Summary List user tweets
// @Tags Twitter
// @Param handle path string true "User handle"
// @Param cursor query string false "Pagination cursor"
// @Success 200 {object} TweetSlice
// @Failure 500 {object} errors.ReasonError
// @Router /twitter/{handle}/statuses [get]
func (h *handler) ListTweets(c *echo.Context) error {
	ctx := c.Request().Context()
	cursor := strings.TrimSpace(c.QueryParam("cursor"))
	result, err := h.service.ListTweets(ctx, c.Param("handle"), cursor)
	return server.JSON(c, result, err)
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}
