package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/pagination"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSave(t *testing.T) {
	gormDB, db, mock := test.NewMockDB(t)
	defer db.Close()

	r := NewRepository(gormDB)
	expectedExec := "INSERT INTO \"articles\""

	cases := []struct {
		name          string
		mockBehaviour func(mock sqlmock.Sqlmock)
		err           error
	}{
		{
			name: "should insert new bookmarked article",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(expectedExec).
					WithArgs(test.TestArticle.ID, test.TestArticle.Title, test.TestArticle.Content, test.TestArticle.ArticleLink,
						test.TestArticle.UserID, test.TestArticle.CreatedAt, test.TestArticle.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			err: nil,
		},
		{
			name: "should return error when title already used",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(expectedExec).
					WithArgs(test.TestArticle.ID, test.TestArticle.Title, test.TestArticle.Content, test.TestArticle.ArticleLink,
						test.TestArticle.UserID, test.TestArticle.CreatedAt, test.TestArticle.UpdatedAt).
					WillReturnError(gorm.ErrDuplicatedKey)
				mock.ExpectRollback()
			},
			err: validation.NewError(validation.BadRequest, "title is already in use"),
		},
		{
			name: "should return error when fail to insert new bookmarked article",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(expectedExec).
					WithArgs(test.TestArticle.ID, test.TestArticle.Title, test.TestArticle.Content, test.TestArticle.ArticleLink,
						test.TestArticle.UserID, test.TestArticle.CreatedAt, test.TestArticle.UpdatedAt).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			err: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockBehaviour(mock)
			err := r.Save(context.Background(), *test.TestArticle)
			assert.Equal(t, c.err, err)
		})
	}
}

func TestList(t *testing.T) {
	gormDB, db, mock := test.NewMockDB(t)
	defer db.Close()

	r := NewRepository(gormDB)
	expectedCountQuery := "^SELECT count(.*) FROM \"articles\""
	expectedSelectQuery := "^SELECT id, title, article_link, created_at, updated_at FROM \"articles\" *"

	cases := []struct {
		name          string
		userID        string
		page          *pagination.Pagination
		mockBehaviour func(mock sqlmock.Sqlmock, userID string, page *pagination.Pagination)
		articles      []*article.Article
		total         int64
		err           error
	}{
		{
			name:   "should return first page of the bookmarked articles",
			userID: test.TestUser.ID,
			page: &pagination.Pagination{
				Limit:  2,
				Offset: 0,
			},
			mockBehaviour: func(mock sqlmock.Sqlmock, userID string, page *pagination.Pagination) {
				mock.ExpectQuery(expectedCountQuery).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
				mock.ExpectQuery(expectedSelectQuery).
					WithArgs(userID, page.Limit).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "article_link", "created_at", "updated_at"}).
						AddRow(test.TestArticle.ID, test.TestArticle.Title, test.TestArticle.ArticleLink, test.TestArticle.CreatedAt, test.TestArticle.UpdatedAt).
						AddRow(test.TestArticle2.ID, test.TestArticle2.Title, test.TestArticle2.ArticleLink, test.TestArticle2.CreatedAt, test.TestArticle2.UpdatedAt))
			},
			articles: []*article.Article{
				test.TestArticle,
				test.TestArticle2,
			},
			total: 2,
			err:   nil,
		},
		{
			name:   "should return second page of the bookmarked articles",
			userID: test.TestUser.ID,
			page: &pagination.Pagination{
				Limit:  2,
				Offset: 2,
			},
			mockBehaviour: func(mock sqlmock.Sqlmock, userID string, page *pagination.Pagination) {
				mock.ExpectQuery(expectedCountQuery).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery(expectedSelectQuery).
					WithArgs(userID, page.Limit, page.Offset).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "article_link", "created_at", "updated_at"}).
						AddRow(test.TestArticle3.ID, test.TestArticle3.Title, test.TestArticle3.ArticleLink, test.TestArticle3.CreatedAt, test.TestArticle3.UpdatedAt))
			},
			articles: []*article.Article{
				test.TestArticle3,
			},
			total: 1,
			err:   nil,
		},
		{
			name:   "should return empty slice when no user's article found",
			userID: test.TestUser.ID,
			page: &pagination.Pagination{
				Limit:  2,
				Offset: 0,
			},
			mockBehaviour: func(mock sqlmock.Sqlmock, userID string, page *pagination.Pagination) {
				mock.ExpectQuery(expectedCountQuery).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectQuery(expectedSelectQuery).
					WithArgs(userID, page.Limit).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			articles: []*article.Article{},
			total:    0,
			err:      nil,
		},
		{
			name:   "should return error when fail to fetch articles",
			userID: test.TestUser.ID,
			page: &pagination.Pagination{
				Limit:  2,
				Offset: 0,
			},
			mockBehaviour: func(mock sqlmock.Sqlmock, userID string, page *pagination.Pagination) {
				mock.ExpectQuery(expectedCountQuery).WillReturnError(gorm.ErrInvalidDB)
			},
			articles: []*article.Article{},
			total:    0,
			err:      gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockBehaviour(mock, c.userID, c.page)

			articles, total, err := r.List(context.Background(), c.userID, *c.page)
			assert.Equal(t, c.err, err)
			if err != nil {
				assert.Zero(t, total)
				assert.Empty(t, articles)
			}

			assert.Equal(t, c.total, total)
			for i, v := range articles {
				assert.Equal(t, c.articles[i].ID, v.ID)
				assert.Equal(t, c.articles[i].Title, v.Title)
				assert.Empty(t, v.Content)
				assert.Equal(t, c.articles[i].ArticleLink, v.ArticleLink)
				assert.Empty(t, v.UserID)
				assert.Equal(t, c.articles[i].CreatedAt, v.CreatedAt)
				assert.Equal(t, c.articles[i].UpdatedAt, v.UpdatedAt)
			}
		})
	}
}
