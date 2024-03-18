package service

import (
	"context"
	"math"
	"time"

	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/pagination"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/pkg/sanitizer"
	"github.com/ryanadiputraa/unclatter/pkg/scrapper"
)

type service struct {
	log        logger.Logger
	scrapper   scrapper.Scrapper
	sanitizer  sanitizer.Sanitizer
	repository article.ArticleRepository
}

func NewService(log logger.Logger, scrapper scrapper.Scrapper, sanitizer sanitizer.Sanitizer, repository article.ArticleRepository) article.ArticleService {
	return &service{
		log:        log,
		scrapper:   scrapper,
		sanitizer:  sanitizer,
		repository: repository,
	}
}

func (s *service) ScrapeContent(ctx context.Context, url string) (content string, err error) {
	content, err = s.scrapper.ScrapeTextContent(url)
	if err != nil {
		s.log.Warn("article service: fail to scrape page", err)
		return
	}
	return
}

func (s *service) BookmarkArticle(ctx context.Context, arg article.BookmarkPayload, userID string) (bookmarked *article.Article, err error) {
	bookmarked = article.NewArticle(article.NewArticleArg{
		Title:       arg.Title,
		Content:     s.sanitizer.Sanitize(arg.Content),
		ArticleLink: arg.ArticleLink,
		UserID:      userID,
	})

	if err = s.repository.Save(ctx, *bookmarked); err != nil {
		return
	}
	return
}

func (s *service) ListBookmarkedArticles(ctx context.Context, userID string, page pagination.Pagination) (articles []*article.Article, meta *pagination.Meta, err error) {
	articles, total, err := s.repository.List(ctx, userID, page)
	if err != nil {
		s.log.Error("article service: fail to fetch user's bookmarked articles", err)
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(page.Limit)))
	}

	meta = &pagination.Meta{
		CurrentPage: page.Offset/page.Limit + 1,
		TotalPages:  totalPages,
		Size:        page.Limit,
		TotalData:   total,
	}

	return
}

func (s *service) UpdateArticle(ctx context.Context, userID, articleID string, arg article.BookmarkPayload) (updated *article.Article, err error) {
	update := article.Article{
		ID:          articleID,
		Title:       arg.Title,
		Content:     arg.Content,
		ArticleLink: arg.ArticleLink,
		UserID:      userID,
		UpdatedAt:   time.Now().UTC(),
	}
	updated, err = s.repository.Update(ctx, update)
	if err != nil {
		s.log.Warn("article service: fail to update bookmarked article", err)
	}
	return
}
