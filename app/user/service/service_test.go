package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ryanadiputraa/unclatter/app/mocks"
	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	dummyEmail = "testuser@mail.com"
)

func TestCreateUser(t *testing.T) {
	cases := []struct {
		name              string
		arg               user.CreateUserArg
		expected          *user.User
		err               error
		mockRepoBehaviour func(mockRepo *mocks.UserRepository)
	}{
		{
			name: "should return created user",
			arg: user.CreateUserArg{
				Email:     dummyEmail,
				FirstName: "test",
				LastName:  "lastname",
			},
			expected: &user.User{
				ID:        "randomid",
				Email:     dummyEmail,
				FirstName: "test",
				LastName:  "lastname",
				CreatedAt: time.Now().UTC(),
			},
			err: nil,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("Save", mock.Anything).Return(nil)
			},
		},
		{
			name: "should fail to create user when given invalid email address",
			arg: user.CreateUserArg{
				Email:     "testusermailcom",
				FirstName: "test",
				LastName:  "lastname",
			},
			expected: nil,
			err:      errors.New("invalid email address"),
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				// no behaviour needed
			},
		},
		{
			name: "should fail to create user when error duplicate email from repository",
			arg: user.CreateUserArg{
				Email:     dummyEmail,
				FirstName: "test",
				LastName:  "lastname",
			},
			expected: nil,
			err:      validation.NewError(validation.BadRequest, "email already registered"),
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("Save", mock.Anything).Return(validation.NewError(validation.BadRequest, "email already registered"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := mocks.NewUserRepository(t)
			c.mockRepoBehaviour(r)

			s := NewService(logger.NewLogger(), r)
			user, err := s.CreateUser(context.Background(), c.arg)

			assert.Equal(t, c.err, err)
			if err != nil {
				return
			}

			assert.NotNil(t, user.ID)
			assert.Equal(t, c.expected.Email, user.Email)
			assert.Equal(t, c.expected.FirstName, user.FirstName)
			assert.Equal(t, c.expected.LastName, user.LastName)
			assert.NotNil(t, user.CreatedAt)
		})
	}
}
