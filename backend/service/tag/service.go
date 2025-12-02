package tag

import (
	"context"
	"errors"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"

	"gorm.io/gorm"
)

// TagService 标签服务接口
type TagService interface {
	// ListTags 获取所有标签，供前台使用
	ListTags(ctx context.Context) ([]*model.TagDTO, error)
	// CreateTag 创建标签
	CreateTag(ctx context.Context, name string) (*model.TagDTO, error)
	// UpdateTag 更新标签
	UpdateTag(ctx context.Context, id int64, name string) (*model.TagDTO, error)
	// DeleteTag 删除标签，删除前检查引用，多于0时返回4002错误码
	DeleteTag(ctx context.Context, id int64) error
	// GetTagByID 根据ID获取标签
	GetTagByID(ctx context.Context, id int64) (*model.TagDTO, error)
	// GetTagsByIDs 根据ID列表获取标签
	GetTagsByIDs(ctx context.Context, ids []int64) ([]*model.TagDTO, error)
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

// ListTags 获取所有标签，供前台使用
func (s *tagService) ListTags(ctx context.Context) ([]*model.TagDTO, error) {
	tags, err := s.tagRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为DTO列表
	dtos := make([]*model.TagDTO, len(tags))
	for i, tag := range tags {
		dtos[i] = tag.ToDTO()
	}

	return dtos, nil
}

// CreateTag 创建标签
func (s *tagService) CreateTag(ctx context.Context, name string) (*model.TagDTO, error) {
	// 验证名称
	if name == "" {
		return nil, errors.New("标签名称不能为空")
	}

	tag := &model.Tag{
		Name: name,
	}

	err := s.tagRepo.Create(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag.ToDTO(), nil
}

// UpdateTag 更新标签
func (s *tagService) UpdateTag(ctx context.Context, id int64, name string) (*model.TagDTO, error) {
	// 验证名称
	if name == "" {
		return nil, errors.New("标签名称不能为空")
	}

	// 检查标签是否存在
	tag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, err
	}

	// 更新名称
	tag.Name = name
	err = s.tagRepo.Update(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag.ToDTO(), nil
}

// DeleteTag 删除标签，删除前检查引用，多于0时返回4002错误码
func (s *tagService) DeleteTag(ctx context.Context, id int64) error {
	// 检查标签是否存在
	_, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return err
	}

	// 检查是否有商品引用该标签
	count, err := s.tagRepo.CountProductsByTag(ctx, id)
	if err != nil {
		return err
	}

	// 如果有商品引用，返回错误码4002
	if count > 0 {
		return errors.New("该标签下存在商品，无法删除")
	}

	// 删除标签
	return s.tagRepo.Delete(ctx, id)
}

// GetTagByID 根据ID获取标签
func (s *tagService) GetTagByID(ctx context.Context, id int64) (*model.TagDTO, error) {
	tag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, err
	}

	return tag.ToDTO(), nil
}

// GetTagsByIDs 根据ID列表获取标签
func (s *tagService) GetTagsByIDs(ctx context.Context, ids []int64) ([]*model.TagDTO, error) {
	if len(ids) == 0 {
		return []*model.TagDTO{}, nil
	}

	tags, err := s.tagRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	// 转换为DTO列表
	dtos := make([]*model.TagDTO, len(tags))
	for i, tag := range tags {
		dtos[i] = tag.ToDTO()
	}

	return dtos, nil
}
