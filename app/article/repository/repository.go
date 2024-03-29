package repository

import (
	"context"
	"errors"

	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/pagination"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *repository) FindByID(ctx context.Context, articleID string) (article *article.Article, err error) {
	err = r.db.First(&article, "id = ?", articleID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = validation.NewError(validation.NotFound, "no article found with given id")
	}
	return
}

func (r *repository) Update(ctx context.Context, arg article.Article) (updated *article.Article, err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&updated, "id = ?", arg.ID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = validation.NewError(validation.NotFound, "no article found with given id")
			}
			return err
		}

		if updated.UserID != arg.UserID {
			return validation.NewError(validation.Forbidden, "forbidden access")
		}

		updated.Title = arg.Title
		updated.Content = arg.Content
		updated.ArticleLink = arg.ArticleLink
		updated.UpdatedAt = arg.UpdatedAt

		return tx.Model(&updated).Updates(article.Article{
			Title:       arg.Title,
			Content:     arg.Content,
			ArticleLink: arg.ArticleLink,
			UpdatedAt:   arg.UpdatedAt,
		}).Error
	})

	return
}

func (r *repository) Delete(ctx context.Context, userID, articleID string) error {
	res := r.db.Where("id = ? AND user_id = ?", articleID, userID).Delete(&article.Article{})
	if res.RowsAffected == 0 && res.Error == nil {
		return validation.NewError(validation.BadRequest, "fail to delete article")
	}
	return res.Error
}
