package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
)

// ViewRecordRepository 浏览记录仓库接口
type ViewRecordRepository interface {
	// AddView 添加浏览记录
	AddView(ctx context.Context, userID, productID int64) error
	// ListRecentViews 获取用户最近浏览记录
	ListRecentViews(ctx context.Context, userID int64, limit int) ([]model.UserRecentView, error)
}

// viewRecordRepository 浏览记录仓库实现
type viewRecordRepository struct {
	db *gorm.DB
}

// NewViewRecordRepository 创建浏览记录仓库实例
func NewViewRecordRepository(db *gorm.DB) ViewRecordRepository {
	return &viewRecordRepository{
		db: db,
	}
}

// AddView 添加浏览记录
// 由于 DB 已有触发器 user_recent_views_prune_after_insert 自动裁剪至 20 条,这里不需手动删除
func (r *viewRecordRepository) AddView(ctx context.Context, userID, productID int64) error {
	record := &model.UserRecentView{
		UserID:    userID,
		ProductID: productID,
		ViewedAt:  time.Now(),
	}

	// 插入记录,触发器会自动裁剪旧记录
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}

	return nil
}

// ListRecentViews 按时间倒序返回用户最近浏览记录
// limit: 返回的最大记录数
func (r *viewRecordRepository) ListRecentViews(ctx context.Context, userID int64, limit int) ([]model.UserRecentView, error) {
	var views []model.UserRecentView

	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("viewed_at DESC, id DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&views).Error; err != nil {
		return nil, err
	}

	return views, nil
}
