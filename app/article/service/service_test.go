package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/pkg/sanitizer"
	"github.com/ryanadiputraa/unclatter/pkg/scrapper"
	"github.com/stretchr/testify/assert"
)

func TestBookmark(t *testing.T) {
	articleID := uuid.NewString()
	userID := uuid.NewString()

	cases := []struct {
		name     string
		arg      article.BookmarkPayload
		expected *article.Article
		err      error
	}{
		{
			name: "should return bookmarked article",
			arg: article.BookmarkPayload{
				Title:       "Content Title",
				Content:     `<div><a onblur="alert(secret)" href="http://www.google.com">Google</a><p>article content</p></div>`,
				ArticleLink: "https://unclatter.com",
			},
			expected: &article.Article{
				ID:          articleID,
				Title:       "Content Title",
				Content:     `<div><a href="http://www.google.com" rel="nofollow">Google</a><p>article content</p></div>`,
				ArticleLink: "https://unclatter.com",
				UserID:      userID,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			},
			err: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := NewService(logger.NewLogger(), scrapper.NewScrapper(), sanitizer.NewSanitizer())
			article, err := s.BookmarkArticle(context.Background(), c.arg, userID)

			assert.Equal(t, c.err, err)
			if err != nil {
				return
			}

			assert.NotNil(t, article.ID)
			assert.Equal(t, c.expected.Title, article.Title)
			assert.Equal(t, c.expected.Content, article.Content)
			assert.Equal(t, c.expected.ArticleLink, article.ArticleLink)
			assert.Equal(t, c.expected.UserID, article.UserID)
			assert.NotEmpty(t, article.CreatedAt)
			assert.NotEmpty(t, article.UpdatedAt)
		})
	}
}
