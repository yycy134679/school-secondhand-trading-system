package model

import "time"

// 商品状态常量
const (
	ProductStatusForSale  = "ForSale"
	ProductStatusSold     = "Sold"
	ProductStatusDelisted = "Delisted"
)

// Product represents a product listing in the database
// 映射products表的字段：ID/SellerID/Title/Description/Price/ConditionID/CategoryID/Status/MainImageURL/CreatedAt/UpdatedAt
type Product struct {
	ID           int64     `json:"id" db:"id"`
	SellerID     int64     `json:"seller_id" db:"seller_id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Price        int64     `json:"price" db:"price"` // 数据库中是NUMERIC(10,2)，Go中使用int64存储分
	ConditionID  int64     `json:"condition_id" db:"condition_id"`
	CategoryID   int64     `json:"category_id" db:"category_id"`
	Status       string    `json:"status" db:"status"` // ForSale, Sold, Delisted
	MainImageURL string    `json:"main_image_url" db:"main_image_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ProductImage represents a product image in the database
// 映射product_images表：ID/ProductID/URL/SortOrder/IsPrimary
// 注意：唯一索引约束每商品最多一张主图
type ProductImage struct {
	ID        int64  `json:"id" db:"id"`
	ProductID int64  `json:"product_id" db:"product_id"`
	URL       string `json:"url" db:"url"`
	SortOrder int    `json:"sort_order" db:"sort_order"`
	IsPrimary bool   `json:"is_primary" db:"is_primary"`
}

// ProductTag represents the many-to-many relationship between products and tags
type ProductTag struct {
	ProductID int64 `json:"product_id" db:"product_id"`
	TagID     int64 `json:"tag_id" db:"tag_id"`
}

// ProductCardDTO represents a product card for listing views
type ProductCardDTO struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Price        int64     `json:"price"`
	MainImageUrl string    `json:"mainImageUrl"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ProductDetailDTO represents detailed product information
type ProductDetailDTO struct {
	ID             int64             `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Price          int64             `json:"price"`
	ConditionID    int64             `json:"conditionId"`
	ConditionName  string            `json:"conditionName,omitempty"`
	CategoryID     int64             `json:"categoryId"`
	Status         string            `json:"status"`
	MainImageUrl   string            `json:"mainImageUrl"`
	Images         []ProductImageDTO `json:"images"`
	TagIds         []int64           `json:"tagIds"`
	Seller         SellerInfoDTO     `json:"seller"`
	ViewerIsSeller bool              `json:"viewerIsSeller"`
	SellerWechat   string            `json:"sellerWechat,omitempty"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

// ProductImageDTO represents product image info in DTO
type ProductImageDTO struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sortOrder"`
	IsPrimary bool   `json:"isPrimary"`
}

// SellerInfoDTO represents seller information in product detail
type SellerInfoDTO struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

// CreateProductRequest represents the request for creating a product
type CreateProductRequest struct {
	Title       string  `json:"title" binding:"required,max=100"`
	Description string  `json:"description" binding:"required"`
	Price       int64   `json:"price" binding:"required,gt=0"`
	ConditionID int64   `json:"conditionId" binding:"required"`
	CategoryID  int64   `json:"categoryId" binding:"required"`
	TagIds      []int64 `json:"tagIds" binding:"required"`
}

// UpdateProductRequest represents the request for updating a product
type UpdateProductRequest struct {
	Title       string  `json:"title" binding:"omitempty,max=100"`
	Description string  `json:"description"`
	Price       int64   `json:"price" binding:"omitempty,gt=0"`
	ConditionID int64   `json:"conditionId"`
	CategoryID  int64   `json:"categoryId"`
	TagIds      []int64 `json:"tagIds"`
}

// ChangeStatusRequest represents the request for changing product status
type ChangeStatusRequest struct {
	Action string `json:"action" binding:"required,oneof=delist relist sold"`
}

// SearchParams represents the search parameters for products
type SearchParams struct {
	Keyword     string `form:"keyword"`
	CategoryID  int64  `form:"categoryId"`
	MinPrice    int64  `form:"minPrice"`
	MaxPrice    int64  `form:"maxPrice"`
	ConditionID int64  `form:"conditionId"`
	TagID       int64  `form:"tagId"`
	SortBy      string `form:"sortBy" binding:"omitempty,oneof=created_at price"`
	SortOrder   string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
	Page        int    `form:"page" binding:"omitempty,min=1"`
	PageSize    int    `form:"pageSize" binding:"omitempty,min=1,max=50"`
}
