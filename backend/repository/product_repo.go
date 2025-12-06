package repository

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
)

// ProductRepository 商品仓库接口
type ProductRepository interface {
	Create(ctx context.Context, product *model.Product, images []model.ProductImage, tagIDs []int64) (int64, error)
	Update(ctx context.Context, product *model.Product, images []model.ProductImage, tagIDs []int64, isAdmin bool) error
	GetByID(ctx context.Context, id int64) (*model.Product, []model.ProductImage, []int64, error)
	ListBySeller(ctx context.Context, sellerID int64, keyword string, page, pageSize int) ([]model.Product, int64, error)
	UpdateStatus(ctx context.Context, id int64, fromStatus, toStatus string) error
	Search(ctx context.Context, params SearchParams) ([]model.Product, int64, error)
	ListLatestForSale(ctx context.Context, excludeIDs []int64, page, pageSize int) ([]model.Product, int64, error)
	ListByCategory(ctx context.Context, categoryID int64, params SearchParams) ([]model.Product, int64, error)
}

// SearchParams 搜索参数
type SearchParams struct {
	Keyword     string
	PriceMin    float64
	PriceMax    float64
	ConditionID int64
	Page        int
	PageSize    int
}

// productRepository 商品仓库实现
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建商品仓库实例
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

// Create 在事务中创建商品，包括商品基本信息、图片和标签
func (r *productRepository) Create(ctx context.Context, product *model.Product, images []model.ProductImage, tagIDs []int64) (int64, error) {
	if r.db == nil {
		return 0, fmt.Errorf("db is nil")
	}

	var createdID int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 插入商品基本信息
		if err := tx.Create(product).Error; err != nil {
			return fmt.Errorf("create product failed: %w", err)
		}
		createdID = product.ID

		// 设置主图URL并插入图片
		mainImageURL := ""
		if len(images) > 0 {
			primaryExists := false
			for i := range images {
				images[i].ProductID = product.ID
				if images[i].IsPrimary {
					primaryExists = true
					mainImageURL = images[i].URL
				}
			}
			if !primaryExists {
				images[0].IsPrimary = true
				mainImageURL = images[0].URL
			}

			if err := tx.Create(&images).Error; err != nil {
				return fmt.Errorf("create product images failed: %w", err)
			}

			if err := tx.Model(product).Update("main_image_url", mainImageURL).Error; err != nil {
				return fmt.Errorf("update main image url failed: %w", err)
			}
		}

		// 处理标签关联
		if len(tagIDs) > 0 {
			relations := make([]map[string]interface{}, 0, len(tagIDs))
			for _, tagID := range tagIDs {
				relations = append(relations, map[string]interface{}{
					"product_id": product.ID,
					"tag_id":     tagID,
				})
			}
			if err := tx.Table("product_tags").Create(&relations).Error; err != nil {
				return fmt.Errorf("create product tag relation failed: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return createdID, nil
}

// Update 更新商品信息，包含权限控制
func (r *productRepository) Update(ctx context.Context, product *model.Product, images []model.ProductImage, tagIDs []int64, isAdmin bool) error {
	// 开始事务
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("begin transaction failed: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取原商品信息进行权限检查
	var originalProduct model.Product
	if err := tx.First(&originalProduct, product.ID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("get product failed: %w", err)
	}

	// 权限检查：普通卖家禁止编辑已售出商品
	if !isAdmin && originalProduct.Status == "Sold" {
		tx.Rollback()
		return fmt.Errorf("cannot update sold product")
	}

	// 管理员在不改变status的前提下可以编辑已售出商品的其他字段
	if isAdmin && originalProduct.Status == "Sold" && product.Status != "Sold" {
		tx.Rollback()
		return fmt.Errorf("admin cannot change status of sold product")
	}

	// 更新商品基本信息
	if err := tx.Model(product).Updates(product).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("update product failed: %w", err)
	}

	// 处理图片更新（这里假设images包含了所有需要保留的图片，会替换原有图片）
	if len(images) > 0 {
		// 删除原有图片
		if err := tx.Where("product_id = ?", product.ID).Delete(&model.ProductImage{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("delete old images failed: %w", err)
		}

		// 设置主图URL
		mainImageURL := ""
		primaryExists := false
		for i := range images {
			images[i].ProductID = product.ID
			if images[i].IsPrimary {
				primaryExists = true
				mainImageURL = images[i].URL
			}
		}
		// 如果没有指定主图，将第一张设为主图
		if !primaryExists && len(images) > 0 {
			images[0].IsPrimary = true
			mainImageURL = images[0].URL
		}

		// 插入新图片
		if err := tx.CreateInBatches(images, len(images)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("create new images failed: %w", err)
		}

		// 更新主图URL
		if mainImageURL != "" {
			if err := tx.Model(product).Update("main_image_url", mainImageURL).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("update main image url failed: %w", err)
			}
		}
	}

	// 更新标签关联
	// 删除原有标签关联
	if err := tx.Table("product_tags").Where("product_id = ?", product.ID).Delete(nil).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("delete old tags relation failed: %w", err)
	}

	// 创建新标签关联
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			productTag := gin.H{
				"product_id": product.ID,
				"tag_id":     tagID,
			}
			if err := tx.Table("product_tags").Create(productTag).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("create new tags relation failed: %w", err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}

// GetByID 根据ID获取商品详情，联合加载图片与tags
func (r *productRepository) GetByID(ctx context.Context, id int64) (*model.Product, []model.ProductImage, []int64, error) {
	// 查询商品基本信息
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("get product failed: %w", err)
	}

	// 查询商品图片
	var images []model.ProductImage
	if err := r.db.WithContext(ctx).Where("product_id = ?", id).Order("sort_order ASC").Find(&images).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("get product images failed: %w", err)
	}

	// 查询商品标签ID
	type TagRelation struct {
		TagID int64 `gorm:"column:tag_id"`
	}
	var tagRelations []TagRelation
	if err := r.db.WithContext(ctx).Table("product_tags").Select("tag_id").Where("product_id = ?", id).Find(&tagRelations).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("get product tags failed: %w", err)
	}

	// 提取标签ID列表
	tagIDs := make([]int64, len(tagRelations))
	for i, relation := range tagRelations {
		tagIDs[i] = relation.TagID
	}

	return &product, images, tagIDs, nil
}

// ListBySeller 获取卖家发布的商品列表，支持关键词搜索和分页
func (r *productRepository) ListBySeller(ctx context.Context, sellerID int64, keyword string, page, pageSize int) ([]model.Product, int64, error) {
	// 构建查询
	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("seller_id = ?", sellerID)

	// 添加关键词搜索
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products failed: %w", err)
	}

	// 分页查询
	var products []model.Product
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("list products failed: %w", err)
	}

	return products, total, nil
}

// UpdateStatus 更新商品状态，带where条件确保状态流转的合法性
// 依赖数据库触发器防止非法流转
func (r *productRepository) UpdateStatus(ctx context.Context, id int64, fromStatus, toStatus string) error {
	// 使用where条件确保只有当前状态为fromStatus的商品才会被更新
	// 这样可以在数据库层面保证状态流转的合法性
	result := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Update("status", toStatus)

	// 检查是否有行被更新
	if result.Error != nil {
		return fmt.Errorf("update product status failed: %w", result.Error)
	}

	// 如果没有行被更新，说明状态不匹配或者商品不存在
	if result.RowsAffected == 0 {
		// 检查商品是否存在
		var count int64
		if err := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Count(&count).Error; err != nil {
			return fmt.Errorf("check product existence failed: %w", err)
		}

		if count == 0 {
			return fmt.Errorf("product not found")
		}

		// 商品存在但状态不匹配
		return fmt.Errorf("invalid status transition")
	}

	return nil
}

// Search 实现关键词+条件组合搜索，仅status=ForSale
func (r *productRepository) Search(ctx context.Context, params SearchParams) ([]model.Product, int64, error) {
	// 构建查询
	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "ForSale")

	// 添加搜索条件
	if params.Keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.ConditionID > 0 {
		query = query.Where("condition_id = ?", params.ConditionID)
	}

	if params.PriceMin > 0 {
		query = query.Where("price >= ?", params.PriceMin)
	}

	if params.PriceMax > 0 {
		query = query.Where("price <= ?", params.PriceMax)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count search products failed: %w", err)
	}

	// 分页查询
	var products []model.Product
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("search products failed: %w", err)
	}

	return products, total, nil
}

// ListLatestForSale 获取最新上架的商品，可排除指定ID
func (r *productRepository) ListLatestForSale(ctx context.Context, excludeIDs []int64, page, pageSize int) ([]model.Product, int64, error) {
	// 构建查询
	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "ForSale")

	// 添加排除条件
	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN (?)", excludeIDs)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count latest products failed: %w", err)
	}

	// 分页查询
	var products []model.Product
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("list latest products failed: %w", err)
	}

	return products, total, nil
}

// ListByCategory 获取指定分类的商品
func (r *productRepository) ListByCategory(ctx context.Context, categoryID int64, params SearchParams) ([]model.Product, int64, error) {
	// 构建查询
	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("status = ? AND category_id = ?", "ForSale", categoryID)

	// 添加搜索条件
	if params.Keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.ConditionID > 0 {
		query = query.Where("condition_id = ?", params.ConditionID)
	}

	if params.PriceMin > 0 {
		query = query.Where("price >= ?", params.PriceMin)
	}

	if params.PriceMax > 0 {
		query = query.Where("price <= ?", params.PriceMax)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count category products failed: %w", err)
	}

	// 分页查询
	var products []model.Product
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("list category products failed: %w", err)
	}

	return products, total, nil
}
