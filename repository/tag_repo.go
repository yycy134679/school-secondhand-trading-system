package repository

import (
	"context"
	"school-secondhand-trading-system/model"

	"gorm.io/gorm"
)

// TagRepository 标签仓库接口
type TagRepository interface {
	ListAll(ctx context.Context) ([]model.Tag, error)
	Create(ctx context.Context, tag *model.Tag) error
	Update(ctx context.Context, tag *model.Tag) error
	Delete(ctx context.Context, id int64) error
	CountProductsByTag(ctx context.Context, id int64) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Tag, error)
}

// tagRepo 标签仓库实现
type tagRepo struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓库实例
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepo{db: db}
}

// ListAll 获取所有标签
func (r *tagRepo) ListAll(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error
	return tags, err
}

// Create 创建标签
func (r *tagRepo) Create(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// Update 更新标签
func (r *tagRepo) Update(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// Delete 删除标签
func (r *tagRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Tag{}, id).Error
}

// CountProductsByTag 统计标签关联的商品数量
func (r *tagRepo) CountProductsByTag(ctx context.Context, id int64) (int64, error) {
	var count int64
	// 通过product_tags关联表查询关联的商品数量
	err := r.db.WithContext(ctx).Table("product_tags").Where("tag_id = ?", id).Count(&count).Error
	return count, err
}

// GetByID 根据ID获取标签
func (r *tagRepo) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}
