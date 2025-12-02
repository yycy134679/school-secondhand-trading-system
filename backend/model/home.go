// Package model 定义首页数据相关的数据结构
package model

// HomeData 首页数据响应结构
type HomeData struct {
	Recommendations []ProductCardDTO `json:"recommendations"` // 推荐商品列表
	Latest          []ProductCardDTO `json:"latest"`          // 最新商品列表
}
