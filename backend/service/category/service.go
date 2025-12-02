package category

import (
	"context"
	"errors"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"

	"gorm.io/gorm"
)

// CategoryService 分类服务接口
type CategoryService interface {
	// ListCategories 获取所有分类，供前台使用
	ListCategories(ctx context.Context) ([]*model.CategoryDTO, error)
	// CreateCategory 创建分类
	CreateCategory(ctx context.Context, name string) (*model.CategoryDTO, error)
	// UpdateCategory 更新分类
	UpdateCategory(ctx context.Context, id int64, name string) (*model.CategoryDTO, error)
	// DeleteCategory 删除分类，删除前检查引用，多于0时返回4001错误码
	DeleteCategory(ctx context.Context, id int64) error
	// GetCategoryByID 根据ID获取分类
	GetCategoryByID(ctx context.Context, id int64) (*model.CategoryDTO, error)
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

// ListCategories 获取所有分类，供前台使用
func (s *categoryService) ListCategories(ctx context.Context) ([]*model.CategoryDTO, error) {
	categories, err := s.categoryRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为DTO列表
	dtos := make([]*model.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = category.ToDTO()
	}

	return dtos, nil
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(ctx context.Context, name string) (*model.CategoryDTO, error) {
	// 验证名称
	if name == "" {
		return nil, errors.New("分类名称不能为空")
	}

	category := &model.Category{
		Name: name,
	}

	err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return category.ToDTO(), nil
}

// UpdateCategory 更新分类
func (s *categoryService) UpdateCategory(ctx context.Context, id int64, name string) (*model.CategoryDTO, error) {
	// 验证名称
	if name == "" {
		return nil, errors.New("分类名称不能为空")
	}

	// 检查分类是否存在
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}

	// 更新名称
	category.Name = name
	err = s.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	return category.ToDTO(), nil
}

// DeleteCategory 删除分类，删除前检查引用，多于0时返回4001错误码
func (s *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	// 检查分类是否存在
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		return err
	}

	// 检查是否有商品引用该分类
	count, err := s.categoryRepo.CountProductsByCategory(ctx, id)
	if err != nil {
		return err
	}

	// 如果有商品引用，返回错误码4001
	if count > 0 {
		return errors.New("该分类下存在商品，无法删除")
	}

	// 删除分类
	return s.categoryRepo.Delete(ctx, id)
}

// GetCategoryByID 根据ID获取分类
func (s *categoryService) GetCategoryByID(ctx context.Context, id int64) (*model.CategoryDTO, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}

	return category.ToDTO(), nil
}
