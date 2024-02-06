package repository

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) user.UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, user user.User) error {
	return r.db.Create(&user).Error
}

func (r *repository) FindByID(ctx context.Context, userID string) (user *user.User, err error) {
	err = r.db.Where("id = ?", userID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		err = validation.NewError(validation.BadRequest, "missing user data")
	}
	return
}

func (r *repository) FindByEmail(ctx context.Context, email string) (user *user.User, err error) {
	err = r.db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		err = validation.NewError(validation.BadRequest, "missing user data")
	}
	return
}
