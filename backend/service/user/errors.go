package user

import "errors"

// Error codes for user service
const (
	// ErrCodeNicknameChangeTooSoon is returned when nickname is changed too soon (within 30 days)
	ErrCodeNicknameChangeTooSoon = 1001
	// ErrCodeInvalidAccountFormat is returned when account format is invalid
	ErrCodeInvalidAccountFormat = 1002
	// ErrCodePasswordTooShort is returned when password is too short (less than 8 characters)
	ErrCodePasswordTooShort = 1003
	// ErrCodeAccountExists is returned when account already exists
	ErrCodeAccountExists = 1004
	// ErrCodeInvalidWechatID is returned when wechat ID format is invalid
	ErrCodeInvalidWechatID = 1005
	// ErrCodeInvalidCredentials is returned when login credentials are invalid
	ErrCodeInvalidCredentials = 1006
	// ErrCodeUserNotFound is returned when user is not found
	ErrCodeUserNotFound = 1007
	// ErrCodeInvalidOldPassword is returned when old password is incorrect
	ErrCodeInvalidOldPassword = 1008
)

// ServiceError represents a service layer error
type ServiceError struct {
	Code    int
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	return e.Message
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

// Common errors
var (
	ErrAccountFormat      = errors.New("账号只能包含字母和数字")
	ErrPasswordTooShort   = errors.New("密码长度至少为 8 个字符")
	ErrAccountExists      = errors.New("账号已存在")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrInvalidCredentials = errors.New("账号或密码不正确")
	ErrInvalidOldPassword = errors.New("原密码错误")
	ErrWechatIDFormat     = errors.New("微信号必须为 4-64 个字符，且只可包含字母、数字、下划线或连字符")
)

// NewNicknameChangeTooSoonError creates a new error for nickname change too soon
func NewNicknameChangeTooSoonError(days int) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeNicknameChangeTooSoon,
		Message: "昵称每 30 天只能修改一次",
		Err:     errors.New("昵称修改过于频繁"),
	}
}

// NewInvalidAccountFormatError creates a new error for invalid account format
func NewInvalidAccountFormatError() *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidAccountFormat,
		Message: ErrAccountFormat.Error(),
		Err:     ErrAccountFormat,
	}
}

// NewPasswordTooShortError creates a new error for password too short
func NewPasswordTooShortError() *ServiceError {
	return &ServiceError{
		Code:    ErrCodePasswordTooShort,
		Message: ErrPasswordTooShort.Error(),
		Err:     ErrPasswordTooShort,
	}
}

// NewAccountExistsError creates a new error for account exists
func NewAccountExistsError() *ServiceError {
	return &ServiceError{
		Code:    ErrCodeAccountExists,
		Message: ErrAccountExists.Error(),
		Err:     ErrAccountExists,
	}
}

// NewInvalidWechatIDError creates a new error for invalid wechat ID
func NewInvalidWechatIDError() *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidWechatID,
		Message: ErrWechatIDFormat.Error(),
		Err:     ErrWechatIDFormat,
	}
}

// NewInvalidCredentialsError creates a new error for invalid credentials
func NewInvalidCredentialsError() *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidCredentials,
		Message: ErrInvalidCredentials.Error(),
		Err:     ErrInvalidCredentials,
	}
}

// NewUserNotFoundError creates a new error for user not found
func NewUserNotFoundError() *ServiceError {
	return &ServiceError{
		Code:    ErrCodeUserNotFound,
		Message: ErrUserNotFound.Error(),
		Err:     ErrUserNotFound,
	}
}

// NewInvalidOldPasswordError creates a new error for invalid old password
func NewInvalidOldPasswordError() *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidOldPassword,
		Message: "原密码错误",
		Err:     ErrInvalidOldPassword,
	}
}
