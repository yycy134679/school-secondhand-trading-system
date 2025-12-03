package category

import (
	"context"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
)

// CategoryService 分类服务接口
type CategoryService interface {
	// ListCategories 获取所有分类，供前台使用
	ListCategories(ctx context.Context) ([]*model.Category, error)
	// CreateCategory 创建分类
	CreateCategory(ctx context.Context, category *model.Category) error
	// UpdateCategory 更新分类
	UpdateCategory(ctx context.Context, category *model.Category) error
	// DeleteCategory 删除分类，删除前检查引用
	DeleteCategory(ctx context.Context, id int64) error
}

// categoryService 分类服务实现
type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// ListCategories 获取所有分类
func (s *categoryService) ListCategories(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.categoryRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	// 转换为指针切片
	result := make([]*model.Category, len(categories))
	for i := range categories {
		result[i] = &categories[i]
	}
	return result, nil
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(ctx context.Context, category *model.Category) error {
	return s.categoryRepo.Create(ctx, category)
}

// UpdateCategory 更新分类
func (s *categoryService) UpdateCategory(ctx context.Context, category *model.Category) error {
	return s.categoryRepo.Update(ctx, category)
}

// DeleteCategory 删除分类，删除前检查引用
func (s *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	// 检查是否有关联的商品
	count, err := s.categoryRepo.CountProductsByCategory(ctx, id)
	if err != nil {
		return err
	}
	// 如果有关联的商品，返回错误码4001
	if count > 0 {
		return ErrCategoryHasProducts
	}
	return s.categoryRepo.Delete(ctx, id)
}
