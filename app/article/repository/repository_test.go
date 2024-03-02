package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
