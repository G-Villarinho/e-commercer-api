package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/OVillas/e-commercer-api/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrHashingPassword       = errors.New("failed to hash password")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrEmailNotConfirmed     = errors.New("email not confirmed")
	ErrNameIsSame            = errors.New("name is same as oldName")
	ErrUserNotFoundInContext = errors.New("user not found in context")
	ErrInvalidOldPassword    = errors.New("invalid old password")
	ErrPasswordIsSame        = errors.New("new password is same as old password")
	ErrEmailAlreadyConfirmed = errors.New("email already confirmed")
	ErrOTPExpires            = errors.New("otp expires")
)

type User struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey;column:id"`
	Name           string         `gorm:"size:100;not null;column:name"`
	Username       string         `gorm:"uniqueIndex;size:100;not null;column:username"`
	Email          string         `gorm:"uniqueIndex;size:100;not null;column:email"`
	PasswordHash   string         `gorm:"size:255;not null;column:passwordHash"`
	EmailConfirmed bool           `gorm:"not null;default:false;column:emailConfirmed"`
	AvatarURL      string         `gorm:"size:255;column:AvatarUrl"`
	CreatedAt      time.Time      `gorm:"column:createdAt"`
	UpdatedAt      time.Time      `gorm:"column:updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index;column:deletedAt"`
}

type UserPayLoad struct {
	Name            string `json:"name" validate:"required,min=1,max=75"`
	Email           string `json:"email" validate:"required,email"`
	ConfirmEmail    string `json:"confirmEmail" validate:"required,email,eqfield=Email"`
	Password        string `json:"password,omitempty" validate:"required,min=8,max=255,containsany=!@#&?"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}

type SignInPayLoad struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateNamePayload struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdatePasswordPayload struct {
	OldPassword     string `json:"oldPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=255,containsany=!@#&?"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
}

type ResendCodePayload struct {
	Email string `json:"email" validate:"required,email"`
}

type ConfirmEmailPayload struct {
	OTP string `json:"otp" validate:"required,numeric,min=6,max=6"`
}

type UserHandler interface {
	Create(ctx echo.Context) error
	SignIn(ctx echo.Context) error
	UpdateName(ctx echo.Context) error
	UpdatePassword(ctx echo.Context) error
	GetUserInfo(ctx echo.Context) error
	ResendCode(ctx echo.Context) error
	ConfirmEmail(ctx echo.Context) error
}

type UserService interface {
	Create(ctx context.Context, user UserPayLoad) error
	SignIn(ctx context.Context, signInPayload SignInPayLoad) (*SessionResponse, error)
	UpdateName(ctx context.Context, name string) error
	UpdatePassword(ctx context.Context, updatePasswordPayload UpdatePasswordPayload) error
	GetUserInfo(ctx context.Context) (*UserResponse, error)
	ResendCode(ctx context.Context, resendCodePayload ResendCodePayload) error
	ConfirmEmail(ctx context.Context, confirmEmailPayload ConfirmEmailPayload) error
}

type UserRepository interface {
	Create(ctx context.Context, user User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateName(ctx context.Context, id uuid.UUID, name string) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	UpdateConfirmEmail(ctx context.Context, id uuid.UUID) error
}

func (u *UserPayLoad) trim() {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	u.ConfirmEmail = strings.TrimSpace(u.ConfirmEmail)
}

func (s *SignInPayLoad) trim() {
	s.Email = strings.TrimSpace(s.Email)
}

func (u *UpdateNamePayload) trim() {
	u.Name = strings.TrimSpace(u.Name)
}

func (r *ResendCodePayload) trim() {
	r.Email = strings.TrimSpace(r.Email)
}

func (u *UserPayLoad) Validate() error {
	u.trim()
	validate := validator.New()
	return validate.Struct(u)
}

func (s *SignInPayLoad) Validate() error {
	s.trim()
	validate := validator.New()
	return validate.Struct(s)
}

func (u *UpdatePasswordPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UpdateNamePayload) Validate() error {
	u.trim()
	validate := validator.New()
	return validate.Struct(u)
}

func (r *ResendCodePayload) Validate() error {
	r.trim()
	validate := validator.New()
	return validate.Struct(r)
}

func (u *ConfirmEmailPayload) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("numeric", util.IsNumeric)
	return validate.Struct(u)
}

func (u *UserPayLoad) ToUser(passwordHash string) *User {
	return &User{
		ID:           uuid.New(),
		Name:         u.Name,
		Email:        u.Email,
		Username:     u.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
	}
}

func (User) TableName() string {
	return "User"
}
