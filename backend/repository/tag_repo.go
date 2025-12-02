package repository

import (
	"context"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"

	"gorm.io/gorm"
)

// TagRepository 标签仓库接口
type TagRepository interface {
	// ListAll 获取所有标签
	ListAll(ctx context.Context) ([]*model.Tag, error)
	// Create 创建标签
	Create(ctx context.Context, tag *model.Tag) error
	// Update 更新标签
	Update(ctx context.Context, tag *model.Tag) error
	// Delete 删除标签
	Delete(ctx context.Context, id int64) error
	// GetByID 根据ID获取标签
	GetByID(ctx context.Context, id int64) (*model.Tag, error)
	// CountProductsByTag 统计该标签下的商品数量
	CountProductsByTag(ctx context.Context, tagID int64) (int64, error)
	// GetByIDs 根据ID列表获取标签
	GetByIDs(ctx context.Context, ids []int64) ([]*model.Tag, error)
}

// tagRepository 标签仓库实现
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓库实例
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{
		db: db,
	}
}

// ListAll 获取所有标签
func (r *tagRepository) ListAll(ctx context.Context) ([]*model.Tag, error) {
	var tags []*model.Tag
	err := r.db.WithContext(ctx).Order("id ASC").Find(&tags).Error
	return tags, err
}

// Create 创建标签
func (r *tagRepository) Create(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// Update 更新标签
func (r *tagRepository) Update(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Model(tag).Update("name", tag.Name).Error
}

// Delete 删除标签
func (r *tagRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Tag{}, id).Error
}

// GetByID 根据ID获取标签
func (r *tagRepository) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// CountProductsByTag 统计该标签下的商品数量
func (r *tagRepository) CountProductsByTag(ctx context.Context, tagID int64) (int64, error) {
	var count int64
	// 统计product_tags表中引用该标签的记录数
	err := r.db.WithContext(ctx).Model(&struct{}{}).Table("product_tags").Where("tag_id = ?", tagID).Count(&count).Error
	return count, err
}

// GetByIDs 根据ID列表获取标签
func (r *tagRepository) GetByIDs(ctx context.Context, ids []int64) ([]*model.Tag, error) {
	var tags []*model.Tag
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}
