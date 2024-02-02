package repository

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/auth"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) auth.AuthProviderRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, authProvider auth.AuthProvider) error {
	var provider auth.AuthProvider
	err := r.db.Where("provider = ? AND provider_user_id = ?", authProvider.Provider, authProvider.ProviderUserID).
		First(&provider).Error

	if err == gorm.ErrRecordNotFound {
		if err = r.db.Create(&authProvider).Error; err != nil {
			return err
		}
		return nil
	}

	return err
}
