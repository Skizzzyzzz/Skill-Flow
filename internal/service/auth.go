package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vern/skillflow/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	deps ServicesDeps
}

func NewAuthService(deps ServicesDeps) *AuthService {
	return &AuthService{deps: deps}
}

type RegisterInput struct {
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*TokenPair, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		IsVerified:   false,
		Role:         "user",
	}

	if err := s.deps.Repos.User.Create(ctx, user); err != nil {
		return nil, err
	}

	profile := &models.Profile{
		UserID:      user.ID,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		DisplayName: input.Username,
	}

	if err := s.deps.Repos.Profile.Create(ctx, profile); err != nil {
		return nil, err
	}

	return s.generateTokenPair(user.ID, user.Role)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.deps.Repos.User.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	now := time.Now()
	user.LastLoginAt = &now
	s.deps.Repos.User.Update(ctx, user)

	return s.generateTokenPair(user.ID, user.Role)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.deps.Config.Auth.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	return s.generateTokenPair(userID, role)
}

func (s *AuthService) GetOIDCAuthURL() string {
	// Implementation for OIDC authentication URL
	return s.deps.Config.Auth.OIDC.IssuerURL + "/protocol/openid-connect/auth?client_id=" + s.deps.Config.Auth.OIDC.ClientID
}

func (s *AuthService) OIDCCallback(ctx context.Context, code string) (*TokenPair, error) {
	// Implementation for OIDC callback
	// This would exchange the code for tokens and create/update user
	return nil, errors.New("not implemented")
}

func (s *AuthService) generateTokenPair(userID uint, role string) (*TokenPair, error) {
	accessToken, err := s.generateAccessToken(userID, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(userID, role)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.deps.Config.Auth.JWTExpiry.Seconds()),
	}, nil
}

func (s *AuthService) generateAccessToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "access",
		"exp":     time.Now().Add(s.deps.Config.Auth.JWTExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.deps.Config.Auth.JWTSecret))
}

func (s *AuthService) generateRefreshToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "refresh",
		"exp":     time.Now().Add(s.deps.Config.Auth.RefreshExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.deps.Config.Auth.JWTSecret))
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
