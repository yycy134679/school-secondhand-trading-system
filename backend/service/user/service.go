package user

import (
	"context"
	"regexp"
	"time"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/auth"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
	"gorm.io/gorm"
)

// UserService handles user business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// UserRegisterRequest represents the request for user registration
type UserRegisterRequest struct {
	Account  string  `json:"account"`
	Nickname string  `json:"nickname"`
	Password string  `json:"password"`
	WechatID *string `json:"wechatId,omitempty"`
}

// UserLoginRequest represents the request for user login
type UserLoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// UserUpdateProfileRequest represents the request for updating user profile
type UserUpdateProfileRequest struct {
	Nickname  string  `json:"nickname"`
	AvatarURL string  `json:"avatarUrl"`
	WechatID  *string `json:"wechatId,omitempty"`
}

// UserChangePasswordRequest represents the request for changing password
type UserChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// UserResponse represents the response for user information
type UserResponse struct {
	ID        int64   `json:"id"`
	Account   string  `json:"account"`
	Nickname  string  `json:"nickname"`
	AvatarUrl string  `json:"avatarUrl"`
	IsAdmin   bool    `json:"isAdmin"`
	WechatID  *string `json:"wechatId,omitempty"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// Register registers a new user and returns authentication response
func (s *UserService) Register(ctx context.Context, account, nickname, password string, wechatID *string) (*AuthResponse, error) {
	// Validate account format (only letters and numbers)
	if !s.isValidAccountFormat(account) {
		return nil, NewInvalidAccountFormatError()
	}

	// Validate password length
	if len(password) < 8 {
		return nil, NewPasswordTooShortError()
	}

	// Check if account already exists
	existingUser, err := s.userRepo.GetByAccount(ctx, account)
	if err != nil {
		// 如果是记录不存在错误，这是正常的，表示账号可用
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		// 记录不存在，继续执行
	} else if existingUser != nil {
		// 记录存在，账号已被使用
		return nil, NewAccountExistsError()
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user model
	now := time.Now()
	user := &model.User{
		Account:   account,
		Nickname:  nickname,
		Password:  hashedPassword,
		IsAdmin:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	// 设置WechatID
	if wechatID != nil {
		user.WechatID = *wechatID
	}

	// Insert user
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := auth.GenerateToken(int64(user.ID))
	if err != nil {
		return nil, err
	}

	// Create response
	return &AuthResponse{
		User:  s.buildUserResponse(user),
		Token: token,
	}, nil
}

// isValidAccountFormat checks if the account format is valid (only letters and numbers)
func (s *UserService) isValidAccountFormat(account string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, account)
	return matched
}

// buildUserResponse converts user model to response
func (s *UserService) buildUserResponse(user *model.User) UserResponse {
	var wechatID *string
	if user.WechatID != "" {
		temp := user.WechatID
		wechatID = &temp
	}
	return UserResponse{
		ID:        int64(user.ID),
		Account:   user.Account,
		Nickname:  user.Nickname,
		AvatarUrl: user.AvatarUrl,
		IsAdmin:   user.IsAdmin,
		WechatID:  wechatID,
	}
}

// Login authenticates a user and returns authentication response
func (s *UserService) Login(ctx context.Context, account, password string, remember bool) (*AuthResponse, error) {
	// Get user by account
	user, err := s.userRepo.GetByAccount(ctx, account)
	if err != nil {
		// 如果是记录不存在错误，返回无效凭证错误
		if err == gorm.ErrRecordNotFound {
			return nil, NewInvalidCredentialsError()
		}
		return nil, err
	}
	if user == nil {
		return nil, NewInvalidCredentialsError()
	}

	// Verify password
	err = auth.ComparePassword(user.Password, password)
	if err != nil {
		return nil, NewInvalidCredentialsError()
	}

	// Generate token
	token, err := auth.GenerateToken(int64(user.ID))
	if err != nil {
		return nil, err
	}

	// Create response
	return &AuthResponse{
		User:  s.buildUserResponse(user),
		Token: token,
	}, nil
}

// GetProfile returns the user profile by user ID
func (s *UserService) GetProfile(ctx context.Context, userID uint) (*UserResponse, error) {
	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, NewUserNotFoundError()
	}

	// Convert to response
	response := s.buildUserResponse(user)
	return &response, nil
}

// UpdateProfile updates the user profile
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, nickname, avatarURL string, wechatID *string) (*UserResponse, error) {
	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, NewUserNotFoundError()
	}

	// Check nickname change interval if nickname is being changed
	if nickname != "" && nickname != user.Nickname {
		// Check if user can change nickname (at least 30 days since last change)
		if user.LastNicknameChangedAt != nil && !user.LastNicknameChangedAt.IsZero() {
			thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
			if !user.LastNicknameChangedAt.Before(thirtyDaysAgo) {
				return nil, NewNicknameChangeTooSoonError(30)
			}
		}
		// Update last nickname changed time
		now := time.Now()
		user.LastNicknameChangedAt = &now
		user.Nickname = nickname
	}

	// Update avatar URL if provided
	if avatarURL != "" {
		user.AvatarUrl = avatarURL
	}

	// Validate and update wechat ID if provided
	if wechatID != nil {
		if !s.isValidWechatID(*wechatID) {
			return nil, NewInvalidWechatIDError()
		}
		user.WechatID = *wechatID
	}

	// Update user
	err = s.userRepo.UpdateProfile(ctx, user)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := s.buildUserResponse(user)
	return &response, nil
}

// 不再需要这个方法，逻辑已内联到UpdateProfile中

// isValidWechatID checks if the wechat ID format is valid
func (s *UserService) isValidWechatID(wechatID string) bool {
	// Wechat ID must be 4-64 characters, containing letters, numbers, underscores or hyphens
	if len(wechatID) < 4 || len(wechatID) > 64 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, wechatID)
	return matched
}

// ChangePassword changes the user's password
func (s *UserService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) (*UserResponse, error) {
	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, NewUserNotFoundError()
	}

	// Verify old password
	err = auth.ComparePassword(user.Password, oldPassword)
	if err != nil {
		return nil, NewInvalidOldPasswordError()
	}

	// Validate new password length
	if len(newPassword) < 8 {
		return nil, NewPasswordTooShortError()
	}

	// Hash new password
	newHashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	// Update password
	user.Password = newHashedPassword
	now := time.Now()
	user.UpdatedAt = now

	// Save updated user
	err = s.userRepo.UpdatePassword(ctx, int64(userID), newHashedPassword)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := s.buildUserResponse(user)
	return &response, nil
}
