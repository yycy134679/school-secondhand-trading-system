package repository

import (
	"context"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"

	"gorm.io/gorm"
)

// ProductConditionRepository 新旧程度仓库接口
type ProductConditionRepository interface {
	ListAll(ctx context.Context) ([]model.ProductCondition, error)
}

// productConditionRepo 仓库实现
type productConditionRepo struct {
	db *gorm.DB
}

// NewProductConditionRepository 创建仓库实例
func NewProductConditionRepository(db *gorm.DB) ProductConditionRepository {
	return &productConditionRepo{db: db}
}

// ListAll 获取所有新旧程度，按排序字段升序
func (r *productConditionRepo) ListAll(ctx context.Context) ([]model.ProductCondition, error) {
	var conditions []model.ProductCondition
	err := r.db.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&conditions).Error
	return conditions, err
}
