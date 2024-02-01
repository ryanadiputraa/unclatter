package repository

import (
	"github.com/ryanadiputraa/unclatter/app/user"
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
	return r.db.Create(arg).Error
}
