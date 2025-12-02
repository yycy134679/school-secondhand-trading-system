package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"gorm.io/gorm"
)

// ProductRepository defines product data access behavior
type ProductRepository interface {
	// 5.2.1 Create: 事务插入 products、product_images、product_tags，设置 main_image_url
	Create(ctx context.Context, product *model.Product, images []model.ProductImage, tagIDs []int64) (int64, error)

	// 5.2.2 Update: 普通卖家禁止编辑 Sold 商品；管理员在不改 status 的前提下可编辑其他字段
	Update(ctx context.Context, product *model.Product, images []model.ProductImage, tagIDs []int64, isAdmin bool) error

	// 5.2.3 GetByID: 联合加载图片与 tags
	GetByID(ctx context.Context, id int64) (*model.Product, []model.ProductImage, []int64, error)

	// 5.2.4 ListBySeller: 用于"我发布的"列表
	ListBySeller(ctx context.Context, sellerID int64, keyword string, page, pageSize int) ([]model.Product, int64, error)

	// 5.2.5 UpdateStatus: 带 where 条件，依赖 DB 触发器防非法流转
	UpdateStatus(ctx context.Context, id int64, fromStatus, toStatus string) error

	// 5.2.6 Search: 实现关键词 + 条件组合搜索，仅 status=ForSale
	Search(ctx context.Context, params model.SearchParams) ([]model.Product, int64, error)

	// 5.2.7 ListLatestForSale: 供首页使用
	ListLatestForSale(ctx context.Context, excludeIDs []int64, page, pageSize int) ([]model.Product, int64, error)

	// 5.2.7 ListByCategory: 供分类页使用
	ListByCategory(ctx context.Context, categoryID int64, params model.SearchParams) ([]model.Product, int64, error)

	// 图片管理相关方法
	CreateImage(ctx context.Context, image *model.ProductImage) (int64, error)
	GetImageByID(ctx context.Context, imageID int64) (*model.ProductImage, error)
	HasPrimaryImage(ctx context.Context, productID int64) (bool, error)
	SetPrimaryImage(ctx context.Context, productID, imageID int64) error
	UpdateImageSort(ctx context.Context, imageID int64, sortOrder int) error
	DeleteImage(ctx context.Context, imageID int64) error
	GetFirstImage(ctx context.Context, productID int64) (*model.ProductImage, error)
	UpdateMainImageURL(ctx context.Context, productID int64, url string) error
}

// productRepoImpl implements ProductRepository
type productRepoImpl struct {
	// 数据库连接
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository instance
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepoImpl{
		db: db,
	}
}

// Create: 事务插入 products、product_images、product_tags，设置 main_image_url
func (r *productRepoImpl) Create(ctx context.Context, product *model.Product, images []model.ProductImage, tagIDs []int64) (int64, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 插入商品
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		// 2. 插入图片
		for i := range images {
			images[i].ProductID = product.ID
		}
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}

			// 3. 设置主图URL
			for _, img := range images {
				if img.IsPrimary {
					product.MainImageURL = img.URL
					break
				}
			}
			// 如果没有设置主图，默认第一张为主图
			if product.MainImageURL == "" && len(images) > 0 {
				product.MainImageURL = images[0].URL
				if err := tx.Model(&model.ProductImage{}).Where("id = ?", images[0].ID).Update("is_primary", true).Error; err != nil {
					return err
				}
			}
			// 更新商品的主图URL
			if err := tx.Model(product).Update("main_image_url", product.MainImageURL).Error; err != nil {
				return err
			}
		}

		// 4. 插入标签关联
		if len(tagIDs) > 0 {
			for _, tagID := range tagIDs {
				if err := tx.Exec("INSERT INTO product_tags (product_id, tag_id) VALUES (?, ?)", product.ID, tagID).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
	return product.ID, err
}

// Update: 普通卖家禁止编辑 Sold 商品；管理员在不改 status 的前提下可编辑其他字段
func (r *productRepoImpl) Update(ctx context.Context, product *model.Product, images []model.ProductImage, tagIDs []int64, isAdmin bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查商品是否存在
		existingProduct := &model.Product{}
		if err := tx.First(existingProduct, product.ID).Error; err != nil {
			return err
		}

		// 检查权限
		if !isAdmin && existingProduct.Status == "Sold" {
			return errors.New("普通卖家禁止编辑已售商品")
		}

		// 如果不是管理员，禁止修改status
		if !isAdmin {
			product.Status = existingProduct.Status
		}

		// 更新商品基本信息
		if err := tx.Model(product).Updates(product).Error; err != nil {
			return err
		}

		// 处理图片（这里简化处理，实际可能需要更复杂的逻辑）
		if len(images) > 0 {
			// 删除原有图片
			if err := tx.Where("product_id = ?", product.ID).Delete(&model.ProductImage{}).Error; err != nil {
				return err
			}

			// 插入新图片
			for i := range images {
				images[i].ProductID = product.ID
			}
			if err := tx.Create(&images).Error; err != nil {
				return err
			}

			// 更新主图URL
			mainImageURL := ""
			for _, img := range images {
				if img.IsPrimary {
					mainImageURL = img.URL
					break
				}
			}
			if mainImageURL != "" {
				if err := tx.Model(product).Update("main_image_url", mainImageURL).Error; err != nil {
					return err
				}
			}
		}

		// 处理标签
		// 删除原有标签关联
		if err := tx.Exec("DELETE FROM product_tags WHERE product_id = ?", product.ID).Error; err != nil {
			return err
		}

		// 插入新标签关联
		if len(tagIDs) > 0 {
			for _, tagID := range tagIDs {
				if err := tx.Exec("INSERT INTO product_tags (product_id, tag_id) VALUES (?, ?)", product.ID, tagID).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetByID: 联合加载图片与 tags
func (r *productRepoImpl) GetByID(ctx context.Context, id int64) (*model.Product, []model.ProductImage, []int64, error) {
	product := &model.Product{}
	if err := r.db.WithContext(ctx).First(product, id).Error; err != nil {
		return nil, nil, nil, err
	}

	// 加载图片
	var images []model.ProductImage
	if err := r.db.WithContext(ctx).Where("product_id = ?", id).Order("sort_order ASC").Find(&images).Error; err != nil {
		return nil, nil, nil, err
	}

	// 加载标签ID
	var tagIDs []int64
	rows, err := r.db.WithContext(ctx).Raw("SELECT tag_id FROM product_tags WHERE product_id = ?", id).Rows()
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tagID int64
		if err := rows.Scan(&tagID); err != nil {
			return nil, nil, nil, err
		}
		tagIDs = append(tagIDs, tagID)
	}

	return product, images, tagIDs, nil
}

// ListBySeller: 用于"我发布的"列表
func (r *productRepoImpl) ListBySeller(ctx context.Context, sellerID int64, keyword string, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("seller_id = ?", sellerID)

	// 关键词搜索
	if keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateStatus: 带 where 条件，依赖 DB 触发器防非法流转
func (r *productRepoImpl) UpdateStatus(ctx context.Context, id int64, fromStatus, toStatus string) error {
	// 带条件更新，确保只在状态匹配时更新
	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ? AND status = ?", id, fromStatus).Update("status", toStatus)
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的更新了记录
	if result.RowsAffected == 0 {
		return errors.New("product not found or status mismatch")
	}

	return nil
}

// Search: 实现关键词 + 条件组合搜索，仅 status=ForSale
func (r *productRepoImpl) Search(ctx context.Context, params model.SearchParams) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "ForSale")

	// 关键词搜索
	if params.Keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// 分类过滤
	if params.CategoryID != 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}

	// 价格范围
	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}

	// 新旧程度
	if params.ConditionID > 0 {
		query = query.Where("condition_id = ?", params.ConditionID)
	}

	// 标签过滤（需要JOIN）
	if params.TagID > 0 {
		query = query.Joins("JOIN product_tags ON products.id = product_tags.product_id").
			Where("product_tags.tag_id = ?", params.TagID).
			Group("products.id")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 根据排序参数设置排序方式
	if params.SortBy != "" {
		order := params.SortBy + " " + params.SortOrder
		if params.SortOrder == "" {
			order = params.SortBy + " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC") // 默认最新
	}

	// 设置默认分页参数
	page := params.Page
	pageSize := params.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListLatestForSale: 供首页使用
func (r *productRepoImpl) ListLatestForSale(ctx context.Context, excludeIDs []int64, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "ForSale")

	// 排除指定ID
	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListByCategory: 供分类页使用
func (r *productRepoImpl) ListByCategory(ctx context.Context, categoryID int64, params model.SearchParams) ([]model.Product, int64, error) {
	// 设置分类ID
	params.CategoryID = categoryID
	return r.Search(ctx, params)
}

// CreateImage 创建商品图片记录
func (r *productRepoImpl) CreateImage(ctx context.Context, image *model.ProductImage) (int64, error) {
	result := r.db.WithContext(ctx).Create(image)
	if result.Error != nil {
		return 0, fmt.Errorf("创建图片记录失败: %w", result.Error)
	}
	return image.ID, nil
}

// GetImageByID 根据ID获取图片信息
func (r *productRepoImpl) GetImageByID(ctx context.Context, imageID int64) (*model.ProductImage, error) {
	var image model.ProductImage
	result := r.db.WithContext(ctx).First(&image, imageID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &image, nil
}

// HasPrimaryImage 检查商品是否已有主图
func (r *productRepoImpl) HasPrimaryImage(ctx context.Context, productID int64) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.ProductImage{}).Where("product_id = ? AND is_primary = ?", productID, true).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("检查主图失败: %w", result.Error)
	}
	return count > 0, nil
}

// SetPrimaryImage 设置主图
func (r *productRepoImpl) SetPrimaryImage(ctx context.Context, productID, imageID int64) error {
	// 使用事务
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 将所有图片设为非主图
	if err := tx.Model(&model.ProductImage{}).Where("product_id = ?", productID).Update("is_primary", false).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("重置主图失败: %w", err)
	}

	// 2. 将指定图片设为主图
	if err := tx.Model(&model.ProductImage{}).Where("id = ? AND product_id = ?", imageID, productID).Update("is_primary", true).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("设置主图失败: %w", err)
	}

	return tx.Commit().Error
}

// UpdateImageSort 更新图片排序
func (r *productRepoImpl) UpdateImageSort(ctx context.Context, imageID int64, sortOrder int) error {
	result := r.db.WithContext(ctx).Model(&model.ProductImage{}).Where("id = ?", imageID).Update("sort_order", sortOrder)
	if result.Error != nil {
		return fmt.Errorf("更新排序失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("图片不存在")
	}
	return nil
}

// DeleteImage 删除图片
func (r *productRepoImpl) DeleteImage(ctx context.Context, imageID int64) error {
	result := r.db.WithContext(ctx).Delete(&model.ProductImage{}, imageID)
	if result.Error != nil {
		return fmt.Errorf("删除图片失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("图片不存在")
	}
	return nil
}

// GetFirstImage 获取商品的第一张图片
func (r *productRepoImpl) GetFirstImage(ctx context.Context, productID int64) (*model.ProductImage, error) {
	var image model.ProductImage
	result := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("sort_order ASC, id ASC").First(&image)
	if result.Error != nil {
		return nil, result.Error
	}
	return &image, nil
}

// UpdateMainImageURL 更新商品主图URL
func (r *productRepoImpl) UpdateMainImageURL(ctx context.Context, productID int64, url string) error {
	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", productID).Update("main_image_url", url)
	if result.Error != nil {
		return fmt.Errorf("更新主图URL失败: %w", result.Error)
	}
	return nil
}
