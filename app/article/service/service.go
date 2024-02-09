package service

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/pkg/scrapper"
)

type service struct {
	log      logger.Logger
	scrapper scrapper.Scrapper
}

func NewService(log logger.Logger, scrapper scrapper.Scrapper) article.ArticleService {
	return &service{
		log:      log,
		scrapper: scrapper,
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
