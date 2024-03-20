package article

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ryanadiputraa/unclatter/app/pagination"
)

type Article struct {
	ID          string    `json:"id" gorm:"type:varchar"`
	Title       string    `json:"title" gorm:"type:varchar;unique;not null"`
	Content     string    `json:"content,omitempty" gorm:"type:text;not null"`
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

type BookmarkPayload struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	ArticleLink string `json:"article_link" validate:"required,http_url"`
}

func NewArticle(arg NewArticleArg) *Article {
	return &Article{
		ID:          uuid.NewString(),
		Title:       arg.Title,
		Content:     arg.Content,
		ArticleLink: arg.ArticleLink,
		UserID:      arg.UserID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

type ArticleService interface {
	ScrapeContent(ctx context.Context, url string) (string, error)
	BookmarkArticle(ctx context.Context, arg BookmarkPayload, userID string) (*Article, error)
	ListBookmarkedArticles(ctx context.Context, userID string, page pagination.Pagination) ([]*Article, *pagination.Meta, error)
	GetBookmarkedArticle(ctx context.Context, userID, articleID string) (*Article, error)
	UpdateArticle(ctx context.Context, userID, articleID string, arg BookmarkPayload) (*Article, error)
	DeleteArticle(ctx context.Context, userID, articleID string) error
}

type ArticleRepository interface {
	Save(ctx context.Context, arg Article) error
	List(ctx context.Context, userID string, page pagination.Pagination) (articles []*Article, total int64, err error)
	FindByID(ctx context.Context, articleID string) (*Article, error)
	Update(ctx context.Context, arg Article) (*Article, error)
	Delete(ctx context.Context, userID, articleID string) error
}
