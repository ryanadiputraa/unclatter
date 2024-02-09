package article

import "context"

type ArticleService interface {
	ScrapeContent(ctx context.Context, url string) (string, error)
}
