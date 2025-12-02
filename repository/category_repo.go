package repository

import (
	"context"
	"school-secondhand-trading-system/model"

	"gorm.io/gorm"
)

// CategoryRepository 分类仓库接口
type CategoryRepository interface {
	ListAll(ctx context.Context) ([]model.Category, error)
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int64) error
	CountProductsByCategory(ctx context.Context, id int64) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Category, error)
}

// categoryRepo 分类仓库实现
type categoryRepo struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库实例
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

// ListAll 获取所有分类
func (r *categoryRepo) ListAll(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	return categories, err
}

// Create 创建分类
func (r *categoryRepo) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// Update 更新分类
func (r *categoryRepo) Update(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete 删除分类
func (r *categoryRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

// CountProductsByCategory 统计分类下的商品数量
func (r *categoryRepo) CountProductsByCategory(ctx context.Context, id int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Product{}).Where("category_id = ?", id).Count(&count).Error
	return count, err
}

// GetByID 根据ID获取分类
func (r *categoryRepo) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}
