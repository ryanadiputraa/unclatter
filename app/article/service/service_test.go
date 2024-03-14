package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/mocks"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/pkg/sanitizer"
	"github.com/ryanadiputraa/unclatter/pkg/scrapper"
	"github.com/ryanadiputraa/unclatter/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookmark(t *testing.T) {
	articleID := uuid.NewString()
	userID := uuid.NewString()

	cases := []struct {
		name              string
		arg               article.BookmarkPayload
		expected          *article.Article
		err               error
		mockRepoBehaviour func(mockRepo *mocks.ArticleRepository)
	}{
		{
			name: "should return bookmarked article",
			arg: article.BookmarkPayload{
				Title:       test.TestArticle.Title,
				Content:     `<div><a onblur="alert(secret)" href="http://www.google.com">Google</a><p>article content</p></div>`,
				ArticleLink: test.TestArticle.ArticleLink,
			},
			expected: &article.Article{
				ID:          articleID,
				Title:       test.TestArticle.Title,
				Content:     `<div><a href="http://www.google.com" rel="nofollow">Google</a><p>article content</p></div>`,
				ArticleLink: test.TestArticle.ArticleLink,
				UserID:      userID,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			},
			err: nil,
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository) {
				mockRepo.On("Save", context.Background(), mock.Anything).Return(nil)
			},
		},
		{
			name: "should return server error when fail to bookmark article",
			arg: article.BookmarkPayload{
				Title:       test.TestArticle.Title,
				Content:     `<div><a onblur="alert(secret)" href="http://www.google.com">Google</a><p>article content</p></div>`,
				ArticleLink: test.TestArticle.ArticleLink,
			},
			expected: &article.Article{
				ID:          articleID,
				Title:       test.TestArticle.Title,
				Content:     `<div><a href="http://www.google.com" rel="nofollow">Google</a><p>article content</p></div>`,
				ArticleLink: test.TestArticle.ArticleLink,
				UserID:      userID,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			},
			err: validation.NewError(validation.BadRequest, "title is already in use"),
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository) {
				mockRepo.On("Save", context.Background(), mock.Anything).Return(validation.NewError(validation.BadRequest, "title is already in use"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := new(mocks.ArticleRepository)
			c.mockRepoBehaviour(r)

			s := NewService(logger.NewLogger(), scrapper.NewScrapper(), sanitizer.NewSanitizer(), r)
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

// TODO: test list & update
