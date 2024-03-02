package repository

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/article"
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

func (r *repository) Save(ctx context.Context, arg article.Article) (err error) {
	return r.db.Create(&arg).Error
}
