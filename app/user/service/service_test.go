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
	"gorm.io/gorm"
)

const (
	dummyEmail         = "testuser@mail.com"
	missingUserDataStr = "missing user data"
)

var testUser, _ = user.NewUser(user.NewUserArg{
	Email:     "testuser@mail.com",
	FirstName: "test",
	LastName:  "user",
})

func TestCreateUser(t *testing.T) {
	cases := []struct {
		name              string
		arg               user.NewUserArg
		expected          *user.User
		err               error
		mockRepoBehaviour func(mockRepo *mocks.UserRepository)
	}{
		{
			name: "should return created user",
			arg: user.NewUserArg{
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
				mockRepo.On("FindByEmail", context.Background(), mock.Anything).Return(nil, validation.NewError(validation.BadRequest, missingUserDataStr))
				mockRepo.On("Save", context.Background(), mock.Anything).Return(nil)
			},
		},
		{
			name: "should fail to create user when given invalid email address",
			arg: user.NewUserArg{
				Email:     "testusermailcom",
				FirstName: "test",
				LastName:  "lastname",
			},
			expected: nil,
			err:      errors.New("invalid email address"),
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", context.Background(), mock.Anything).Return(nil, validation.NewError(validation.BadRequest, missingUserDataStr))
			},
		},
		{
			name: "should fail to create user and return error from repository",
			arg: user.NewUserArg{
				Email:     dummyEmail,
				FirstName: "test",
				LastName:  "lastname",
			},
			expected: nil,
			err:      gorm.ErrInvalidDB,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", context.Background(), mock.Anything).Return(nil, validation.NewError(validation.BadRequest, missingUserDataStr))
				mockRepo.On("Save", context.Background(), mock.Anything).
					Return(gorm.ErrInvalidDB)
			},
		},
		{
			name: "should fail to create user and return error from repository",
			arg: user.NewUserArg{
				Email:     dummyEmail,
				FirstName: "test",
				LastName:  "lastname",
			},
			expected: nil,
			err:      gorm.ErrInvalidDB,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", context.Background(), mock.Anything).Return(nil, validation.NewError(validation.BadRequest, missingUserDataStr))
				mockRepo.On("Save", context.Background(), mock.Anything).
					Return(gorm.ErrInvalidDB)
			},
		},
		{
			name: "should fail to create user and return error from find user by email repository",
			arg: user.NewUserArg{
				Email:     dummyEmail,
				FirstName: "test",
				LastName:  "lastname",
			},
			expected: nil,
			err:      gorm.ErrInvalidDB,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", context.Background(), mock.Anything).Return(nil, gorm.ErrInvalidDB)
			},
		},
		{
			name: "should return already saved user",
			arg: user.NewUserArg{
				Email:     dummyEmail,
				FirstName: "test",
				LastName:  "lastname",
			},
			expected: testUser,
			err:      nil,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", context.Background(), mock.Anything).Return(testUser, nil)
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

func TestGetUserInfo(t *testing.T) {
	cases := []struct {
		name              string
		arg               string
		expected          *user.User
		err               error
		mockRepoBehaviour func(mockRepo *mocks.UserRepository)
	}{
		{
			name:     "should return user info with given id",
			arg:      testUser.ID,
			expected: testUser,
			err:      nil,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByID", context.Background(), mock.Anything).Return(testUser, nil)
			},
		},
		{
			name:     "should return error when user not found",
			arg:      testUser.ID,
			expected: nil,
			err:      validation.NewError(validation.BadRequest, missingUserDataStr),
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByID", context.Background(), mock.Anything).
					Return(nil, validation.NewError(validation.BadRequest, missingUserDataStr))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := mocks.NewUserRepository(t)
			c.mockRepoBehaviour(r)

			s := NewService(logger.NewLogger(), r)
			user, err := s.GetUserInfo(context.Background(), c.arg)

			assert.Equal(t, c.err, err)
			if err != nil {
				assert.Empty(t, user)
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
