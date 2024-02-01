package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ryanadiputraa/unclatter/internal/domain/validation"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"-"`
}

func NewUser(email, firstName, lastName string) (*User, error) {
	if isValid := validation.IsValidEmail(email); !isValid {
		return nil, errors.New("invalid email address")
	}

	return &User{
		ID:        uuid.NewString(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now().UTC(),
	}, nil
}
