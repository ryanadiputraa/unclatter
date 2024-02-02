package repository

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) user.UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) SaveOrUpdate(ctx context.Context, arg user.User) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"first_name", "last_name",
		}),
	}).Create(&arg).Error
}
