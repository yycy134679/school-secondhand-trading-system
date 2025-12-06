package model

import (
	"time"
)

// Tag 商品标签模型
type Tag struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"type:varchar(50);not null;uniqueIndex"`
	CategoryID int64     `json:"categoryId" gorm:"column:category_id;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
