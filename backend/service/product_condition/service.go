package productcondition

import (
	"context"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
)

// Service 新旧程度服务接口
type Service interface {
	// ListProductConditions 获取所有新旧程度
	ListProductConditions(ctx context.Context) ([]*model.ProductCondition, error)
}

type service struct {
	repo repository.ProductConditionRepository
}

// NewService 创建服务实例
func NewService(repo repository.ProductConditionRepository) Service {
	return &service{repo: repo}
}

// ListProductConditions 获取所有新旧程度
func (s *service) ListProductConditions(ctx context.Context) ([]*model.ProductCondition, error) {
	conditions, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.ProductCondition, len(conditions))
	for i := range conditions {
		result[i] = &conditions[i]
	}
	return result, nil
}
