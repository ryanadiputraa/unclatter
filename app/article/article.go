package article

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID          string    `json:"id" gorm:"type:varchar"`
	Title       string    `json:"title" gorm:"type:varchar;not null"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	ArticleLink string    `json:"article_link" gorm:"type:varchar;not null"`
	UserID      string    `json:"-" gorm:"type:varchar;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;not null"`
}

type NewArticleArg struct {
	Title       string
	Content     string
	ArticleLink string
	UserID      string
}

func NewArticle(arg NewArticleArg) *Article {
	return &Article{
		ID:        uuid.NewString(),
		Title:     arg.Title,
		Content:   arg.Content,
		UserID:    arg.UserID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

type ArticleService interface {
	ScrapeContent(ctx context.Context, url string) (string, error)
}
