package model

import "time"

// Product 商品模型
type Product struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	CategoryID   int64     `json:"categoryId"`
	ConditionID  int64     `json:"conditionId"`
	SellerID     int64     `json:"sellerId"`
	Status       string    `json:"status"`
	MainImageURL string    `json:"mainImageUrl" gorm:"column:main_image_url"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ProductImage 商品图片模型
type ProductImage struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	ProductID int64  `json:"productId"`
	URL       string `json:"url"`
	IsPrimary bool   `json:"isPrimary"`
	SortOrder int    `json:"sortOrder"`
}

// ProductDetailDTO 商品详情DTO
type ProductDetailDTO struct {
	ID             int64          `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Price          float64        `json:"price"`
	CategoryID     int64          `json:"categoryId"`
	ConditionID    int64          `json:"conditionId"`
	ConditionName  string         `json:"conditionName"`
	MainImageURL   string         `json:"mainImageUrl"`
	Images         []ProductImage `json:"images"`
	TagIDs         []int64        `json:"tagIds"`
	Seller         SellerInfo     `json:"seller"`
	ViewerIsSeller bool           `json:"viewerIsSeller"`
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	SellerWechat   *string        `json:"sellerWechat,omitempty"`
}

// ProductCardDTO 商品卡片DTO
type ProductCardDTO struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Price       float64   `json:"price"`
	MainImage   string    `json:"mainImageUrl"`
	Status      string    `json:"status"`
	SellerID    int64     `json:"sellerId"`
	CategoryID  int64     `json:"categoryId"`
	ConditionID int64     `json:"conditionId"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// SellerInfo 卖家简要信息
type SellerInfo struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarUrl string `json:"avatarUrl"`
}
