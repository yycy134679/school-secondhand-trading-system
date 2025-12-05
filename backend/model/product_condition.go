package model

import "time"

// ProductCondition 商品新旧程度模型
type ProductCondition struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"type:varchar(32);not null;uniqueIndex"`
	Name      string    `json:"name" gorm:"type:varchar(50);not null"`
	SortOrder int32     `json:"sortOrder" gorm:"type:smallint;not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (ProductCondition) TableName() string {
	return "product_conditions"
}
