package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthProvider struct {
	ID             string    `json:"id" gorm:"type:varchar"`
	Provider       string    `json:"provider" gorm:"type:varchar;not null"`
	ProviderUserID string    `json:"provider_user_id" gorm:"type:varchar;not null"`
	UserID         string    `json:"user_id" gorm:"type:varchar;not null"`
	CreatedAt      time.Time `json:"-" gorm:"not null"`
}

type NewAuthProviderArg struct {
	Provider       string
	ProviderUserID string
	UserID         string
}

func NewAuthProvider(arg NewAuthProviderArg) *AuthProvider {
	return &AuthProvider{
		ID:             uuid.NewString(),
		Provider:       arg.Provider,
		ProviderUserID: arg.ProviderUserID,
		UserID:         arg.UserID,
		CreatedAt:      time.Now().UTC(),
	}
}

type AuthService interface {
	AddUserAuthProvider(ctx context.Context, arg NewAuthProviderArg) (*AuthProvider, error)
}

type AuthProviderRepository interface {
	Save(ctx context.Context, authProvider AuthProvider) error
}
