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

func (h *handler) GetUser(c *echo.Context) error {
	ctx := c.Request().Context()
	user, err := h.service.GetUser(ctx, c.Param("handle"))
	return server.JSON(c, user, err)
}

func (h *handler) GetTweet(c *echo.Context) error {
	ctx := c.Request().Context()
	tweet, err := h.service.GetTweet(ctx, c.Param("handle"), c.Param("id"))
	return server.JSON(c, tweet, err)
}

func (h *handler) ListTweets(c *echo.Context) error {
	ctx := c.Request().Context()
	cursor := strings.TrimSpace(c.QueryParam("cursor"))
	result, err := h.service.ListTweets(ctx, c.Param("handle"), cursor)
	return server.JSON(c, result, err)
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}
