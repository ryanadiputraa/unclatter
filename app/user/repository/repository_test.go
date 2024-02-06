package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSave(t *testing.T) {
	gormDB, db, mock := test.NewMockDB(t)
	defer db.Close()

	r := NewRepository(gormDB)
	expectedInsertQuery := "INSERT INTO \"users\""

	cases := []struct {
		name          string
		mockBehaviour func(mock sqlmock.Sqlmock)
		err           error
	}{
		{
			name: "should insert new user",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(expectedInsertQuery).
					WithArgs(test.TestUser.ID, test.TestUser.Email, test.TestUser.FirstName, test.TestUser.LastName, test.TestUser.CreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			err: nil,
		},
		{
			name: "should fail to insert user and return error",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(expectedInsertQuery).
					WithArgs(test.TestUser.ID, test.TestUser.Email, test.TestUser.FirstName, test.TestUser.LastName, test.TestUser.CreatedAt).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			err: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockBehaviour(mock)
			err := r.Save(context.Background(), *test.TestUser)
			assert.Equal(t, c.err, err)
		})
	}
}

func TestFindByID(t *testing.T) {
	gormDB, db, mock := test.NewMockDB(t)
	defer db.Close()

	r := NewRepository(gormDB)
	expectedQuery := "^SELECT (.+) FROM \"users\" *"

	cases := []struct {
		name          string
		mockBehaviour func(mock sqlmock.Sqlmock)
		arg           string
		user          *user.User
		err           error
	}{
		{
			name: "should return user with given id",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"email",
						"first_name",
						"last_name",
						"created_at",
					}).AddRow(
						test.TestUser.ID,
						test.TestUser.Email,
						test.TestUser.FirstName,
						test.TestUser.LastName,
						test.TestUser.CreatedAt,
					),
				)
			},
			arg:  test.TestUser.ID,
			user: test.TestUser,
			err:  nil,
		},
		{
			name: "should return validation error missing user data when no user found with given id",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WillReturnError(gorm.ErrRecordNotFound)
			},
			arg:  test.TestUser.ID,
			user: nil,
			err:  validation.NewError(validation.BadRequest, "missing user data"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockBehaviour(mock)
			user, err := r.FindByID(context.Background(), test.TestUser.ID)
			assert.Equal(t, c.err, err)
			if err != nil {
				assert.Empty(t, user)
				return
			}

			assert.Equal(t, c.user.ID, user.ID)
			assert.Equal(t, c.user.Email, user.Email)
			assert.Equal(t, c.user.FirstName, user.FirstName)
			assert.Equal(t, c.user.LastName, user.LastName)
			assert.NotNil(t, user.CreatedAt)
		})
	}
}

func TestFindByEmail(t *testing.T) {
	gormDB, db, mock := test.NewMockDB(t)
	defer db.Close()

	r := NewRepository(gormDB)
	expectedQuery := "^SELECT (.+) FROM \"users\" *"

	cases := []struct {
		name          string
		mockBehaviour func(mock sqlmock.Sqlmock)
		arg           string
		user          *user.User
		err           error
	}{
		{
			name: "should return user with given id",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"email",
						"first_name",
						"last_name",
						"created_at",
					}).AddRow(
						test.TestUser.ID,
						test.TestUser.Email,
						test.TestUser.FirstName,
						test.TestUser.LastName,
						test.TestUser.CreatedAt,
					),
				)
			},
			arg:  test.TestUser.ID,
			user: test.TestUser,
			err:  nil,
		},
		{
			name: "should return validation error missing user data when no user found with given id",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WillReturnError(gorm.ErrRecordNotFound)
			},
			arg:  test.TestUser.ID,
			user: nil,
			err:  validation.NewError(validation.BadRequest, "missing user data"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockBehaviour(mock)
			user, err := r.FindByEmail(context.Background(), test.TestUser.Email)
			assert.Equal(t, c.err, err)
			if err != nil {
				assert.Empty(t, user)
				return
			}

			assert.Equal(t, c.user.ID, user.ID)
			assert.Equal(t, c.user.Email, user.Email)
			assert.Equal(t, c.user.FirstName, user.FirstName)
			assert.Equal(t, c.user.LastName, user.LastName)
			assert.NotNil(t, user.CreatedAt)
		})
	}
}
