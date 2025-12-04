package model

import "time"

// UserRecentView 用户最近浏览记录模型
type UserRecentView struct {
	ID        int64     `json:"id" gorm:"primaryKey;column:id"`
	UserID    int64     `json:"userId" gorm:"column:user_id;not null;index:idx_views_user_time"`
	ProductID int64     `json:"productId" gorm:"column:product_id;not null"`
	ViewedAt  time.Time `json:"viewedAt" gorm:"column:viewed_at;not null;default:NOW();index:idx_views_user_time"`
}

// TableName 指定表名
func (UserRecentView) TableName() string {
	return "user_recent_views"
}
