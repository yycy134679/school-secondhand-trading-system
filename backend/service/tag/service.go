package tag

import (
	"context"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
)

// TagService 标签服务接口
type TagService interface {
	// ListTags 获取所有标签，供前台使用
	ListTags(ctx context.Context) ([]*model.Tag, error)
	// CreateTag 创建标签
	CreateTag(ctx context.Context, tag *model.Tag) error
	// UpdateTag 更新标签
	UpdateTag(ctx context.Context, tag *model.Tag) error
	// DeleteTag 删除标签，删除前检查引用
	DeleteTag(ctx context.Context, id int64) error
}

// tagService 标签服务实现
type tagService struct {
	tagRepo repository.TagRepository
}

// NewTagService 创建标签服务实例
func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

// ListTags 获取所有标签
func (s *tagService) ListTags(ctx context.Context) ([]*model.Tag, error) {
	tags, err := s.tagRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	// 转换为指针切片
	result := make([]*model.Tag, len(tags))
	for i := range tags {
		result[i] = &tags[i]
	}
	return result, nil
}

// CreateTag 创建标签
func (s *tagService) CreateTag(ctx context.Context, tag *model.Tag) error {
	return s.tagRepo.Create(ctx, tag)
}

// UpdateTag 更新标签
func (s *tagService) UpdateTag(ctx context.Context, tag *model.Tag) error {
	return s.tagRepo.Update(ctx, tag)
}

// DeleteTag 删除标签，删除前检查引用
func (s *tagService) DeleteTag(ctx context.Context, id int64) error {
	// 检查是否有关联的商品
	count, err := s.tagRepo.CountProductsByTag(ctx, id)
	if err != nil {
		return err
	}
	// 如果有关联的商品，返回错误码4002
	if count > 0 {
		return ErrTagHasProducts
	}
	return s.tagRepo.Delete(ctx, id)
}
