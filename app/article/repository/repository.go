package repository

import (
	"context"
	"errors"

	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/pagination"
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
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = validation.NewError(validation.BadRequest, "title is already in use")
	}
	return err
}

func (r *repository) List(ctx context.Context, userID string, page pagination.Pagination) (articles []*article.Article, total int64, err error) {
	err = r.db.Model(&article.Article{}).Count(&total).Error
	if err != nil {
		return
	}

	err = r.db.
		Select("id, title, article_link, created_at, updated_at").
		Where("user_id = ?", userID).
		Order("updated_at DESC, created_at DESC").
		Limit(page.Limit).Offset(page.Offset).
		Find(&articles).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		articles = []*article.Article{}
		err = nil
		return
	}

	return
}
