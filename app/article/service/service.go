package service

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/pkg/sanitizer"
	"github.com/ryanadiputraa/unclatter/pkg/scrapper"
)

type service struct {
	log       logger.Logger
	scrapper  scrapper.Scrapper
	sanitizer sanitizer.Sanitizer
}

func NewService(log logger.Logger, scrapper scrapper.Scrapper, sanitizer sanitizer.Sanitizer) article.ArticleService {
	return &service{
		log:       log,
		scrapper:  scrapper,
		sanitizer: sanitizer,
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
	return
}
