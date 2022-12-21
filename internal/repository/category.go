package repository

import (
	"context"

	"github.com/SemmiDev/kanbanapp/internal/entity"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	categories := []entity.Category{}
	err := r.db.WithContext(ctx).Where("user_id = ?", id).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	err = r.db.WithContext(ctx).Create(category).Error
	if err != nil {
		return 0, err
	}
	return category.ID, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	total := len(categories)
	err := r.db.WithContext(ctx).CreateInBatches(categories, total).Error
	return err
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	user := entity.Category{}
	err := r.db.WithContext(ctx).First(&user, id).Error
	return user, err
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	err := r.db.WithContext(ctx).Save(category).Error
	return err
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
	return err
}
