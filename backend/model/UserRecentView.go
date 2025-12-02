// Package model 定义了系统中的数据模型（实体类）
package model

import "time"

// UserRecentView 用户最近浏览记录模型，对应数据库中的 user_recent_views 表
//
// 字段说明：
//   - ID: 记录的唯一标识（主键，自增）
//   - UserID: 用户ID（外键，关联users表）
//   - ProductID: 商品ID（外键，关联products表）
//   - ViewedAt: 浏览时间（自动设置为当前时间）
//
// 数据库表结构对应：
//   CREATE TABLE user_recent_views (
//     id         BIGSERIAL PRIMARY KEY,
//     user_id    BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
//     product_id BIGINT       NOT NULL REFERENCES products(id) ON DELETE CASCADE,
//     viewed_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
//   );
//
// 业务规则：
//   - 每个用户最多保留最近20条浏览记录，由数据库触发器自动裁剪
//   - 用于生成"猜你喜欢"推荐列表
//   - 当用户或商品被删除时，相关浏览记录会被级联删除
type UserRecentView struct {
	ID        int64     `json:"id" gorm:"primaryKey"`                      // 主键ID
	UserID    int64     `json:"user_id" gorm:"not null;index"`            // 用户ID（外键）
	ProductID int64     `json:"product_id" gorm:"not null;index"`         // 商品ID（外键）
	ViewedAt  time.Time `json:"viewed_at" gorm:"not null;default:now()"` // 浏览时间
}

// TableName 指定UserRecentView模型对应的数据库表名
func (UserRecentView) TableName() string {
	return "user_recent_views"
}
