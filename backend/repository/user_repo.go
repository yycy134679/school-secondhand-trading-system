package repository

import (
	"context"

	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"gorm.io/gorm"
)

// UserRepository defines user data access behavior. Implement with GORM.
type UserRepository interface {
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*model.User, error)
	// GetByAccount retrieves a user by account
	GetByAccount(ctx context.Context, account string) (*model.User, error)
	// Create creates a new user
	Create(ctx context.Context, user *model.User) error
	// UpdateProfile updates user profile (nickname/avatar/wechat)
	UpdateProfile(ctx context.Context, user *model.User) error
	// UpdatePassword updates user password
	UpdatePassword(ctx context.Context, userID int64, newHash string) error
}

// userRepo implements UserRepository
// It uses GORM for database operations

type userRepo struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository implementation
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// GetByID retrieves a user by ID
func (r *userRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByAccount retrieves a user by account
func (r *userRepo) GetByAccount(ctx context.Context, account string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("account = ?", account).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create creates a new user
func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// UpdateProfile updates user profile (nickname/avatar/wechat)
func (r *userRepo) UpdateProfile(ctx context.Context, user *model.User) error {
	updates := map[string]interface{}{
		"avatar_url": user.AvatarUrl,
		"wechat_id":  user.WechatID,
	}

	// 如果更新了昵称，同时更新最后昵称修改时间
	updates["nickname"] = user.Nickname
	if user.LastNicknameChangedAt != nil {
		updates["last_nickname_changed_at"] = user.LastNicknameChangedAt
	}

	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(updates).Error
}

// UpdatePassword updates user password
func (r *userRepo) UpdatePassword(ctx context.Context, userID int64, newHash string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("password_hash", newHash).Error
}
