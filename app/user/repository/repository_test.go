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

func TestSaveOrUpdate(t *testing.T) {
	gormDB, db, mock := test.NewMockDB(t)
	defer db.Close()

	r := NewRepository(gormDB)
	expectedInsertQuery := "INSERT INTO \"users\""

	user, err := user.NewUser(user.NewUserArg{
		Email:     "testuser@mail.com",
		FirstName: "test",
		LastName:  "user",
	})
	assert.NoError(t, err)

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
					WithArgs(user.ID, user.Email, user.FirstName, user.LastName, user.CreatedAt).
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
					WithArgs(user.ID, user.Email, user.FirstName, user.LastName, user.CreatedAt).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			err: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockBehaviour(mock)
			err := r.SaveOrUpdate(context.Background(), *user)
			assert.Equal(t, c.err, err)
		})
	}
}

func TestFindByID(t *testing.T) {
	gormDB, db, mock := test.NewMockDB(t)
	defer db.Close()

	r := NewRepository(gormDB)
	expectedQuery := "^SELECT (.+) FROM \"users\" *"
	testUser, err := user.NewUser(user.NewUserArg{
		Email:     "testuser@mail.com",
		FirstName: "test",
		LastName:  "user",
	})
	assert.NoError(t, err)

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
						testUser.ID,
						testUser.Email,
						testUser.FirstName,
						testUser.LastName,
						testUser.CreatedAt,
					),
				)
			},
			arg:  testUser.ID,
			user: testUser,
			err:  nil,
		},
		{
			name: "should return validation error missing user data when no user found with given id",
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WillReturnError(gorm.ErrRecordNotFound)
			},
			arg:  testUser.ID,
			user: nil,
			err:  validation.NewError(validation.BadRequest, "missing user data"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockBehaviour(mock)
			user, err := r.FindByID(context.Background(), testUser.ID)
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
