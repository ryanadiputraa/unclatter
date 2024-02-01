package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSave(t *testing.T) {
	gormDB, db, mock := test.NewMockDB(t)
	defer db.Close()

	r := NewRepository(gormDB)

	user, err := user.NewUser(user.CreateUserArg{
		Email:     "testuser@mail.com",
		FirstName: "test",
		LastName:  "user",
	})
	assert.Nil(t, err)

	cases := []struct {
		name              string
		mockRepoBehaviour func(mock sqlmock.Sqlmock)
		err               error
	}{
		{
			name: "should insert new user",
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").
					WithArgs(user.ID, user.Email, user.FirstName, user.LastName, user.CreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			err: nil,
		},
		{
			name: "should fail to insert user and return error",
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").
					WithArgs(user.ID, user.Email, user.FirstName, user.LastName, user.CreatedAt).
					WillReturnError(gorm.ErrDuplicatedKey)
				mock.ExpectRollback()
			},
			err: gorm.ErrDuplicatedKey,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockRepoBehaviour(mock)
			err := r.Save(*user)
			assert.Equal(t, c.err, err)
		})
	}
}
