package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ryanadiputraa/unclatter/app/mocks"
	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const (
	missingUserDataStr = "missing user data"
)

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
				Email:     test.TestUser.Email,
				FirstName: test.TestUser.FirstName,
				LastName:  test.TestUser.LastName,
			},
			expected: test.TestUser,
			err:      nil,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", context.Background(), mock.Anything).Return(nil, validation.NewError(validation.BadRequest, missingUserDataStr))
				mockRepo.On("Save", context.Background(), mock.Anything).Return(nil)
			},
		},
		{
			name: "should fail to create user when given invalid email address",
			arg: user.NewUserArg{
				Email:     "invalidmail.com",
				FirstName: test.TestUser.FirstName,
				LastName:  test.TestUser.LastName,
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
				Email:     test.TestUser.Email,
				FirstName: test.TestUser.FirstName,
				LastName:  test.TestUser.LastName,
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
				Email:     test.TestUser.Email,
				FirstName: test.TestUser.FirstName,
				LastName:  test.TestUser.LastName,
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
				Email:     test.TestUser.Email,
				FirstName: test.TestUser.FirstName,
				LastName:  test.TestUser.LastName,
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
				Email:     test.TestUser.Email,
				FirstName: test.TestUser.FirstName,
				LastName:  test.TestUser.LastName,
			},
			expected: test.TestUser,
			err:      nil,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", context.Background(), mock.Anything).Return(test.TestUser, nil)
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
			arg:      test.TestUser.ID,
			expected: test.TestUser,
			err:      nil,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByID", context.Background(), mock.Anything).Return(test.TestUser, nil)
			},
		},
		{
			name:     "should return error when user not found",
			arg:      test.TestUser.ID,
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
