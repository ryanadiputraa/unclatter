package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/middleware"
)

type handler struct {
	articleService article.ArticleService
}

func NewHandler(r *echo.Group, articleService article.ArticleService, authMiddleware middleware.AuthMiddleware) {
	h := &handler{
		articleService: articleService,
	}

	r.GET("", h.ScrapeContent(), authMiddleware.ParseJWTToken)
}

func (h *handler) ScrapeContent() echo.HandlerFunc {
	return func(c echo.Context) error {
		url := c.QueryParam("url")
		if len(url) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "missing url query param",
			})
		}

		content, err := h.articleService.ScrapeContent(c.Request().Context(), url)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "fail to get page content",
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"data": content,
		})
	}
}
