// Package repository 定义了数据访问层接口和实现
package repository

import (
	"context"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"gorm.io/gorm"
)

// ViewRecordRepository 浏览记录仓库接口
// 负责用户浏览记录的数据库操作，用于推荐系统和最近浏览功能
type ViewRecordRepository interface {
	// AddView 添加一条浏览记录
	// 参数：
	//   - ctx: 上下文
	//   - userID: 用户ID
	//   - productID: 商品ID
	// 返回：
	//   - error: 错误信息
	// 说明：
	//   - 插入一条新的浏览记录，数据库会自动设置 viewed_at 为当前时间
	//   - 数据库触发器会自动维护每个用户最多20条记录的限制
	//   - 如果用户重复浏览同一商品，会产生多条记录（用于计算浏览频次）
	AddView(ctx context.Context, userID int64, productID int64) error

	// ListRecentViews 获取用户最近浏览记录
	// 参数：
	//   - ctx: 上下文
	//   - userID: 用户ID
	//   - limit: 限制返回条数（建议不超过20）
	// 返回：
	//   - []*model.UserRecentView: 浏览记录列表，按时间倒序
	//   - error: 错误信息
	// 说明：
	//   - 返回指定用户的最近浏览记录，按 viewed_at 时间倒序排列
	//   - 主要用于推荐系统分析用户偏好和生成"最近浏览"列表
	ListRecentViews(ctx context.Context, userID int64, limit int) ([]*model.UserRecentView, error)
}

// viewRecordRepositoryImpl 浏览记录仓库实现
type viewRecordRepositoryImpl struct {
	db *gorm.DB
}

// NewViewRecordRepository 创建浏览记录仓库实例
func NewViewRecordRepository(db *gorm.DB) ViewRecordRepository {
	return &viewRecordRepositoryImpl{
		db: db,
	}
}

// AddView 实现添加浏览记录
func (r *viewRecordRepositoryImpl) AddView(ctx context.Context, userID int64, productID int64) error {
	// 创建浏览记录实例
	view := &model.UserRecentView{
		UserID:    userID,
		ProductID: productID,
		// ViewedAt 会由数据库自动设置为 NOW()
	}

	// 插入数据库
	// 注意：数据库有触发器 user_recent_views_prune_after_insert
	// 会在插入后自动删除超过20条的老记录
	if err := r.db.WithContext(ctx).Create(view).Error; err != nil {
		return err
	}

	return nil
}

// ListRecentViews 实现获取最近浏览记录
func (r *viewRecordRepositoryImpl) ListRecentViews(ctx context.Context, userID int64, limit int) ([]*model.UserRecentView, error) {
	var views []*model.UserRecentView

	// 查询指定用户的浏览记录，按时间倒序，限制条数
	// 使用复合索引 idx_views_user_time (user_id, viewed_at DESC, id DESC)
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("viewed_at DESC, id DESC").
		Limit(limit).
		Find(&views).Error

	if err != nil {
		return nil, err
	}

	return views, nil
}