package product

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/util"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
	"gorm.io/gorm"
)

// ProductService handles product business logic
type ProductService struct {
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	db          *gorm.DB
	redisClient interface { // Redis client interface
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Get(ctx context.Context, key string) (string, error)
		Del(ctx context.Context, key string) error
	}
}

// isAdmin checks if the user is an admin
func (s *ProductService) isAdmin(userID int64) bool {
	// Get user by ID
	user, err := s.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return false
	}

	// Check if user is admin using IsAdmin field
	return user.IsAdmin
}

// NewProductService creates a new product service instance
func NewProductService(
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	db *gorm.DB,
	redisClient interface {
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Get(ctx context.Context, key string) (string, error)
		Del(ctx context.Context, key string) error
	},
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		userRepo:    userRepo,
		db:          db,
		redisClient: redisClient,
	}
}

// CreateProductRequest represents the request for creating a product
type CreateProductRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       int64    `json:"price" binding:"required,gt=0"`
	CategoryID  int64    `json:"categoryId" binding:"required"`
	ConditionID int64    `json:"conditionId" binding:"required"`
	TagIDs      []int64  `json:"tagIds"`
	Images      []string `json:"images" binding:"required,min=1"`
}

// UpdateProductRequest represents the request for updating a product
type UpdateProductRequest struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Price       *int64    `json:"price" binding:"omitempty,gt=0"`
	CategoryID  *int64    `json:"categoryId"`
	ConditionID *int64    `json:"conditionId"`
	TagIDs      *[]int64  `json:"tagIds"`
	Images      *[]string `json:"images"`
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, sellerID int64, req CreateProductRequest) (*model.ProductDetailDTO, error) {
	// 校验必填字段与价格 > 0
	if req.Title == "" || req.Description == "" || req.Price <= 0 {
		return nil, errors.New("title, description and price must be valid")
	}

	// 校验卖家信息
	seller, err := s.userRepo.GetByID(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	// 校验卖家wechat_id非空
	if seller.WechatID == "" {
		return nil, fmt.Errorf("1001:请先完善微信号")
	}

	// 校验categoryId存在
	// 实际应查询分类表确认存在，这里简化实现
	// 假设分类ID是有效的，实际项目中应该调用分类仓库验证
	if req.CategoryID <= 0 {
		return nil, errors.New("invalid category id")
	}

	// 校验conditionId来自product_conditions
	// 实际应查询条件表确认存在，这里简化实现
	// 假设条件ID是有效的，实际项目中应该调用条件仓库验证
	if req.ConditionID <= 0 {
		return nil, errors.New("invalid condition id")
	}

	// 校验tagIds合法
	// 实际应查询标签表确认存在，这里简化实现
	// 假设标签ID是有效的，实际项目中应该调用标签仓库验证
	for _, tagID := range req.TagIDs {
		if tagID <= 0 {
			return nil, errors.New("invalid tag id")
		}
	}

	// 处理图片上传 - 调用统一文件服务获取URL
	// 这里简化实现，实际应该调用文件服务上传并获取URL
	// 假设req.Images中已经是有效的URL
	if len(req.Images) == 0 {
		return nil, errors.New("at least one image is required")
	}

	// 构建Product模型
	product := &model.Product{
		SellerID:     sellerID,
		Title:        req.Title,
		Description:  req.Description,
		Price:        req.Price,
		CategoryID:   req.CategoryID,
		ConditionID:  req.ConditionID,
		Status:       "ForSale",     // 设置初始状态为在售
		MainImageURL: req.Images[0], // 第一张图片作为主图
	}

	// 构建ProductImage列表
	var images []model.ProductImage
	for i, imageURL := range req.Images {
		isMain := i == 0 // 第一张作为主图
		images = append(images, model.ProductImage{
			ProductID: 0, // 会在插入后由数据库生成
			URL:       imageURL,
			SortOrder: i,
			IsPrimary: isMain,
		})
	}

	// 写入数据库
	productID, err := s.productRepo.Create(ctx, product, images, req.TagIDs)
	if err != nil {
		return nil, err
	}

	// 重新获取完整信息构建DTO
	createdProduct, productImages, tagIDs, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 构建返回DTO
	dto := &model.ProductDetailDTO{
		ID:           createdProduct.ID,
		Title:        createdProduct.Title,
		Description:  createdProduct.Description,
		Price:        createdProduct.Price,
		MainImageUrl: createdProduct.MainImageURL,
		Status:       createdProduct.Status,
		CategoryID:   createdProduct.CategoryID,
		ConditionID:  createdProduct.ConditionID,
		TagIds:       tagIDs,
		Seller: model.SellerInfoDTO{
			ID:        seller.ID,
			Nickname:  seller.Nickname,
			AvatarURL: seller.AvatarUrl,
		},
		CreatedAt: createdProduct.CreatedAt,
		UpdatedAt: createdProduct.UpdatedAt,
		Images:    make([]model.ProductImageDTO, 0, len(productImages)),
	}

	// 填充图片信息
	for _, img := range productImages {
		dto.Images = append(dto.Images, model.ProductImageDTO{
			ID:        img.ID,
			URL:       img.URL,
			SortOrder: img.SortOrder,
			IsPrimary: img.IsPrimary,
		})
	}

	return dto, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, actorID, productID int64, req UpdateProductRequest, isAdmin bool) (*model.ProductDetailDTO, error) {
	// 获取产品信息
	product, currentImages, currentTagIDs, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 检查权限
	if !isAdmin && product.SellerID != actorID {
		return nil, errors.New("permission denied")
	}

	// 检查状态 - 普通卖家禁止编辑已售商品
	// 管理员仅在status不变时可编辑Sold的其他字段
	if product.Status == "Sold" {
		if !isAdmin {
			return nil, errors.New("cannot edit sold product")
		}
		// 管理员可以编辑已售商品的其他字段，但不能更改状态
		// 这里我们假设在repo层不会修改状态，所以不需要特殊处理
	}

	// 验证更新字段
	if req.Title != nil && *req.Title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if req.Description != nil && *req.Description == "" {
		return nil, errors.New("description cannot be empty")
	}
	if req.Price != nil && *req.Price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	if req.CategoryID != nil && *req.CategoryID <= 0 {
		return nil, errors.New("invalid category id")
	}
	if req.ConditionID != nil && *req.ConditionID <= 0 {
		return nil, errors.New("invalid condition id")
	}
	if req.TagIDs != nil {
		for _, tagID := range *req.TagIDs {
			if tagID <= 0 {
				return nil, errors.New("invalid tag id")
			}
		}
	}

	// 更新允许的字段
	if req.Title != nil {
		product.Title = *req.Title
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.ConditionID != nil {
		product.ConditionID = *req.ConditionID
	}

	// 处理图片更新 - 支持图片替换或增删
	var updatedImages []model.ProductImage
	if req.Images != nil {
		// 验证图片不为空
		if len(*req.Images) == 0 {
			return nil, errors.New("at least one image is required")
		}

		// 构建新的图片列表
		for i, imageURL := range *req.Images {
			isMain := i == 0
			updatedImages = append(updatedImages, model.ProductImage{
				ProductID: productID, // 关联到当前产品
				URL:       imageURL,
				SortOrder: i,
				IsPrimary: isMain,
			})

			// 更新主图URL
			if isMain {
				product.MainImageURL = imageURL
			}
		}
	} else {
		// 如果没有提供新图片，保留原有图片
		updatedImages = currentImages
	}

	// 处理标签更新
	tagIDs := currentTagIDs
	if req.TagIDs != nil {
		tagIDs = *req.TagIDs
	}

	// 执行更新
	err = s.productRepo.Update(ctx, product, updatedImages, tagIDs, isAdmin)
	if err != nil {
		return nil, err
	}

	// 重新获取更新后的数据构建DTO
	updatedProduct, productImages, newTagIDs, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 获取卖家信息
	seller, err := s.userRepo.GetByID(ctx, updatedProduct.SellerID)
	if err != nil {
		return nil, err
	}

	// 构建返回DTO
	dto := &model.ProductDetailDTO{
		ID:           updatedProduct.ID,
		Title:        updatedProduct.Title,
		Description:  updatedProduct.Description,
		Price:        updatedProduct.Price,
		MainImageUrl: updatedProduct.MainImageURL,
		Status:       updatedProduct.Status,
		CategoryID:   updatedProduct.CategoryID,
		ConditionID:  updatedProduct.ConditionID,
		TagIds:       newTagIDs,
		Seller: model.SellerInfoDTO{
			ID:        seller.ID,
			Nickname:  seller.Nickname,
			AvatarURL: seller.AvatarUrl,
		},
		CreatedAt: updatedProduct.CreatedAt,
		UpdatedAt: updatedProduct.UpdatedAt,
		Images:    make([]model.ProductImageDTO, 0, len(productImages)),
	}

	// 填充图片信息
	for _, img := range productImages {
		dto.Images = append(dto.Images, model.ProductImageDTO{
			ID:        img.ID,
			URL:       img.URL,
			SortOrder: img.SortOrder,
			IsPrimary: img.IsPrimary,
		})
	}

	return dto, nil
}

// ChangeStatus changes the product status
func (s *ProductService) ChangeStatus(ctx context.Context, sellerID, productID int64, action string) error {
	// 获取产品信息
	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// 检查权限
	if product.SellerID != sellerID {
		return errors.New("permission denied")
	}

	// 检查是否为终态
	if product.Status == "Sold" {
		return errors.New("3004:product is already sold")
	}

	// 验证动作参数
	if action != "delist" && action != "relist" && action != "sold" {
		return errors.New("invalid action")
	}

	// 根据当前状态和动作映射合法流转
	var fromStatus, toStatus string
	fromStatus = product.Status
	validTransition := false

	switch {
	case action == "delist" && product.Status == "ForSale":
		toStatus = "Delisted"
		validTransition = true
	case action == "relist" && product.Status == "Delisted":
		toStatus = "ForSale"
		validTransition = true
	case action == "sold" && product.Status == "ForSale":
		toStatus = "Sold"
		validTransition = true
	}

	if !validTransition {
		return errors.New("invalid status transition")
	}

	// 执行状态更新
	err = s.productRepo.UpdateStatus(ctx, productID, fromStatus, toStatus)
	if err != nil {
		return err
	}

	// 对成功的delist/relist写入Redis撤销记录
	// 键包含productID + sellerID，TTL为3秒
	if action == "delist" || action == "relist" {
		redisKey := fmt.Sprintf("product:status:undo:%d:%d", productID, sellerID)
		redisValue := fmt.Sprintf("%s:%s", fromStatus, toStatus)
		// 设置TTL为3秒
		s.redisClient.Set(ctx, redisKey, redisValue, 3*time.Second)
	}

	return nil
}

// UndoLastStatusChange undoes the last status change
func (s *ProductService) UndoLastStatusChange(ctx context.Context, sellerID, productID int64) error {
	// 从Redis读取缓存的fromStatus/toStatus
	redisKey := fmt.Sprintf("product:status:undo:%d:%d", productID, sellerID)
	cachedValue, err := s.redisClient.Get(ctx, redisKey)

	// 若不存在或已超时，返回错误码3005
	if err != nil {
		return errors.New("3005:undo not available")
	}

	// 解析缓存值
	var fromStatus, toStatus string
	if _, err = fmt.Sscanf(cachedValue, "%s:%s", &fromStatus, &toStatus); err != nil {
		// 缓存格式错误，删除无效缓存
		s.redisClient.Del(ctx, redisKey)
		return errors.New("3005:undo not available")
	}

	// 获取当前产品状态
	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// 检查权限
	if product.SellerID != sellerID {
		return errors.New("permission denied")
	}

	// 检查是否为终态
	if product.Status == "Sold" {
		return errors.New("3004:cannot undo sold status")
	}

	// 检查当前状态是否与缓存toStatus一致
	if product.Status != toStatus {
		// 状态已变化，删除缓存并返回错误码3005
		s.redisClient.Del(ctx, redisKey)
		return errors.New("3005:status has changed")
	}

	// 执行状态回滚
	err = s.productRepo.UpdateStatus(ctx, productID, toStatus, fromStatus)
	if err != nil {
		return err
	}

	// 删除缓存
	s.redisClient.Del(ctx, redisKey)

	return nil
}

// GetProductDetail gets product detail
func (s *ProductService) GetProductDetail(ctx context.Context, productID int64, viewerID *int64) (*model.ProductDetailDTO, error) {
	// 查询详情与关联图片/tags/seller
	product, productImages, tagIDs, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 获取卖家信息
	seller, err := s.userRepo.GetByID(ctx, product.SellerID)
	if err != nil {
		return nil, err
	}

	// 构建返回DTO
	dto := &model.ProductDetailDTO{
		ID:           product.ID,
		Title:        product.Title,
		Description:  product.Description,
		Price:        product.Price,
		MainImageUrl: product.MainImageURL,
		Status:       product.Status,
		CategoryID:   product.CategoryID,
		ConditionID:  product.ConditionID,
		TagIds:       tagIDs,
		Seller: model.SellerInfoDTO{
			ID:        seller.ID,
			Nickname:  seller.Nickname,
			AvatarURL: seller.AvatarUrl,
		},
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
		Images:    make([]model.ProductImageDTO, 0, len(productImages)),
	}

	// 填充图片信息
	for _, img := range productImages {
		dto.Images = append(dto.Images, model.ProductImageDTO{
			ID:        img.ID,
			URL:       img.URL,
			SortOrder: img.SortOrder,
			IsPrimary: img.IsPrimary,
		})
	}

	// 若viewerID非空，调用推荐模块记录浏览
	// 根据API设计，使用独立的/contact接口获取sellerWechat
	// 这里使用goroutine异步记录浏览，不阻塞主流程
	if viewerID != nil {
		go func() {
			// 在实际项目中，这里应该调用推荐服务的API
			// 例如: s.recommendService.RecordView(*viewerID, productID)
			// 这里简化处理
			_ = *viewerID // 避免未使用警告
			_ = productID
		}()
	}

	return dto, nil
}

// ListMyProducts lists products by seller with pagination
func (s *ProductService) ListMyProducts(ctx context.Context, sellerID int64, keyword string, page, pageSize int) ([]model.ProductCardDTO, int64, error) {
	// 参数验证和默认值处理
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // 默认页大小为20，最大不超过100
	}

	// 调用仓库层方法获取数据
	products, total, err := s.productRepo.ListBySeller(ctx, sellerID, keyword, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为DTO
	dtos := make([]model.ProductCardDTO, 0, len(products))
	for _, product := range products {
		dto := model.ProductCardDTO{
			ID:           product.ID,
			Title:        product.Title,
			Price:        product.Price,
			MainImageUrl: product.MainImageURL,
			Status:       product.Status,
			CreatedAt:    product.CreatedAt,
		}
		dtos = append(dtos, dto)
	}

	// 确保返回正确的总数，这将用于前端分页计算
	return dtos, total, nil
}

// UploadProductImage 上传商品图片
func (s *ProductService) UploadProductImage(ctx context.Context, userID, productID int64, file multipart.File, header *multipart.FileHeader) (int64, string, error) {
	// 1. 检查商品是否存在以及权限
	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return 0, "", fmt.Errorf("商品不存在")
	}

	// 检查是否是卖家或管理员
	if product.SellerID != userID && !s.isAdmin(userID) {
		return 0, "", errors.New("无权限操作此商品")
	}

	// 检查商品状态
	if product.Status == model.ProductStatusSold {
		return 0, "", errors.New("已售出商品不允许添加图片")
	}

	// 2. 保存图片
	imageURL, err := util.SaveImage(file, header)
	if err != nil {
		return 0, "", fmt.Errorf("图片保存失败: %w", err)
	}

	// 3. 检查是否需要设置为主图
	hasPrimary, err := s.productRepo.HasPrimaryImage(ctx, productID)
	if err != nil {
		return 0, "", fmt.Errorf("检查主图失败: %w", err)
	}

	// 4. 创建图片记录
	image := &model.ProductImage{
		ProductID: productID,
		URL:       imageURL,
		SortOrder: 0, // 暂时设置为0，后续可以优化
		IsPrimary: !hasPrimary,
	}

	// 5. 保存图片记录
	imageID, err := s.productRepo.CreateImage(ctx, image)
	if err != nil {
		return 0, "", fmt.Errorf("保存图片记录失败: %w", err)
	}

	// 如果是主图，更新商品的主图URL
	if image.IsPrimary {
		if err := s.productRepo.UpdateMainImageURL(ctx, productID, imageURL); err != nil {
			return 0, "", fmt.Errorf("更新主图URL失败: %w", err)
		}
	}

	return imageID, imageURL, nil
}

// SetPrimaryImage 设置主图
func (s *ProductService) SetPrimaryImage(ctx context.Context, userID, productID, imageID int64) error {
	// 1. 检查商品是否存在以及权限
	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("商品不存在")
	}

	// 检查是否是卖家或管理员
	if product.SellerID != userID && !s.isAdmin(userID) {
		return errors.New("无权限操作此商品")
	}

	// 检查商品状态
	if product.Status == model.ProductStatusSold {
		return errors.New("已售出商品不允许修改图片")
	}

	// 2. 检查图片是否存在且属于该商品
	image, err := s.productRepo.GetImageByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("图片不存在")
	}

	if image.ProductID != productID {
		return errors.New("图片不属于该商品")
	}

	// 3. 事务设置主图
	if err := s.productRepo.SetPrimaryImage(ctx, productID, imageID); err != nil {
		return fmt.Errorf("设置主图失败: %w", err)
	}

	// 4. 更新商品的主图URL
	if err := s.productRepo.UpdateMainImageURL(ctx, productID, image.URL); err != nil {
		return fmt.Errorf("更新主图URL失败: %w", err)
	}

	return nil
}

// UpdateImageSort 更新图片排序
func (s *ProductService) UpdateImageSort(ctx context.Context, userID, productID, imageID int64, sortOrder int) error {
	// 1. 检查商品是否存在以及权限
	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("商品不存在")
	}

	// 检查是否是卖家或管理员
	if product.SellerID != userID && !s.isAdmin(userID) {
		return errors.New("无权限操作此商品")
	}

	// 检查商品状态
	if product.Status == model.ProductStatusSold {
		return errors.New("已售出商品不允许修改图片")
	}

	// 2. 检查图片是否存在且属于该商品
	image, err := s.productRepo.GetImageByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("图片不存在")
	}

	if image.ProductID != productID {
		return errors.New("图片不属于该商品")
	}

	// 3. 更新排序
	if err := s.productRepo.UpdateImageSort(ctx, imageID, sortOrder); err != nil {
		return fmt.Errorf("更新排序失败: %w", err)
	}

	return nil
}

// DeleteProductImage 删除商品图片
func (s *ProductService) DeleteProductImage(ctx context.Context, userID, productID, imageID int64) error {
	// 1. 检查商品是否存在以及权限
	product, _, _, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("商品不存在")
	}

	// 检查是否是卖家或管理员
	if product.SellerID != userID && !s.isAdmin(userID) {
		return errors.New("无权限操作此商品")
	}

	// 检查商品状态
	if product.Status == model.ProductStatusSold {
		return errors.New("已售出商品不允许修改图片")
	}

	// 2. 检查图片是否存在且属于该商品
	image, err := s.productRepo.GetImageByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("图片不存在")
	}

	if image.ProductID != productID {
		return errors.New("图片不属于该商品")
	}

	// 3. 检查是否是主图
	isPrimary := image.IsPrimary

	// 4. 删除图片
	if err := s.productRepo.DeleteImage(ctx, imageID); err != nil {
		return fmt.Errorf("删除图片失败: %w", err)
	}

	// 5. 如果删除的是主图，选择新的主图
	if isPrimary {
		// 获取下一个应该设为主图的图片
		nextPrimary, err := s.productRepo.GetFirstImage(ctx, productID)
		if err != nil {
			// 如果没有图片了，将主图URL设为空
			if err := s.productRepo.UpdateMainImageURL(ctx, productID, ""); err != nil {
				return fmt.Errorf("更新主图URL失败: %w", err)
			}
		} else {
			// 设置新的主图
			if err := s.productRepo.SetPrimaryImage(ctx, productID, nextPrimary.ID); err != nil {
				return fmt.Errorf("设置新主图失败: %w", err)
			}
			if err := s.productRepo.UpdateMainImageURL(ctx, productID, nextPrimary.URL); err != nil {
				return fmt.Errorf("更新主图URL失败: %w", err)
			}
		}
	}

	return nil
}

// Search 搜索商品
func (s *ProductService) Search(ctx context.Context, params model.SearchParams) ([]model.ProductCardDTO, int64, error) {
	// 参数验证和默认值处理
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20 // 默认页大小为20，最大不超过100
	}

	// 查询商品列表
	products, total, err := s.productRepo.Search(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// 转换为DTO
	dtos := make([]model.ProductCardDTO, 0, len(products))
	for _, product := range products {
		dto := model.ProductCardDTO{
			ID:           product.ID,
			Title:        product.Title,
			Price:        product.Price,
			MainImageUrl: product.MainImageURL,
			Status:       product.Status,
			CreatedAt:    product.CreatedAt,
		}
		dtos = append(dtos, dto)
	}

	return dtos, total, nil
}

// ListByCategory 根据分类ID查询商品
func (s *ProductService) ListByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]model.ProductCardDTO, int64, error) {
	// 参数验证和默认值处理
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // 默认页大小为20，最大不超过100
	}

	// 创建搜索参数
	params := model.SearchParams{
		CategoryID: categoryID,
		Page:       page,
		PageSize:   pageSize,
	}

	// 查询商品列表
	products, total, err := s.productRepo.ListByCategory(ctx, categoryID, params)
	if err != nil {
		return nil, 0, err
	}

	// 转换为DTO
	dtos := make([]model.ProductCardDTO, 0, len(products))
	for _, product := range products {
		dto := model.ProductCardDTO{
			ID:           product.ID,
			Title:        product.Title,
			Price:        product.Price,
			MainImageUrl: product.MainImageURL,
			Status:       product.Status,
			CreatedAt:    product.CreatedAt,
		}
		dtos = append(dtos, dto)
	}

	return dtos, total, nil
}

// ListLatestForSale 获取最新在售商品列表
func (s *ProductService) ListLatestForSale(ctx context.Context, excludeIDs []int64, page, pageSize int) ([]model.ProductCardDTO, int64, error) {
	// 参数验证和默认值处理
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // 默认页大小为20，最大不超过100
	}

	// 查询最新在售商品
	products, total, err := s.productRepo.ListLatestForSale(ctx, excludeIDs, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为DTO
	dtos := make([]model.ProductCardDTO, 0, len(products))
	for _, product := range products {
		dto := model.ProductCardDTO{
			ID:           product.ID,
			Title:        product.Title,
			Price:        product.Price,
			MainImageUrl: product.MainImageURL,
			Status:       product.Status,
			CreatedAt:    product.CreatedAt,
		}
		dtos = append(dtos, dto)
	}

	return dtos, total, nil
}
