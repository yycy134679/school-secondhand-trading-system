package product

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/cache"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/util"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
)

// ProductService 商品服务结构体
type ProductService struct {
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	db          *gorm.DB
	cache       *cache.MemoryCache
}

// NewProductService 创建商品服务实例
func NewProductService(
	db *gorm.DB,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	cache *cache.MemoryCache,
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		userRepo:    userRepo,
		db:          db,
		cache:       cache,
	}
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Title       string
	Description string
	Price       float64
	ConditionID int64
	CategoryID  int64
	TagIDs      []int64
	Images      []*multipart.FileHeader
	// PrimaryImageIndex 前端标记的主图索引，可为空
	PrimaryImageIndex *int
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	CategoryID  *int64   `json:"categoryId"`
	ConditionID *int64   `json:"conditionId"`
	// 前端以逗号分隔字符串传递
	TagIDs    string   `json:"tagIds"`
	ImageURLs []string `json:"imageUrls"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword     string
	PriceMin    *float64
	PriceMax    *float64
	ConditionID *int64
	Page        int
	PageSize    int
}

// SearchParams 搜索参数结构体（与controller中使用的名称保持一致）
type SearchParams struct {
	Keyword      string
	MinPrice     *float64
	MaxPrice     *float64
	ConditionID  *int64
	ConditionIDs []int64
	CategoryID   *int64
	TagID        *int64
	Sort         string
	Page         int
	PageSize     int
}

type statusChangeRecord struct {
	From   string
	To     string
	UserID int64
}

// ContactResponse 联系卖家的返回结构
type ContactResponse struct {
	CanContact   bool    `json:"canContact"`
	SellerWechat *string `json:"sellerWechat"`
	Tips         string  `json:"tips,omitempty"`
}

const (
	detailCacheTTL = 5 * time.Minute
)

// CreateProduct 创建商品
func (s *ProductService) CreateProduct(ctx context.Context, userID int64, req *CreateProductRequest) (interface{}, error) {
	if s.productRepo == nil || s.userRepo == nil {
		return nil, fmt.Errorf("服务未初始化")
	}

	// 检查用户及微信号
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if strings.TrimSpace(user.WechatID) == "" {
		return nil, fmt.Errorf("请先完善微信号")
	}

	// 至少一张图片
	if len(req.Images) == 0 {
		return nil, fmt.Errorf("请至少上传一张图片")
	}

	primaryIndex := 0
	if req.PrimaryImageIndex != nil && *req.PrimaryImageIndex >= 0 && *req.PrimaryImageIndex < len(req.Images) {
		primaryIndex = *req.PrimaryImageIndex
	}

	// 保存图片文件
	images := make([]model.ProductImage, 0, len(req.Images))
	for i, header := range req.Images {
		file, err := header.Open()
		if err != nil {
			return nil, fmt.Errorf("读取图片失败: %w", err)
		}
		url, err := util.SaveImage(file, header)
		file.Close()
		if err != nil {
			return nil, err
		}

		images = append(images, model.ProductImage{
			URL:       url,
			SortOrder: i + 1,
			IsPrimary: i == primaryIndex,
		})
	}

	product := &model.Product{
		Title:        req.Title,
		Description:  req.Description,
		Price:        req.Price,
		CategoryID:   req.CategoryID,
		ConditionID:  req.ConditionID,
		SellerID:     userID,
		Status:       "ForSale",
		MainImageURL: images[primaryIndex].URL,
	}

	if _, err := s.productRepo.Create(ctx, product, images, req.TagIDs); err != nil {
		return nil, err
	}

	// 确保响应包含主图
	product.MainImageURL = images[primaryIndex].URL

	dto, err := s.buildDetailDTO(ctx, product, images, req.TagIDs, &userID)
	if err == nil && s.cache != nil {
		_ = s.cache.Set(ctx, buildDetailCacheKey(product.ID), dto, detailCacheTTL)
	}

	return product, nil
}

// UpdateProduct 更新商品
func (s *ProductService) UpdateProduct(ctx context.Context, userID, productID int64, req *UpdateProductRequest, isAdmin bool) (interface{}, error) {
	if s.productRepo == nil || s.userRepo == nil {
		return nil, fmt.Errorf("服务未初始化")
	}

	product, images, tagIDs, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("商品不存在")
		}
		return nil, err
	}

	if !isAdmin && product.SellerID != userID {
		return nil, fmt.Errorf("无权限操作该商品")
	}

	if product.Status == "Sold" && !isAdmin {
		return nil, fmt.Errorf("已售出的商品不能修改")
	}

	if req.Title != nil {
		product.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		product.Description = strings.TrimSpace(*req.Description)
	}
	if req.Price != nil && *req.Price > 0 {
		product.Price = *req.Price
	}
	if req.CategoryID != nil && *req.CategoryID > 0 {
		product.CategoryID = *req.CategoryID
	}
	if req.ConditionID != nil && *req.ConditionID > 0 {
		product.ConditionID = *req.ConditionID
	}

	if req.TagIDs != "" {
		tagIDs = parseTagIDs(req.TagIDs)
	}

	if len(req.ImageURLs) > 0 {
		images = make([]model.ProductImage, 0, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			images = append(images, model.ProductImage{
				ProductID: product.ID,
				URL:       url,
				SortOrder: i + 1,
				IsPrimary: i == 0,
			})
		}
		product.MainImageURL = images[0].URL
	}

	if err := s.productRepo.Update(ctx, product, images, tagIDs, isAdmin); err != nil {
		return nil, err
	}

	return product, nil
}

// ChangeStatus 变更商品状态
func (s *ProductService) ChangeStatus(ctx context.Context, userID, productID int64, action string) error {
	if s.productRepo == nil || s.userRepo == nil {
		return fmt.Errorf("服务未初始化")
	}

	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("商品不存在")
		}
		return err
	}

	if product.SellerID != userID {
		return fmt.Errorf("无权限操作该商品")
	}

	fromStatus := product.Status
	toStatus := ""
	switch action {
	case "delist":
		if fromStatus != "ForSale" {
			return fmt.Errorf("终态或状态不匹配，无法下架")
		}
		toStatus = "Delisted"
	case "relist":
		if fromStatus != "Delisted" {
			return fmt.Errorf("状态不匹配，无法重新上架")
		}
		toStatus = "ForSale"
	case "sold":
		if fromStatus != "ForSale" {
			return fmt.Errorf("终态或状态不匹配，无法标记已售")
		}
		toStatus = "Sold"
	default:
		return fmt.Errorf("无效的动作")
	}

	if err := s.productRepo.UpdateStatus(ctx, productID, fromStatus, toStatus); err != nil {
		return err
	}

	if s.cache != nil {
		record := statusChangeRecord{
			From:   fromStatus,
			To:     toStatus,
			UserID: userID,
		}
		_ = s.cache.Set(ctx, buildStatusCacheKey(productID), record, 3*time.Second)
	}

	return nil
}

// UndoLastStatusChange 撤销状态变更
func (s *ProductService) UndoLastStatusChange(ctx context.Context, userID, productID int64) error {
	if s.productRepo == nil || s.cache == nil {
		return fmt.Errorf("撤销记录不存在或超时")
	}

	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("商品不存在")
		}
		return err
	}

	if product.SellerID != userID {
		return fmt.Errorf("无权限操作该商品")
	}

	val, err := s.cache.Get(ctx, buildStatusCacheKey(productID))
	if err != nil {
		return fmt.Errorf("撤销记录不存在或超时")
	}

	record, ok := val.(statusChangeRecord)
	if !ok || record.UserID != userID {
		return fmt.Errorf("撤销记录不存在或超时")
	}

	if product.Status != record.To {
		return fmt.Errorf("撤销记录不存在或超时")
	}

	if err := s.productRepo.UpdateStatus(ctx, productID, record.To, record.From); err != nil {
		return err
	}

	_ = s.cache.Delete(ctx, buildStatusCacheKey(productID))
	return nil
}

// GetProductDetail 获取商品详情
func (s *ProductService) GetProductDetail(ctx context.Context, productID int64, viewerID *int64) (*model.ProductDetailDTO, error) {
	if s.productRepo == nil || s.userRepo == nil {
		return nil, fmt.Errorf("服务未初始化")
	}

	if s.cache != nil {
		if val, err := s.cache.Get(ctx, buildDetailCacheKey(productID)); err == nil {
			if dto, ok := val.(*model.ProductDetailDTO); ok && dto != nil {
				return dto, nil
			}
		}
	}

	product, images, tagIDs, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("商品不存在")
		}
		return nil, err
	}

	dto, err := s.buildDetailDTO(ctx, product, images, tagIDs, viewerID)
	if err == nil && s.cache != nil {
		_ = s.cache.Set(ctx, buildDetailCacheKey(productID), dto, detailCacheTTL)
	}

	return dto, nil
}

// GetProductContact 获取联系卖家信息（微信号或提示）
func (s *ProductService) GetProductContact(ctx context.Context, productID int64, viewerID *int64) (*ContactResponse, error) {
	if s.productRepo == nil || s.userRepo == nil {
		return nil, fmt.Errorf("服务未初始化")
	}

	if viewerID == nil {
		return &ContactResponse{
			CanContact:   false,
			SellerWechat: nil,
			Tips:         "请先登录后联系卖家",
		}, nil
	}

	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("商品不存在")
		}
		return nil, err
	}

	if product.SellerID == *viewerID {
		return &ContactResponse{
			CanContact:   false,
			SellerWechat: nil,
		}, nil
	}

	if product.Status != "ForSale" {
		return &ContactResponse{
			CanContact:   false,
			SellerWechat: nil,
			Tips:         "商品当前不可联系卖家",
		}, nil
	}

	seller, err := s.userRepo.GetByID(ctx, product.SellerID)
	if err != nil {
		return nil, fmt.Errorf("获取卖家信息失败: %w", err)
	}

	wechat := strings.TrimSpace(seller.WechatID)
	if wechat == "" {
		return &ContactResponse{
			CanContact:   false,
			SellerWechat: nil,
			Tips:         "卖家联系方式暂不可用，请稍后再试",
		}, nil
	}

	return &ContactResponse{
		CanContact:   true,
		SellerWechat: &wechat,
		Tips:         "请线下交易，注意安全",
	}, nil
}

// ListMyProducts 获取我的商品列表
func (s *ProductService) ListMyProducts(ctx context.Context, userID int64, keyword string, page, pageSize int) ([]model.ProductCardDTO, int64, error) {
	if s.productRepo == nil {
		return nil, 0, fmt.Errorf("服务未初始化")
	}

	products, total, err := s.productRepo.ListBySeller(ctx, userID, keyword, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	dtoList := make([]model.ProductCardDTO, 0, len(products))
	for i := range products {
		dto, err := s.toCardDTO(ctx, &products[i])
		if err != nil {
			return nil, 0, err
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, total, nil
}

// SearchProducts 搜索商品
func (s *ProductService) SearchProducts(ctx context.Context, params *SearchRequest) ([]model.ProductCardDTO, int64, error) {
	searchParams := &SearchParams{
		Keyword:  params.Keyword,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	if params.PriceMin != nil {
		searchParams.MinPrice = params.PriceMin
	}
	if params.PriceMax != nil {
		searchParams.MaxPrice = params.PriceMax
	}
	if params.ConditionID != nil {
		searchParams.ConditionID = params.ConditionID
	}

	return s.Search(ctx, searchParams)
}

// Search 搜索商品（与controller中使用的方法名保持一致）
func (s *ProductService) Search(ctx context.Context, params *SearchParams) ([]model.ProductCardDTO, int64, error) {
	if s.productRepo == nil {
		return nil, 0, fmt.Errorf("服务未初始化")
	}

	repoParams := repository.SearchParams{
		Keyword:      params.Keyword,
		Page:         params.Page,
		PageSize:     params.PageSize,
		PriceMin:     valueOrZero(params.MinPrice),
		PriceMax:     valueOrZero(params.MaxPrice),
		ConditionID:  valueOrZeroInt64(params.ConditionID),
		ConditionIDs: params.ConditionIDs,
		Sort:         params.Sort,
	}

	products, total, err := s.productRepo.Search(ctx, repoParams)
	if err != nil {
		return nil, 0, err
	}

	dtoList := make([]model.ProductCardDTO, 0, len(products))
	for i := range products {
		dto, err := s.toCardDTO(ctx, &products[i])
		if err != nil {
			return nil, 0, err
		}
		dtoList = append(dtoList, dto)
	}
	return dtoList, total, nil
}

// GetProductsByCategory 获取分类商品
func (s *ProductService) GetProductsByCategory(ctx context.Context, categoryID int64, params *SearchRequest) ([]model.ProductCardDTO, int64, error) {
	searchParams := &SearchParams{
		Keyword:  params.Keyword,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	if params.PriceMin != nil {
		searchParams.MinPrice = params.PriceMin
	}
	if params.PriceMax != nil {
		searchParams.MaxPrice = params.PriceMax
	}
	if params.ConditionID != nil {
		searchParams.ConditionID = params.ConditionID
	}
	searchParams.CategoryID = &categoryID

	return s.ListByCategory(ctx, categoryID, searchParams)
}

// ListByCategory 按分类列出商品（与controller中使用的方法名保持一致）
func (s *ProductService) ListByCategory(ctx context.Context, categoryID int64, params *SearchParams) ([]model.ProductCardDTO, int64, error) {
	if s.productRepo == nil {
		return nil, 0, fmt.Errorf("服务未初始化")
	}

	repoParams := repository.SearchParams{
		Keyword:      params.Keyword,
		Page:         params.Page,
		PageSize:     params.PageSize,
		PriceMin:     valueOrZero(params.MinPrice),
		PriceMax:     valueOrZero(params.MaxPrice),
		ConditionID:  valueOrZeroInt64(params.ConditionID),
		ConditionIDs: params.ConditionIDs,
		Sort:         params.Sort,
	}

	products, total, err := s.productRepo.ListByCategory(ctx, categoryID, repoParams)
	if err != nil {
		return nil, 0, err
	}

	dtoList := make([]model.ProductCardDTO, 0, len(products))
	for i := range products {
		dto, err := s.toCardDTO(ctx, &products[i])
		if err != nil {
			return nil, 0, err
		}
		dtoList = append(dtoList, dto)
	}
	return dtoList, total, nil
}

// AddProductImage 添加商品图片
func (s *ProductService) AddProductImage(ctx context.Context, userID, productID int64, file multipart.File, header *multipart.FileHeader) (*model.ProductImage, error) {
	product, images, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("商品不存在")
		}
		return nil, err
	}

	if product.SellerID != userID {
		return nil, fmt.Errorf("无权限操作该商品")
	}
	if product.Status == "Sold" {
		return nil, fmt.Errorf("已售出的商品不能修改")
	}

	url, err := util.SaveImage(file, header)
	if err != nil {
		return nil, err
	}

	sortOrder := len(images) + 1
	isPrimary := len(images) == 0
	image := &model.ProductImage{
		ProductID: productID,
		URL:       url,
		SortOrder: sortOrder,
		IsPrimary: isPrimary,
	}

	if s.db == nil {
		return nil, fmt.Errorf("服务未初始化")
	}

	if err := s.db.WithContext(ctx).Create(image).Error; err != nil {
		return nil, err
	}

	if isPrimary {
		_ = s.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", productID).Update("main_image_url", url).Error
	}

	return image, nil
}

// SetPrimaryImage 设置主图
func (s *ProductService) SetPrimaryImage(ctx context.Context, userID, productID, imageID int64) error {
	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("商品不存在")
		}
		return err
	}
	if product.SellerID != userID {
		return fmt.Errorf("无权限操作该商品")
	}
	if product.Status == "Sold" {
		return fmt.Errorf("已售出的商品不能修改")
	}
	if s.db == nil {
		return fmt.Errorf("服务未初始化")
	}

	var image model.ProductImage
	if err := s.db.WithContext(ctx).Where("id = ? AND product_id = ?", imageID, productID).First(&image).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("图片不存在")
		}
		return err
	}

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Model(&model.ProductImage{}).Where("product_id = ?", productID).Update("is_primary", false).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&model.ProductImage{}).Where("id = ?", imageID).Update("is_primary", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&model.Product{}).Where("id = ?", productID).Update("main_image_url", image.URL).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// UpdateImageSortOrder 更新图片排序
func (s *ProductService) UpdateImageSortOrder(ctx context.Context, userID, productID, imageID int64, sortOrder int) error {
	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("商品不存在")
		}
		return err
	}
	if product.SellerID != userID {
		return fmt.Errorf("无权限操作该商品")
	}
	if product.Status == "Sold" {
		return fmt.Errorf("已售出的商品不能修改")
	}
	if s.db == nil {
		return fmt.Errorf("服务未初始化")
	}

	result := s.db.WithContext(ctx).Model(&model.ProductImage{}).
		Where("id = ? AND product_id = ?", imageID, productID).
		Update("sort_order", sortOrder)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("图片不存在")
	}
	return nil
}

// DeleteProductImage 删除商品图片
func (s *ProductService) DeleteProductImage(ctx context.Context, userID, productID, imageID int64) error {
	product, images, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("商品不存在")
		}
		return err
	}

	if product.SellerID != userID {
		return fmt.Errorf("无权限操作该商品")
	}
	if product.Status == "Sold" {
		return fmt.Errorf("已售出的商品不能修改")
	}
	if len(images) <= 1 {
		return fmt.Errorf("至少保留一张图片")
	}
	if s.db == nil {
		return fmt.Errorf("服务未初始化")
	}

	tx := s.db.WithContext(ctx).Begin()
	var target model.ProductImage
	if err := tx.Where("id = ? AND product_id = ?", imageID, productID).First(&target).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("图片不存在")
		}
		return err
	}

	if err := tx.Delete(&target).Error; err != nil {
		tx.Rollback()
		return err
	}

	if target.IsPrimary {
		var newPrimary model.ProductImage
		if err := tx.Where("product_id = ?", productID).
			Order("sort_order ASC, id ASC").
			First(&newPrimary).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&model.ProductImage{}).Where("id = ?", newPrimary.ID).Update("is_primary", true).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&model.Product{}).Where("id = ?", productID).Update("main_image_url", newPrimary.URL).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func pickSellerWechat(wechat string, viewerIsSeller bool) *string {
	if strings.TrimSpace(wechat) == "" {
		return nil
	}
	if !viewerIsSeller {
		return nil
	}
	w := strings.TrimSpace(wechat)
	return &w
}

func buildStatusCacheKey(productID int64) string {
	return fmt.Sprintf("product:status:%d:last", productID)
}

func buildDetailCacheKey(productID int64) string {
	return fmt.Sprintf("product:detail:%d", productID)
}

func valueOrZero(val *float64) float64 {
	if val == nil {
		return 0
	}
	return *val
}

func valueOrZeroInt64(val *int64) int64 {
	if val == nil {
		return 0
	}
	return *val
}

func parseTagIDs(tagIDs string) []int64 {
	if strings.TrimSpace(tagIDs) == "" {
		return nil
	}
	parts := strings.Split(tagIDs, ",")
	result := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if id, err := strconv.ParseInt(p, 10, 64); err == nil {
			result = append(result, id)
		}
	}
	return result
}

func (s *ProductService) toCardDTO(ctx context.Context, p *model.Product) (model.ProductCardDTO, error) {
	main := p.MainImageURL
	if main == "" && s.db != nil {
		_ = s.db.WithContext(ctx).
			Table("product_images").
			Select("url").
			Where("product_id = ?", p.ID).
			Order("sort_order ASC, id ASC").
			Limit(1).
			Scan(&main)
	}

	return model.ProductCardDTO{
		ID:          p.ID,
		Title:       p.Title,
		Price:       p.Price,
		MainImage:   main,
		Status:      p.Status,
		SellerID:    p.SellerID,
		CategoryID:  p.CategoryID,
		ConditionID: p.ConditionID,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func (s *ProductService) buildDetailDTO(ctx context.Context, product *model.Product, images []model.ProductImage, tagIDs []int64, viewerID *int64) (*model.ProductDetailDTO, error) {
	// 查询新旧程度名称（失败仅记录日志，继续返回数据）
	var conditionName string
	if s.db != nil {
		if err := s.db.WithContext(ctx).
			Table("product_conditions").
			Select("name").
			Where("id = ?", product.ConditionID).
			Limit(1).
			Scan(&conditionName).Error; err != nil {
			log.Printf("warn: get condition name failed for product %d: %v", product.ID, err)
		}
	}

	// 查询卖家信息（失败返回占位）
	seller, err := s.userRepo.GetByID(ctx, product.SellerID)
	if err != nil {
		log.Printf("warn: get seller info failed for product %d seller %d: %v", product.ID, product.SellerID, err)
		seller = &model.User{ID: product.SellerID}
	}

	viewerIsSeller := viewerID != nil && *viewerID == product.SellerID

	// 计算主图
	mainImage := product.MainImageURL
	if mainImage == "" && len(images) > 0 {
		for _, img := range images {
			if img.IsPrimary {
				mainImage = img.URL
				break
			}
		}
		if mainImage == "" {
			mainImage = images[0].URL
		}
	}

	return &model.ProductDetailDTO{
		ID:             product.ID,
		Title:          product.Title,
		Description:    product.Description,
		Price:          product.Price,
		CategoryID:     product.CategoryID,
		ConditionID:    product.ConditionID,
		ConditionName:  conditionName,
		MainImageURL:   mainImage,
		Images:         images,
		TagIDs:         tagIDs,
		Seller:         model.SellerInfo{ID: seller.ID, Nickname: seller.Nickname, AvatarUrl: seller.AvatarUrl},
		ViewerIsSeller: viewerIsSeller,
		Status:         product.Status,
		CreatedAt:      product.CreatedAt,
		UpdatedAt:      product.UpdatedAt,
		SellerWechat:   pickSellerWechat(seller.WechatID, viewerIsSeller),
	}, nil
}
