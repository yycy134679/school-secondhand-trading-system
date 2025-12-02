package repository

import (
	"context"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"

	"gorm.io/gorm"
)

// CategoryRepository 分类仓库接口
type CategoryRepository interface {
	// ListAll 获取所有分类
	ListAll(ctx context.Context) ([]*model.Category, error)
	// Create 创建分类
	Create(ctx context.Context, category *model.Category) error
	// Update 更新分类
	Update(ctx context.Context, category *model.Category) error
	// Delete 删除分类
	Delete(ctx context.Context, id int64) error
	// GetByID 根据ID获取分类
	GetByID(ctx context.Context, id int64) (*model.Category, error)
	// CountProductsByCategory 统计该分类下的商品数量
	CountProductsByCategory(ctx context.Context, categoryID int64) (int64, error)
}

// categoryRepository 分类仓库实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库实例
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

// ListAll 获取所有分类
func (r *categoryRepository) ListAll(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.WithContext(ctx).Order("id ASC").Find(&categories).Error
	return categories, err
}

// Create 创建分类
func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// Update 更新分类
func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Model(category).Update("name", category.Name).Error
}

// Delete 删除分类
func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

// GetByID 根据ID获取分类
func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// CountProductsByCategory 统计该分类下的商品数量
func (r *categoryRepository) CountProductsByCategory(ctx context.Context, categoryID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}
