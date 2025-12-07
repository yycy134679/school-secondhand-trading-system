package recommend

import (
	"context"
	"sort"
	"time"

	"gorm.io/gorm"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
)

// RedisClient Redis客户端接口（用于解耦）
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

// RecommendService 推荐服务
type RecommendService struct {
	viewRecordRepo repository.ViewRecordRepository
	productRepo    repository.ProductRepository
	db             *gorm.DB
	redis          RedisClient // 使用接口，可选
}

// NewRecommendService 创建推荐服务实例
// redis 参数可以为 nil，此时不使用缓存功能
func NewRecommendService(
	viewRecordRepo repository.ViewRecordRepository,
	productRepo repository.ProductRepository,
	db *gorm.DB,
	redis RedisClient,
) *RecommendService {
	return &RecommendService{
		viewRecordRepo: viewRecordRepo,
		productRepo:    productRepo,
		db:             db,
		redis:          redis,
	}
}

// RecordView 记录浏览
func (s *RecommendService) RecordView(ctx context.Context, userID, productID int64) error {
	return s.viewRecordRepo.AddView(ctx, userID, productID)
}

// GetRecommendations 获取推荐商品
// maxCount: 返回的最大推荐数量
func (s *RecommendService) GetRecommendations(ctx context.Context, userID int64, maxCount int) ([]model.Product, error) {
	// 注意：Redis 缓存功能已移除，始终计算推荐
	// 如需启用缓存，请实现 RedisClient 接口并在初始化时传入

	// 计算推荐
	recommendations, err := s.calculateRecommendations(ctx, userID, maxCount)
	if err != nil {
		return nil, err
	}

	return recommendations, nil
}

// calculateRecommendations 计算推荐商品
func (s *RecommendService) calculateRecommendations(ctx context.Context, userID int64, maxCount int) ([]model.Product, error) {
	if maxCount <= 0 {
		return []model.Product{}, nil
	}

	// 1. 获取用户最近浏览的商品(最多20条)，按浏览时间倒序
	recentViews, err := s.viewRecordRepo.ListRecentViews(ctx, userID, 20)
	if err != nil {
		return nil, err
	}

	if len(recentViews) == 0 {
		return []model.Product{}, nil
	}

	// 2. 提取浏览过的商品ID
	viewedProductIDs := make([]int64, len(recentViews))
	for i, view := range recentViews {
		viewedProductIDs[i] = view.ProductID
	}

	// 3. 过滤已删除/下架的浏览记录，基于仍在售的商品抽取分类
	var viewedProducts []model.Product
	if err := s.db.WithContext(ctx).
		Where("id IN ? AND status = ?", viewedProductIDs, "ForSale").
		Find(&viewedProducts).Error; err != nil {
		return nil, err
	}

	if len(viewedProducts) == 0 {
		return []model.Product{}, nil
	}

	productMap := make(map[int64]*model.Product, len(viewedProducts))
	for i := range viewedProducts {
		productMap[viewedProducts[i].ID] = &viewedProducts[i]
	}

	// 4. 按浏览时间倒序选出最新的两个不同分类
	selectedCategories := make([]int64, 0, 2)
	categorySeen := make(map[int64]struct{})
	for _, view := range recentViews {
		product, ok := productMap[view.ProductID]
		if !ok {
			continue
		}
		if _, exists := categorySeen[product.CategoryID]; exists {
			continue
		}
		categorySeen[product.CategoryID] = struct{}{}
		selectedCategories = append(selectedCategories, product.CategoryID)
		if len(selectedCategories) == 2 {
			break
		}
	}

	if len(selectedCategories) == 0 {
		return []model.Product{}, nil
	}

	// 5. 按分类查询在售商品，排除已浏览和本人发布，最多4条（两类各2，单类最多4）
	categoryProducts := make(map[int64][]model.Product, len(selectedCategories))
	for _, categoryID := range selectedCategories {
		products, err := s.fetchCategoryProducts(ctx, categoryID, userID, viewedProductIDs, maxCount)
		if err != nil {
			return nil, err
		}
		categoryProducts[categoryID] = products
	}

	var results []model.Product
	if len(selectedCategories) == 1 {
		products := categoryProducts[selectedCategories[0]]
		if len(products) > maxCount {
			products = products[:maxCount]
		}
		results = append(results, products...)
	} else {
		// 先从每个分类各取2条（不足则取全部），再用剩余的按分类顺序补齐到上限
		categoryTaken := make(map[int64]int, len(selectedCategories))
		for _, categoryID := range selectedCategories {
			products := categoryProducts[categoryID]
			take := 2
			if len(products) < take {
				take = len(products)
			}
			results = append(results, products[:take]...)
			categoryTaken[categoryID] = take
		}

		remaining := maxCount - len(results)
		if remaining > 0 {
			for _, categoryID := range selectedCategories {
				if remaining == 0 {
					break
				}
				products := categoryProducts[categoryID]
				start := categoryTaken[categoryID]
				if start >= len(products) {
					continue
				}
				end := start + remaining
				if end > len(products) {
					end = len(products)
				}
				if start < end {
					results = append(results, products[start:end]...)
					categoryTaken[categoryID] = end
					remaining = maxCount - len(results)
				}
			}
		}
	}

	// 6. 全局按创建时间倒序排序并裁剪到上限
	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})

	if len(results) > maxCount {
		results = results[:maxCount]
	}

	return results, nil
}

// getProductsByIDs 根据ID列表获取商品
func (s *RecommendService) getProductsByIDs(ctx context.Context, productIDs []int64) ([]model.Product, error) {
	if len(productIDs) == 0 {
		return []model.Product{}, nil
	}

	var products []model.Product
	err := s.db.WithContext(ctx).
		Where("id IN ? AND status = ?", productIDs, "ForSale").
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

// HomeData 首页数据结构
type HomeData struct {
	Recommendations []model.ProductCardDTO `json:"recommendations"`
	Latest          []model.ProductCardDTO `json:"latest"`
	TotalCount      int64                  `json:"totalCount"`
}

// GetHomeData 获取首页数据
func (s *RecommendService) GetHomeData(ctx context.Context, userID *int64, page, pageSize int) (*HomeData, error) {
	const maxRecommendations = 4

	var recommendations []model.Product
	recommendIDSet := make(map[int64]struct{})
	viewedIDSet := make(map[int64]struct{})

	// 如果用户已登录,获取推荐商品
	if userID != nil {
		var err error
		recommendations, err = s.GetRecommendations(ctx, *userID, maxRecommendations)
		if err != nil {
			return nil, err
		}
		for _, p := range recommendations {
			recommendIDSet[p.ID] = struct{}{}
		}

		// 获取浏览过的商品，用于补齐时继续排除
		recentViews, err := s.viewRecordRepo.ListRecentViews(ctx, *userID, 20)
		if err != nil {
			return nil, err
		}
		for _, view := range recentViews {
			viewedIDSet[view.ProductID] = struct{}{}
		}
	}

	// 获取最新在售商品（先排除已推荐的，以便最新列表去重）
	excludeIDs := make([]int64, 0, len(recommendIDSet))
	for id := range recommendIDSet {
		excludeIDs = append(excludeIDs, id)
	}
	latestProducts, total, err := s.productRepo.ListLatestForSale(ctx, excludeIDs, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 登录用户：如果推荐数不足上限,用最新商品补充（仍排除本人发布和已浏览）
	usedFromLatest := make(map[int64]struct{})
	if userID != nil && len(recommendations) < maxRecommendations && len(latestProducts) > 0 {
		for i := range latestProducts {
			if len(recommendations) >= maxRecommendations {
				break
			}
			if latestProducts[i].SellerID == *userID {
				continue
			}
			if _, exists := recommendIDSet[latestProducts[i].ID]; exists {
				continue
			}
			if _, viewed := viewedIDSet[latestProducts[i].ID]; viewed {
				continue
			}
			recommendations = append(recommendations, latestProducts[i])
			recommendIDSet[latestProducts[i].ID] = struct{}{}
			usedFromLatest[latestProducts[i].ID] = struct{}{}
		}
	}

	// 推荐按创建时间倒序
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].CreatedAt.After(recommendations[j].CreatedAt)
	})

	recommendDTOs := make([]model.ProductCardDTO, len(recommendations))
	for i, p := range recommendations {
		recommendDTOs[i] = s.toProductCardDTO(&p)
	}

	// 过滤掉已用于补齐的最新商品，保持最新列表与推荐去重
	filteredLatest := make([]model.Product, 0, len(latestProducts))
	for i := range latestProducts {
		if _, used := usedFromLatest[latestProducts[i].ID]; used {
			continue
		}
		filteredLatest = append(filteredLatest, latestProducts[i])
	}
	latestDTOs := make([]model.ProductCardDTO, len(filteredLatest))
	for i, p := range filteredLatest {
		latestDTOs[i] = s.toProductCardDTO(&p)
	}

	// 调整最新列表总数（移除补齐占用的商品）
	adjustedTotal := total - int64(len(usedFromLatest))
	if adjustedTotal < 0 {
		adjustedTotal = 0
	}

	return &HomeData{
		Recommendations: recommendDTOs,
		Latest:          latestDTOs,
		TotalCount:      adjustedTotal,
	}, nil
}

// fetchCategoryProducts 获取某分类的在售商品，排除已浏览和本人发布，按创建时间倒序
func (s *RecommendService) fetchCategoryProducts(ctx context.Context, categoryID int64, userID int64, excludeIDs []int64, limit int) ([]model.Product, error) {
	if limit <= 0 {
		return []model.Product{}, nil
	}

	query := s.db.WithContext(ctx).
		Where("status = ? AND category_id = ?", "ForSale", categoryID).
		Where("seller_id <> ?", userID)

	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	var products []model.Product
	if err := query.Order("created_at DESC").Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// toProductCardDTO 转换为ProductCardDTO
func (s *RecommendService) toProductCardDTO(product *model.Product) model.ProductCardDTO {
	// 获取主图URL
	var mainImage string
	s.db.Table("product_images").
		Select("url").
		Where("product_id = ? AND is_primary = ?", product.ID, true).
		Limit(1).
		Scan(&mainImage)

	return model.ProductCardDTO{
		ID:          product.ID,
		Title:       product.Title,
		Price:       product.Price,
		MainImage:   mainImage,
		Status:      product.Status,
		SellerID:    product.SellerID,
		CategoryID:  product.CategoryID,
		ConditionID: product.ConditionID,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// RecentViewWithProduct 浏览记录带商品信息
type RecentViewWithProduct struct {
	ViewedAt time.Time            `json:"viewedAt"`
	Product  model.ProductCardDTO `json:"product"`
}

// GetRecentViewsWithProducts 获取用户最近浏览记录并关联商品信息
func (s *RecommendService) GetRecentViewsWithProducts(ctx context.Context, userID int64, limit int) ([]RecentViewWithProduct, error) {
	// 获取浏览记录
	views, err := s.viewRecordRepo.ListRecentViews(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	if len(views) == 0 {
		return []RecentViewWithProduct{}, nil
	}

	// 提取商品ID
	productIDs := make([]int64, len(views))
	for i, view := range views {
		productIDs[i] = view.ProductID
	}

	// 获取商品信息
	var products []model.Product
	err = s.db.WithContext(ctx).
		Where("id IN ?", productIDs).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	// 构建商品ID到商品的映射
	productMap := make(map[int64]*model.Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	// 组装结果
	result := make([]RecentViewWithProduct, 0, len(views))
	for _, view := range views {
		if product, exists := productMap[view.ProductID]; exists {
			result = append(result, RecentViewWithProduct{
				ViewedAt: view.ViewedAt,
				Product:  s.toProductCardDTO(product),
			})
		}
	}

	return result, nil
}
