package repository

import (
	"context"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
)

// UserRepository defines user data access behavior. Implement with GORM.
type UserRepository interface {
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*model.User, error)
}
