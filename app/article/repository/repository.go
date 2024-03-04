package repository

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) article.ArticleRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, arg article.Article) error {
	err := r.db.Create(&arg).Error
	if err == gorm.ErrDuplicatedKey {
		err = validation.NewError(validation.BadRequest, "title is already in use")
	}
	return err
}
