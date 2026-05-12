package services

import (
	"errors"
	"time"

	"github.com/mustafazeren/go-ecommerce-course/internal/config"
	"github.com/mustafazeren/go-ecommerce-course/internal/dto"
	"github.com/mustafazeren/go-ecommerce-course/internal/models"
	"github.com/mustafazeren/go-ecommerce-course/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user exists
	var existingUser models.User
	err := s.db.Where("email = ?", req.Email).First(&existingUser).Error

	if err == nil {
		return nil, errors.New("user already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash pass
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	// Create User
	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("user creation failed")
	}
	// Create cart
	cart := models.Cart{
		UserID: user.ID,
	}
	if err := s.db.Create(&cart).Error; err != nil {
		return nil, errors.New("cart not found")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ? and is_active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, &s.config.JWT)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? and expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}
	s.db.Delete(&refreshToken)
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.db.Where("token=?", refreshToken).Delete(&models.RefreshToken{}).Error
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		&s.config.JWT,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpires),
	}

	if err := s.db.Create(&refreshTokenModel).Error; err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
