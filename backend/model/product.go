package model

import "time"

// Product represents a product listing.
type Product struct {
	ID          int64     `json:"id"`
	SellerID    int64     `json:"seller_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
