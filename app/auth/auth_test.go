package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUserAuthProvider(t *testing.T) {
	uuid1 := uuid.NewString()
	uuid2 := uuid.NewString()

	cases := []struct {
		name     string
		arg      NewAuthProviderArg
		expected *AuthProvider
	}{
		{
			name: "should return an auth provider instance",
			arg: NewAuthProviderArg{
				Provider:       "google",
				ProviderUserID: "273120381907",
				UserID:         uuid1,
			},
			expected: &AuthProvider{
				ID:             uuid2,
				Provider:       "google",
				ProviderUserID: "273120381907",
				UserID:         uuid1,
				CreatedAt:      time.Now().UTC(),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			user := NewAuthProvider(c.arg)

			assert.NotEmpty(t, user.ID)
			assert.Equal(t, c.expected.Provider, user.Provider)
			assert.Equal(t, c.expected.ProviderUserID, user.ProviderUserID)
			assert.Equal(t, c.expected.UserID, user.UserID)
			assert.NotEmpty(t, user.CreatedAt)
		})
	}
}
