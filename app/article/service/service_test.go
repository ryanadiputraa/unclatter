package service

import (
	"context"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/mocks"
	"github.com/ryanadiputraa/unclatter/app/pagination"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/pkg/sanitizer"
	"github.com/ryanadiputraa/unclatter/pkg/scrapper"
	"github.com/ryanadiputraa/unclatter/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestScrapeContent(t *testing.T) {
	cases := []struct {
		name                  string
		url                   string
		expected              string
		err                   error
		mockScrapperBehaviour func(mockScrapper *mocks.Scrapper, url string)
	}{
		{
			name:     "should return scrapped article content",
			url:      test.TestArticle.ArticleLink,
			expected: test.TestArticle.Content,
			err:      nil,
			mockScrapperBehaviour: func(mockScrapper *mocks.Scrapper, url string) {
				mockScrapper.On("ScrapeTextContent", url).Return(test.TestArticle.Content, nil)
			},
		},
		{
			name:     "should return err when fail to scrape article content",
			url:      test.TestArticle.ArticleLink,
			expected: "",
			err:      validation.NewError(validation.BadRequest, "fail to scrape any article content"),
			mockScrapperBehaviour: func(mockScrapper *mocks.Scrapper, url string) {
				mockScrapper.On("ScrapeTextContent", url).Return("", nil)
			},
		},
		{
			name:     "should return err when fail to scrape article content",
			url:      test.TestArticle.ArticleLink,
			expected: "",
			err:      colly.ErrForbiddenURL,
			mockScrapperBehaviour: func(mockScrapper *mocks.Scrapper, url string) {
				mockScrapper.On("ScrapeTextContent", url).Return("", colly.ErrForbiddenURL)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			scrapperPkg := new(mocks.Scrapper)
			c.mockScrapperBehaviour(scrapperPkg, c.url)

			r := new(mocks.ArticleRepository)
			s := NewService(logger.NewLogger(), scrapperPkg, sanitizer.NewSanitizer(), r)
			content, err := s.ScrapeContent(context.Background(), c.url)

			assert.Equal(t, c.err, err)
			assert.Equal(t, c.expected, content)
		})
	}
}

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

func TestListArticle(t *testing.T) {
	cases := []struct {
		name              string
		userID            string
		page              pagination.Pagination
		expected          []*article.Article
		meta              *pagination.Meta
		err               error
		mockRepoBehaviour func(mockRepo *mocks.ArticleRepository, userID string, page pagination.Pagination)
	}{
		{
			name:   "should return list of user's bookmarked articles",
			userID: test.TestUser.ID,
			page: pagination.Pagination{
				Limit:  2,
				Offset: 0,
			},
			expected: []*article.Article{
				test.TestArticle,
				test.TestArticle2,
				test.TestArticle3,
			},
			meta: &pagination.Meta{
				CurrentPage: 1,
				TotalPages:  2,
				Size:        2,
				TotalData:   3,
			},
			err: nil,
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository, userID string, page pagination.Pagination) {
				mockRepo.On("List", context.Background(), userID, page).
					Return(
						[]*article.Article{test.TestArticle, test.TestArticle2, test.TestArticle3},
						int64(3),
						nil,
					)
			},
		},
		{
			name:   "should return empty list of user's bookmarked articles",
			userID: test.TestUser.ID,
			page: pagination.Pagination{
				Limit:  2,
				Offset: 0,
			},
			expected: []*article.Article{},
			meta: &pagination.Meta{
				CurrentPage: 1,
				TotalPages:  0,
				Size:        2,
				TotalData:   0,
			},
			err: nil,
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository, userID string, page pagination.Pagination) {
				mockRepo.On("List", context.Background(), userID, page).
					Return(
						[]*article.Article{},
						int64(0),
						nil,
					)
			},
		},
		{
			name:   "should return err when fail to fetch user's bookmarked articles",
			userID: test.TestUser.ID,
			page: pagination.Pagination{
				Limit:  2,
				Offset: 0,
			},
			expected: []*article.Article{},
			meta:     nil,
			err:      gorm.ErrInvalidDB,
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository, userID string, page pagination.Pagination) {
				mockRepo.On("List", context.Background(), userID, page).Return(nil, int64(0), gorm.ErrInvalidDB)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := new(mocks.ArticleRepository)
			c.mockRepoBehaviour(r, c.userID, c.page)

			s := NewService(logger.NewLogger(), scrapper.NewScrapper(), sanitizer.NewSanitizer(), r)
			articles, meta, err := s.ListBookmarkedArticles(context.Background(), c.userID, c.page)

			assert.Equal(t, c.err, err)
			if err != nil {
				return
			}

			assert.Equal(t, len(c.expected), len(articles))
			for i, article := range articles {
				assert.Equal(t, c.expected[i].ID, article.ID)
				assert.Equal(t, c.expected[i].Title, article.Title)
				assert.Equal(t, c.expected[i].Content, article.Content)
				assert.Equal(t, c.expected[i].ArticleLink, article.ArticleLink)
				assert.Equal(t, c.expected[i].UserID, article.UserID)
				assert.NotEmpty(t, article.CreatedAt)
				assert.NotEmpty(t, article.UpdatedAt)
			}

			assert.Equal(t, c.meta.CurrentPage, meta.CurrentPage)
			assert.Equal(t, c.meta.TotalPages, meta.TotalPages)
			assert.Equal(t, c.meta.Size, meta.Size)
			assert.Equal(t, c.meta.TotalData, meta.TotalData)
		})
	}
}

func TestGetArticle(t *testing.T) {
	cases := []struct {
		name              string
		userID            string
		articleID         string
		expected          *article.Article
		err               error
		mockRepoBehaviour func(mockRepo *mocks.ArticleRepository, articleID string)
	}{
		{
			name:      "should return bookmarked article with given id",
			userID:    test.TestArticle.UserID,
			articleID: test.TestArticle.ID,
			expected:  test.TestArticle,
			err:       nil,
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository, articleID string) {
				mockRepo.On("FindByID", context.Background(), articleID).
					Return(test.TestArticle, nil)
			},
		},
		{
			name:      "should return err when accessing other user's bookmarked article",
			userID:    uuid.NewString(),
			articleID: test.TestArticle.ID,
			expected:  test.TestArticle,
			err:       validation.NewError(validation.Forbidden, "forbidden access"),
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository, articleID string) {
				mockRepo.On("FindByID", context.Background(), articleID).
					Return(test.TestArticle, nil)
			},
		},
		{
			name:      "should return err when fail to fetch bookmarked article with given id",
			userID:    test.TestUser.ID,
			articleID: uuid.NewString(),
			expected:  nil,
			err:       validation.NewError(validation.NotFound, "no article found with given id"),
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository, articleID string) {
				mockRepo.On("FindByID", context.Background(), articleID).
					Return(nil, validation.NewError(validation.NotFound, "no article found with given id"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := new(mocks.ArticleRepository)
			c.mockRepoBehaviour(r, c.articleID)

			s := NewService(logger.NewLogger(), scrapper.NewScrapper(), sanitizer.NewSanitizer(), r)
			article, err := s.GetBookmarkedArticle(context.Background(), c.userID, c.articleID)

			assert.Equal(t, c.err, err)
			if err != nil {
				assert.Equal(t, c.expected, c.expected)
				return
			}

			assert.Equal(t, c.expected.ID, article.ID)
			assert.Equal(t, c.expected.Title, article.Title)
			assert.Equal(t, c.expected.Content, article.Content)
			assert.Equal(t, c.expected.ArticleLink, article.ArticleLink)
			assert.Equal(t, c.expected.UserID, article.UserID)
			assert.NotEmpty(t, article.CreatedAt)
			assert.NotEmpty(t, article.UpdatedAt)
		})
	}
}

func TestUpdateArticle(t *testing.T) {
	updatePayload := article.BookmarkPayload{
		Title:       "Updated Title",
		Content:     "<p>Updated Content</p>",
		ArticleLink: "https://newlink.com",
	}
	updatedTime := time.Now().UTC()

	cases := []struct {
		name              string
		userID            string
		articleID         string
		arg               article.BookmarkPayload
		expected          *article.Article
		err               error
		mockRepoBehaviour func(mockRepo *mocks.ArticleRepository)
	}{
		{
			name:      "should return updated article",
			userID:    test.TestArticle.UserID,
			articleID: test.TestArticle.ID,
			arg:       updatePayload,
			expected: &article.Article{
				ID:          test.TestArticle.ID,
				Title:       updatePayload.Title,
				Content:     updatePayload.Content,
				ArticleLink: updatePayload.ArticleLink,
				UserID:      test.TestArticle.UserID,
				CreatedAt:   test.TestArticle.CreatedAt,
				UpdatedAt:   updatedTime,
			},
			err: nil,
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository) {
				mockRepo.On("Update", context.Background(), mock.Anything).
					Return(
						&article.Article{
							ID:          test.TestArticle.ID,
							Title:       updatePayload.Title,
							Content:     updatePayload.Content,
							ArticleLink: updatePayload.ArticleLink,
							UserID:      test.TestArticle.UserID,
							CreatedAt:   test.TestArticle.CreatedAt,
							UpdatedAt:   updatedTime,
						},
						nil,
					)
			},
		},
		{
			name:      "should return err when fail to update article",
			userID:    test.TestArticle.UserID,
			articleID: test.TestArticle.ID,
			arg:       updatePayload,
			expected:  nil,
			err:       validation.NewError(validation.Forbidden, "forbidden access"),
			mockRepoBehaviour: func(mockRepo *mocks.ArticleRepository) {
				mockRepo.On("Update", context.Background(), mock.Anything).
					Return(nil, validation.NewError(validation.Forbidden, "forbidden access"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := new(mocks.ArticleRepository)
			c.mockRepoBehaviour(r)

			s := NewService(logger.NewLogger(), scrapper.NewScrapper(), sanitizer.NewSanitizer(), r)
			article, err := s.UpdateArticle(context.Background(), c.userID, c.articleID, c.arg)

			assert.Equal(t, c.err, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.expected.ID, article.ID)
			assert.Equal(t, c.expected.Title, article.Title)
			assert.Equal(t, c.expected.Content, article.Content)
			assert.Equal(t, c.expected.ArticleLink, article.ArticleLink)
			assert.Equal(t, c.expected.UserID, article.UserID)
			assert.NotEmpty(t, article.CreatedAt)
			assert.NotEmpty(t, article.UpdatedAt)
		})
	}
}
