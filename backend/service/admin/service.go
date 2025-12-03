package admin

import (
	"context"
	"fmt"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"gorm.io/gorm"
)

// UserListResponse 用户列表响应结构
type UserListResponse struct {
	Total int64         `json:"total"` // 总用户数
	Users []*model.User `json:"users"` // 用户列表
}

// AdminService 管理后台服务接口
type AdminService struct {
	db *gorm.DB
}

// NewAdminService 创建管理后台服务实例
func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		db: db,
	}
}

// DashboardStats 仪表盘统计数据结构
type DashboardStats struct {
	UserCount    int64 `json:"userCount"`    // 用户总数
	ProductCount int64 `json:"productCount"` // 商品总数
	ForSaleCount int64 `json:"forSaleCount"` // 在售商品数
	SoldCount    int64 `json:"soldCount"`    // 已售商品数
}

// ListUsers 获取用户列表，支持根据账号/昵称模糊搜索和分页
//
// 功能说明：
//   - 根据关键词对用户账号和昵称进行模糊搜索
//   - 支持分页功能，返回用户总数和当前页的用户列表
//
// 参数：
//   - ctx: 上下文，用于控制查询超时
//   - keyword: 搜索关键词，根据账号或昵称进行模糊匹配
//   - page: 页码，从1开始
//   - pageSize: 每页数量
//
// 返回值：
//   - UserListResponse: 包含用户总数和用户列表的结构体
//   - error: 错误信息，查询失败时返回
func (s *AdminService) ListUsers(ctx context.Context, keyword string, page, pageSize int) (*UserListResponse, error) {
	var users []*model.User
	var total int64

	// 构建查询
	query := s.db.WithContext(ctx).Model(&model.User{})

	// 如果有关键词，添加模糊搜索条件
	if keyword != "" {
		query = query.Where("account LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	return &UserListResponse{
		Total: total,
		Users: users,
	}, nil
}

// GetDashboardStats 获取仪表盘统计数据
//
// 功能说明：
//   - 统计系统中的用户总数、商品总数、在售商品数、已售商品数
//   - 使用COUNT聚合查询提高性能
//
// 参数：
//   - ctx: 上下文，用于控制查询超时
//
// 返回值：
//   - DashboardStats: 包含统计数据的结构体
//   - error: 错误信息，查询失败时返回
func (s *AdminService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	var stats DashboardStats

	// 1. 统计用户总数
	if err := s.db.WithContext(ctx).Model(&model.User{}).Count(&stats.UserCount).Error; err != nil {
		return nil, err
	}

	// 2. 统计商品总数
	if err := s.db.WithContext(ctx).Model(&model.Product{}).Count(&stats.ProductCount).Error; err != nil {
		return nil, err
	}

	// 3. 统计在售商品数
	if err := s.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "for_sale").Count(&stats.ForSaleCount).Error; err != nil {
		return nil, err
	}

	// 4. 统计已售商品数
	if err := s.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "sold").Count(&stats.SoldCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// ProductListResponse 商品列表响应数据结构
type ProductListResponse struct {
	Total    int64             `json:"total"`
	Products []AdminProductDTO `json:"products"`
}

// AdminProductDTO 管理后台商品数据传输对象
type AdminProductDTO struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Price          int64  `json:"price"`
	Status         string `json:"status"`
	SellerID       int64  `json:"sellerId"`
	SellerAccount  string `json:"sellerAccount"`
	SellerNickname string `json:"sellerNickname"`
	CategoryID     int64  `json:"categoryId"`
	CategoryName   string `json:"categoryName"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// ListProductsAdmin 获取商品列表，支持状态/卖家ID/关键词过滤
func (s *AdminService) ListProductsAdmin(ctx context.Context, status string, sellerId int64, keyword string, page, pageSize int) (*ProductListResponse, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询条件
	baseQuery := `SELECT 
		p.id, p.title, p.price, p.status, p.seller_id, 
		u.account as seller_account, u.nickname as seller_nickname,
		p.category_id, c.name as category_name,
		p.created_at, p.updated_at
	FROM products p
	LEFT JOIN users u ON p.seller_id = u.id
	LEFT JOIN categories c ON p.category_id = c.id`

	countQuery := `SELECT COUNT(*) FROM products p`
	whereClause := ""
	args := []interface{}{}

	// 添加过滤条件
	if status != "" || sellerId > 0 || keyword != "" {
		whereClause = " WHERE"
		if status != "" {
			whereClause += " p.status = ?"
			args = append(args, status)
		}
		if sellerId > 0 {
			if status != "" {
				whereClause += " AND"
			}
			whereClause += " p.seller_id = ?"
			args = append(args, sellerId)
		}
		if keyword != "" {
			if status != "" || sellerId > 0 {
				whereClause += " AND"
			}
			whereClause += " (p.title LIKE ? OR p.description LIKE ?)"
			likeKeyword := "%" + keyword + "%"
			args = append(args, likeKeyword, likeKeyword)
		}
	}

	// 查询总数
	var total int64
	if err := s.db.WithContext(ctx).Raw(countQuery+whereClause, args...).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("查询商品总数失败: %w", err)
	}

	// 查询商品列表
	query := baseQuery + whereClause + " ORDER BY p.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	products := make([]AdminProductDTO, 0)
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&products).Error; err != nil {
		return nil, fmt.Errorf("查询商品列表失败: %w", err)
	}

	return &ProductListResponse{
		Total:    total,
		Products: products,
	}, nil
}

// UpdateProductRequest 管理后台更新商品请求结构
type UpdateProductRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       int64   `json:"price"`
	ConditionID int64   `json:"conditionId"`
	CategoryID  int64   `json:"categoryId"`
	TagIDs      []int64 `json:"tagIds"`
	Status      *string `json:"status,omitempty"` // 可选字段，但会被禁止修改
}

// UpdateProductAsAdmin 管理员更新商品，禁止修改status字段
// 如果请求体携带status或试图改变状态，返回3004错误
func (s *AdminService) UpdateProductAsAdmin(ctx context.Context, productID int64, req UpdateProductRequest) error {
	// 检查请求是否携带status字段
	if req.Status != nil {
		return fmt.Errorf("3004:禁止修改商品状态字段")
	}

	// 开始事务
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %w", tx.Error)
	}
	defer tx.Rollback()

	// 检查商品是否存在
	var existingProduct model.Product
	query := `SELECT id, status FROM products WHERE id = ?`
	if err := tx.WithContext(ctx).Raw(query, productID).Scan(&existingProduct).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("商品不存在")
		}
		return fmt.Errorf("查询商品信息失败: %w", err)
	}

	// 更新商品基本信息（排除status字段）
	updateQuery := `UPDATE products 
		SET title = ?, description = ?, price = ?, 
		    condition_id = ?, category_id = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	var err error
	if err = tx.WithContext(ctx).Exec(updateQuery,
		req.Title, req.Description, req.Price,
		req.ConditionID, req.CategoryID, productID).Error; err != nil {
		return fmt.Errorf("更新商品基本信息失败: %s", err)
	}

	// 更新标签关联
	// 先删除现有标签关联
	if err = tx.WithContext(ctx).Exec("DELETE FROM product_tags WHERE product_id = ?", productID).Error; err != nil {
		return fmt.Errorf("删除现有标签关联失败: %s", err)
	}

	// 添加新的标签关联
	if len(req.TagIDs) > 0 {
		insertQuery := "INSERT INTO product_tags (product_id, tag_id) VALUES (?, ?)"
		for _, tagID := range req.TagIDs {
			if err = tx.WithContext(ctx).Exec(insertQuery, productID, tagID).Error; err != nil {
				return fmt.Errorf("添加标签关联失败: %w", err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}
