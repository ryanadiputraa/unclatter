package service

import (
	"context"
	"testing"

	"github.com/ryanadiputraa/unclatter/app/auth"
	"github.com/ryanadiputraa/unclatter/app/mocks"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const (
	missingUserDataStr = "missing user data"
)

func TestAddUserAuthProvider(t *testing.T) {
	cases := []struct {
		name              string
		arg               auth.NewAuthProviderArg
		expected          *auth.AuthProvider
		err               error
		mockRepoBehaviour func(mockRepo *mocks.AuthProviderRepository)
	}{
		{
			name: "should return added user auth provider",
			arg: auth.NewAuthProviderArg{
				Provider:       test.TestAuthProvider.Provider,
				ProviderUserID: test.TestAuthProvider.ProviderUserID,
				UserID:         test.TestAuthProvider.UserID,
			},
			expected: test.TestAuthProvider,
			err:      nil,
			mockRepoBehaviour: func(mockRepo *mocks.AuthProviderRepository) {
				mockRepo.On("Save", context.Background(), mock.Anything).Return(nil)
			},
		},
		{
			name: "should return error when fail to save auth provider",
			arg: auth.NewAuthProviderArg{
				Provider:       test.TestAuthProvider.Provider,
				ProviderUserID: test.TestAuthProvider.ProviderUserID,
				UserID:         test.TestAuthProvider.UserID,
			},
			expected: nil,
			err:      gorm.ErrInvalidDB,
			mockRepoBehaviour: func(mockRepo *mocks.AuthProviderRepository) {
				mockRepo.On("Save", context.Background(), mock.Anything).Return(gorm.ErrInvalidDB)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := mocks.NewAuthProviderRepository(t)
			c.mockRepoBehaviour(r)

			s := NewService(logger.NewLogger(), r)
			authProvider, err := s.AddUserAuthProvider(context.Background(), c.arg)

			assert.Equal(t, c.err, err)
			if err != nil {
				return
			}
			assert.NotNil(t, authProvider.ID)
			assert.Equal(t, c.expected.Provider, authProvider.Provider)
			assert.Equal(t, c.expected.ProviderUserID, authProvider.ProviderUserID)
			assert.Equal(t, c.expected.UserID, authProvider.UserID)
			assert.NotNil(t, authProvider.CreatedAt)
		})
	}
}
