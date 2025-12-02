package model

import (
	"time"
)

// Tag 标签模型，映射tags表
type Tag struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// TagDTO 标签数据传输对象，用于API响应
type TagDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ToDTO 将Tag模型转换为TagDTO
func (t *Tag) ToDTO() *TagDTO {
	return &TagDTO{
		ID:   t.ID,
		Name: t.Name,
	}
}
