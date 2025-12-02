package model

import (
	"time"
)

// Category 分类模型，映射categories表
type Category struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// CategoryDTO 分类数据传输对象，用于API响应
type CategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ToDTO 将Category模型转换为CategoryDTO
func (c *Category) ToDTO() *CategoryDTO {
	return &CategoryDTO{
		ID:   c.ID,
		Name: c.Name,
	}
}
