package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ryanadiputraa/unclatter/app/validation"
)

type User struct {
	ID        string    `json:"id" gorm:"type:varchar"`
	Email     string    `json:"email" gorm:"type:varchar;unique;not null"`
	FirstName string    `json:"first_name" gorm:"type:varchar;not null"`
	LastName  string    `json:"last_name" gorm:"type:varchar;not null"`
	CreatedAt time.Time `json:"-" gorm:"type:varchar;not null"`
}

type CreateUserArg struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewUser(arg CreateUserArg) (*User, error) {
	if isValid := validation.IsValidEmail(arg.Email); !isValid {
		return nil, errors.New("invalid email address")
	}

	return &User{
		ID:        uuid.NewString(),
		Email:     arg.Email,
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		CreatedAt: time.Now().UTC(),
	}, nil
}

type UserService interface {
	CreateUser(ctx context.Context, arg CreateUserArg) (*User, error)
}

type UserRepository interface {
	SaveOrUpdate(c context.Context, arg User) error
}
