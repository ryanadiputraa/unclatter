package article

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewArticle(t *testing.T) {
	uuid := uuid.NewString()

	cases := []struct {
		name     string
		arg      NewArticleArg
		expected *Article
		error
	}{
		{
			name: "should return a valid user",
			arg: NewArticleArg{
				Title:       "Title",
				Content:     "<p>Sample Content Body</p>",
				ArticleLink: "https://unclatter.com",
				UserID:      uuid,
			},
			expected: &Article{
				ID:          uuid,
				Title:       "Title",
				Content:     "<p>Sample Content Body</p>",
				ArticleLink: "https://unclatter.com",
				UserID:      uuid,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			},
			error: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			user := NewArticle(c.arg)

			assert.NotEmpty(t, user.ID)
			assert.Equal(t, c.expected.Title, user.Title)
			assert.Equal(t, c.expected.Content, user.Content)
			assert.Equal(t, c.expected.ArticleLink, user.ArticleLink)
			assert.Equal(t, c.expected.UserID, user.UserID)
			assert.NotEmpty(t, user.CreatedAt)
			assert.NotEmpty(t, user.UpdatedAt)
		})
	}
}
