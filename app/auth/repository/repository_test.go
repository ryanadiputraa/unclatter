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
	expectedQuery := "^SELECT (.+) FROM \"auth_providers\" *"
	expectedExec := "INSERT INTO \"auth_providers\""

	cases := []struct {
		name          string
		mockBehaviour func(mock sqlmock.Sqlmock)
		err           error
	}{
		{
			name: "should insert new user auth provider",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectExec(expectedExec).
					WithArgs(test.TestAuthProvider.ID, test.TestAuthProvider.Provider, test.TestAuthProvider.ProviderUserID, test.TestAuthProvider.UserID, test.TestAuthProvider.CreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			err: nil,
		},
		{
			name: "should return error when fail to insert new user auth provider",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectExec(expectedExec).
					WithArgs(test.TestAuthProvider.ID, test.TestAuthProvider.Provider, test.TestAuthProvider.ProviderUserID, test.TestAuthProvider.UserID, test.TestAuthProvider.CreatedAt).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			err: gorm.ErrInvalidDB,
		},
		{
			name: "should do nothing when user auth provider already recorded",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id",
					"provider",
					"provider_user_id",
					"user_id",
					"created_at",
				}).AddRow(
					test.TestAuthProvider.ID,
					test.TestAuthProvider.Provider,
					test.TestAuthProvider.ProviderUserID,
					test.TestAuthProvider.UserID,
					test.TestAuthProvider.CreatedAt)

				mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
			},
			err: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockBehaviour(mock)
			err := r.Save(context.Background(), *test.TestAuthProvider)
			assert.Equal(t, c.err, err)
		})
	}
}
