package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/middleware"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/pkg/validator"
)

type handler struct {
	articleService article.ArticleService
	validator      validator.Validator
}

func NewHandler(r *echo.Group, articleService article.ArticleService, authMiddleware middleware.AuthMiddleware, validator validator.Validator) {
	h := &handler{
		articleService: articleService,
		validator:      validator,
	}

	r.GET("", h.ScrapeContent(), authMiddleware.ParseJWTToken)
	r.POST("/bookmarks", h.BookmarkArticle(), authMiddleware.ParseJWTToken)
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

func (h *handler) BookmarkArticle() echo.HandlerFunc {
	return func(c echo.Context) error {
		rc := c.(*middleware.RequestContext)

		var payload article.BookmarkPayload
		c.Bind(&payload)
		if err, errMap := h.validator.Validate(payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "invalid params",
				"error":   errMap,
			})
		}

		bookmarked, err := h.articleService.BookmarkArticle(c.Request().Context(), payload, rc.UserID)
		if err != nil {
			if vErr, ok := err.(*validation.Error); ok {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": vErr.Message,
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "internal server error",
			})
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"data": bookmarked,
		})
	}
}
