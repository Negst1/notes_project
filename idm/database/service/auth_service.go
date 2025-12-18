package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"idm/database/models"
	"idm/database/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, register *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, login *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, token string) (*models.AuthResponse, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, tokenString string) (*jwt.Token, error)
	GetUserFromToken(ctx context.Context, token string) (*models.User, error)
}

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type authService struct {
	repo     repository.AuthRepository
	repoUser repository.UserRepository
	config   JWTConfig
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(repo repository.AuthRepository, repoUser repository.UserRepository, config JWTConfig) AuthService {
	return &authService{
		repo:     repo,
		repoUser: repoUser,
		config:   config,
	}
}

func (r *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	_, err := r.repoUser.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:   req.Username,
		Email:      req.Email,
		IsActive:   true,
		IsVerified: false,
	}

	err = user.HashPassword(user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return r.generateAuthResponse(ctx, user)
}

func (r *authService) Login(ctx context.Context, login *models.LoginRequest) (*models.AuthResponse, error) {
	return nil, nil
}

func (r *authService) RefreshToken(ctx context.Context, token string) (*models.AuthResponse, error) {
	return nil, nil
}
func (r *authService) Logout(ctx context.Context, token string) error {
	return nil
}
func (r *authService) ValidateToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	return nil, nil
}
func (r *authService) GetUserFromToken(ctx context.Context, token string) (*models.User, error) {
	return nil, nil
}
func (s *authService) generateAuthResponse(ctx context.Context, user *models.User) (*models.AuthResponse, error) {

	accessToken, _, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExp, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveRefreshToken(ctx, user.ID, refreshToken, refreshExp); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.config.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *authService) generateAccessToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.config.AccessTokenTTL)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "idm-api",
			Subject:   strconv.Itoa(user.ID),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *authService) generateRefreshToken() (string, time.Time, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(s.config.RefreshTokenTTL)
	token := hex.EncodeToString(bytes)

	return token, expiresAt, nil
}
