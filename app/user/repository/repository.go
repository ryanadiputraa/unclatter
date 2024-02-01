package repository

import (
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

func (r *repository) Save(arg user.User) error {
	err := r.db.Create(arg).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			serviceErr := validation.NewError(validation.BadRequest, "email already registered")
			return serviceErr
		}
	}
	return nil
}
