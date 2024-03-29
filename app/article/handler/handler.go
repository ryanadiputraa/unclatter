package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/middleware"
	"github.com/ryanadiputraa/unclatter/app/pagination"
	"github.com/ryanadiputraa/unclatter/app/validation"
	_http "github.com/ryanadiputraa/unclatter/pkg/http"
	"github.com/ryanadiputraa/unclatter/pkg/validator"
)

type handler struct {
	rw             _http.ResponseWriter
	articleService article.ArticleService
	validator      validator.Validator
}

func NewHandler(web *http.ServeMux, rw _http.ResponseWriter, articleService article.ArticleService, authMiddleware middleware.AuthMiddleware, validator validator.Validator) {
	h := &handler{
		rw:             rw,
		articleService: articleService,
		validator:      validator,
	}

	web.Handle("GET /api/articles", authMiddleware.ParseJWTToken(h.ScrapeContent()))
	web.Handle("POST /api/articles/bookmarks", authMiddleware.ParseJWTToken(h.BookmarkArticle()))
	web.Handle("GET /api/articles/bookmarks", authMiddleware.ParseJWTToken(h.ListBookmarkedArticles()))
	web.Handle("GET /api/articles/bookmarks/{id}", authMiddleware.ParseJWTToken(h.GetBookmarkedArticle()))
	web.Handle("PUT /api/articles/bookmarks/{id}", authMiddleware.ParseJWTToken(h.UpdateArticle()))
	web.Handle("DELETE /api/articles/bookmarks/{id}", authMiddleware.ParseJWTToken(h.DeleteArticle()))
}

func (h *handler) ScrapeContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if len(url) == 0 {
			h.rw.WriteErrMessage(w, http.StatusBadRequest, "missing 'url' query param")
			return
		}

		content, err := h.articleService.ScrapeContent(r.Context(), url)
		if err != nil {
			h.rw.WriteErrMessage(w, http.StatusBadRequest, "fail to get page content")
			return
		}

		h.rw.WriteResponseData(w, http.StatusOK, content)
	}
}

func (h *handler) BookmarkArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ac := r.Context().(*middleware.AuthContext)
		var payload article.BookmarkPayload

		json.NewDecoder(r.Body).Decode(&payload)
		if err, errMap := h.validator.Validate(payload); err != nil {
			h.rw.WriteErrDetails(w, http.StatusBadRequest, "invalid params", errMap)
			return
		}

		bookmarked, err := h.articleService.BookmarkArticle(ac.Context, payload, ac.UserID)
		if err != nil {
			if vErr, ok := err.(*validation.Error); ok {
				h.rw.WriteErrMessage(w, validation.HttpErrMap[vErr.Err], vErr.Message)
				return
			}
			h.rw.WriteErrMessage(w, http.StatusInternalServerError, "internal server error")
			return
		}

		h.rw.WriteResponseData(w, http.StatusCreated, bookmarked)
	}
}

func (h *handler) ListBookmarkedArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ac := r.Context().(*middleware.AuthContext)
		query := r.URL.Query()
		page := query.Get("page")
		size := query.Get("size")

		pagination, errMap, err := pagination.ValidateParam(page, size)
		if err != nil {
			h.rw.WriteErrDetails(w, http.StatusBadRequest, "invalid params", errMap)
			return
		}

		articles, meta, err := h.articleService.ListBookmarkedArticles(ac.Context, ac.UserID, *pagination)
		if err != nil {
			if vErr, ok := err.(*validation.Error); ok {
				h.rw.WriteErrMessage(w, validation.HttpErrMap[vErr.Err], vErr.Message)
				return
			}
			h.rw.WriteErrMessage(w, http.StatusInternalServerError, "internal server error")
			return
		}

		h.rw.WriteResponseDataWithPagination(w, http.StatusOK, articles, *meta)
	}
}

func (h *handler) GetBookmarkedArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ac := r.Context().(*middleware.AuthContext)
		articleID := r.PathValue("id")

		article, err := h.articleService.GetBookmarkedArticle(ac.Context, ac.UserID, articleID)
		if err != nil {
			if vErr, ok := err.(*validation.Error); ok {
				h.rw.WriteErrMessage(w, validation.HttpErrMap[vErr.Err], vErr.Message)
				return
			}
			h.rw.WriteErrMessage(w, http.StatusInternalServerError, "internal server error")
			return
		}

		h.rw.WriteResponseData(w, http.StatusOK, article)
	}
}

func (h *handler) UpdateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ac := r.Context().(*middleware.AuthContext)
		id := r.PathValue("id")
		var payload article.BookmarkPayload

		json.NewDecoder(r.Body).Decode(&payload)
		if err, errMap := h.validator.Validate(payload); err != nil {
			h.rw.WriteErrDetails(w, http.StatusBadRequest, "invalid params", errMap)
			return
		}

		article, err := h.articleService.UpdateArticle(ac, ac.UserID, id, payload)
		if err != nil {
			if vErr, ok := err.(*validation.Error); ok {
				h.rw.WriteErrMessage(w, validation.HttpErrMap[vErr.Err], vErr.Message)
				return
			}
			h.rw.WriteErrMessage(w, http.StatusInternalServerError, "internal server error")
			return
		}

		h.rw.WriteResponseData(w, http.StatusOK, article)
	}
}

func (h *handler) DeleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ac := r.Context().(*middleware.AuthContext)
		id := r.PathValue("id")

		err := h.articleService.DeleteArticle(ac, ac.UserID, id)
		if err != nil {
			if vErr, ok := err.(*validation.Error); ok {
				h.rw.WriteErrMessage(w, validation.HttpErrMap[vErr.Err], vErr.Message)
				return
			}
			h.rw.WriteErrMessage(w, http.StatusInternalServerError, "internal server error")
			return
		}

		h.rw.WriteResponseData(w, http.StatusOK, nil)
	}
}
