// Package recommend 提供基于用户浏览记录的商品推荐服务
package recommend

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
)

// RecommendationCache 推荐缓存结构
type RecommendationCache struct {
	ProductIDs []int64   // 推荐商品ID列表
	ExpiresAt  time.Time // 缓存过期时间
}

// UndoRecord 撤销记录结构
type UndoRecord struct {
	FromStatus string    // 原状态
	ToStatus   string    // 目标状态
	ExpiresAt  time.Time // 过期时间
}

// RecommendService 推荐服务
// 使用内存缓存替代Redis，适用于单实例部署
type RecommendService struct {
	viewRecordRepo repository.ViewRecordRepository
	productRepo    repository.ProductRepository

	// 内存缓存 - 推荐结果
	recommendCache map[int64]*RecommendationCache // key: userID
	cacheMutex     sync.RWMutex

	// 内存缓存 - 商品状态撤销记录
	undoRecords map[string]*UndoRecord // key: "productID:sellerID"
	undoMutex   sync.RWMutex
}

// NewRecommendService 创建推荐服务实例
func NewRecommendService(
	viewRecordRepo repository.ViewRecordRepository,
	productRepo repository.ProductRepository,
) *RecommendService {
	service := &RecommendService{
		viewRecordRepo: viewRecordRepo,
		productRepo:    productRepo,
		recommendCache: make(map[int64]*RecommendationCache),
		undoRecords:    make(map[string]*UndoRecord),
	}

	// 启动定期清理过期缓存的后台任务
	go service.cleanupExpiredCache()

	return service
}

// cleanupExpiredCache 定期清理过期的缓存数据
func (s *RecommendService) cleanupExpiredCache() {
	ticker := time.NewTicker(1 * time.Minute) // 每分钟清理一次
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		// 清理推荐缓存
		s.cacheMutex.Lock()
		for userID, cache := range s.recommendCache {
			if now.After(cache.ExpiresAt) {
				delete(s.recommendCache, userID)
			}
		}
		s.cacheMutex.Unlock()

		// 清理撤销记录
		s.undoMutex.Lock()
		for key, record := range s.undoRecords {
			if now.After(record.ExpiresAt) {
				delete(s.undoRecords, key)
			}
		}
		s.undoMutex.Unlock()
	}
}

// RecordView 记录用户浏览行为
func (s *RecommendService) RecordView(ctx context.Context, userID int64, productID int64) error {
	return s.viewRecordRepo.AddView(ctx, userID, productID)
}

// GetRecommendations 获取用户推荐商品列表
func (s *RecommendService) GetRecommendations(ctx context.Context, userID int64, maxCount int) ([]int64, error) {
	// 先尝试从缓存获取
	s.cacheMutex.RLock()
	if cache, exists := s.recommendCache[userID]; exists && time.Now().Before(cache.ExpiresAt) {
		s.cacheMutex.RUnlock()
		if len(cache.ProductIDs) > maxCount {
			return cache.ProductIDs[:maxCount], nil
		}
		return cache.ProductIDs, nil
	}
	s.cacheMutex.RUnlock()

	// 缓存未命中，生成推荐
	recommendations, err := s.generateRecommendations(ctx, userID, maxCount)
	if err != nil {
		return nil, err
	}

	// 将结果缓存10分钟
	s.cacheMutex.Lock()
	s.recommendCache[userID] = &RecommendationCache{
		ProductIDs: recommendations,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	}
	s.cacheMutex.Unlock()

	return recommendations, nil
}

// generateRecommendations 生成推荐列表的核心逻辑
func (s *RecommendService) generateRecommendations(ctx context.Context, userID int64, maxCount int) ([]int64, error) {
	// 获取用户最近浏览记录
	recentViews, err := s.viewRecordRepo.ListRecentViews(ctx, userID, 20)
	if err != nil {
		return nil, err
	}

	if len(recentViews) == 0 {
		// 用户无浏览记录，返回热门商品
		return s.getPopularProducts(ctx, maxCount)
	}

	// 基于浏览记录的推荐逻辑
	// 这里可以实现更复杂的推荐算法，目前使用简化版本
	return s.getRecommendationsByViews(ctx, recentViews, userID, maxCount)
}

// getPopularProducts 获取热门商品（用于无浏览记录的用户）
func (s *RecommendService) getPopularProducts(ctx context.Context, maxCount int) ([]int64, error) {
	// 简化实现：使用最新在售商品作为热门商品近似
	products, _, err := s.productRepo.ListLatestForSale(ctx, nil, 1, maxCount)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(products))
	for _, p := range products {
		ids = append(ids, p.ID)
	}
	return ids, nil
}

// getRecommendationsByViews 基于浏览记录生成推荐
func (s *RecommendService) getRecommendationsByViews(ctx context.Context, views []*model.UserRecentView, userID int64, maxCount int) ([]int64, error) {
	// 统计最近浏览商品的标签频次
	tagFreq := make(map[int64]int)
	excludeIDs := make(map[int64]struct{})
	for _, v := range views {
		excludeIDs[v.ProductID] = struct{}{}
		prod, _, tagIDs, err := s.productRepo.GetByID(ctx, v.ProductID)
		if err != nil {
			continue // 忽略个别失败
		}
		// 统计标签
		for _, t := range tagIDs {
			tagFreq[t]++
		}
		// 也避免推荐用户自己的商品
		if prod != nil && prod.SellerID == userID {
			excludeIDs[prod.ID] = struct{}{}
		}
	}

	// 若无标签数据，回退到热门商品
	if len(tagFreq) == 0 {
		return s.getPopularProducts(ctx, maxCount)
	}

	// 将标签按频次排序
	type kv struct {
		tagID int64
		freq  int
	}
	sorted := make([]kv, 0, len(tagFreq))
	for id, f := range tagFreq {
		sorted = append(sorted, kv{tagID: id, freq: f})
	}
	// 简单排序（频次降序）
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].freq > sorted[i].freq {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// 依据标签依次搜索在售商品，合并去重，排除自身商品与已浏览商品
	result := make([]int64, 0, maxCount)
	seen := make(map[int64]struct{})
	for _, it := range sorted {
		if len(result) >= maxCount {
			break
		}
		params := model.SearchParams{
			TagID:    it.tagID,
			Page:     1,
			PageSize: maxCount,
		}
		prods, _, err := s.productRepo.Search(ctx, params)
		if err != nil {
			continue
		}
		for _, p := range prods {
			if len(result) >= maxCount {
				break
			}
			if p.SellerID == userID {
				continue
			}
			if _, ex := excludeIDs[p.ID]; ex {
				continue
			}
			if _, ok := seen[p.ID]; ok {
				continue
			}
			seen[p.ID] = struct{}{}
			result = append(result, p.ID)
		}
	}

	// 不足则用最新在售补齐
	if len(result) < maxCount {
		// 构造排除列表
		exList := make([]int64, 0, len(seen)+len(excludeIDs))
		for id := range seen {
			exList = append(exList, id)
		}
		for id := range excludeIDs {
			exList = append(exList, id)
		}
		latest, _, err := s.productRepo.ListLatestForSale(ctx, exList, 1, maxCount-len(result))
		if err == nil {
			for _, p := range latest {
				if p.SellerID == userID {
					continue
				}
				result = append(result, p.ID)
			}
		}
	}

	return result, nil
}

// GetHomeData 构建首页数据（推荐 + 最新）
func (s *RecommendService) GetHomeData(ctx context.Context, userID *int64, page, pageSize int) (*model.HomeData, error) {
	var recIDs []int64
	if userID != nil {
		ids, err := s.GetRecommendations(ctx, *userID, 5)
		if err != nil {
			// 忽略推荐错误，继续返回最新
		} else {
			recIDs = ids
		}
	}

	// 查询最新在售，排除推荐ID
	latestProds, _, err := s.productRepo.ListLatestForSale(ctx, recIDs, page, pageSize)
	if err != nil {
		return nil, err
	}

	toCard := func(p model.Product) model.ProductCardDTO {
		return model.ProductCardDTO{
			ID:           p.ID,
			Title:        p.Title,
			Price:        p.Price,
			MainImageUrl: p.MainImageURL,
			Status:       p.Status,
			CreatedAt:    p.CreatedAt,
		}
	}

	// 构建推荐商品卡片
	recCards := make([]model.ProductCardDTO, 0, len(recIDs))
	if len(recIDs) > 0 {
		// 简化：逐个加载商品信息
		for _, id := range recIDs {
			prod, _, _, err := s.productRepo.GetByID(ctx, id)
			if err == nil && prod != nil {
				recCards = append(recCards, toCard(*prod))
			}
		}
	}

	// 若推荐不足5条，用最新在售补齐（不重复）
	if len(recCards) < 5 {
		seen := make(map[int64]struct{})
		for _, c := range recCards {
			seen[c.ID] = struct{}{}
		}
		for _, p := range latestProds {
			if len(recCards) >= 5 {
				break
			}
			if _, ok := seen[p.ID]; ok {
				continue
			}
			recCards = append(recCards, toCard(p))
		}
	}

	// 最新列表卡片
	latestCards := make([]model.ProductCardDTO, 0, len(latestProds))
	for _, p := range latestProds {
		latestCards = append(latestCards, toCard(p))
	}

	return &model.HomeData{
		Recommendations: recCards,
		Latest:          latestCards,
	}, nil
}

// SaveUndoRecord 保存撤销记录
func (s *RecommendService) SaveUndoRecord(productID int64, sellerID int64, fromStatus, toStatus string) {
	key := fmt.Sprintf("%d:%d", productID, sellerID)

	s.undoMutex.Lock()
	s.undoRecords[key] = &UndoRecord{
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		ExpiresAt:  time.Now().Add(3 * time.Second), // 3秒过期
	}
	s.undoMutex.Unlock()
}

// GetUndoRecord 获取撤销记录
func (s *RecommendService) GetUndoRecord(productID int64, sellerID int64) (fromStatus, toStatus string, exists bool) {
	key := fmt.Sprintf("%d:%d", productID, sellerID)

	s.undoMutex.RLock()
	defer s.undoMutex.RUnlock()

	record, exists := s.undoRecords[key]
	if !exists || time.Now().After(record.ExpiresAt) {
		return "", "", false
	}

	return record.FromStatus, record.ToStatus, true
}

// DeleteUndoRecord 删除撤销记录
func (s *RecommendService) DeleteUndoRecord(productID int64, sellerID int64) {
	key := fmt.Sprintf("%d:%d", productID, sellerID)

	s.undoMutex.Lock()
	delete(s.undoRecords, key)
	s.undoMutex.Unlock()
}
