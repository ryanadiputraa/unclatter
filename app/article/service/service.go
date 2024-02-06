package service

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/article"
)

type service struct{}

func NewService() article.ArticleService {
	return &service{}
}

func (s *service) ScrapeContent(ctx context.Context, url string) (string, error) {
	// TODO: scrape content
	return "", nil
}
