package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vern/skillflow/internal/config"
	"github.com/vern/skillflow/internal/service"
	"github.com/vern/skillflow/pkg/logger"
)

type AuthHandler struct {
	services *service.Services
	config   *config.Config
	logger   *logger.Logger
}

func NewAuthHandler(services *service.Services, cfg *config.Config, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		services: services,
		config:   cfg,
		logger:   log,
	}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.services.Auth.Register(c.Request.Context(), service.RegisterInput{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		h.logger.Error("Failed to register user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.services.Auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to login", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.services.Auth.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh token", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) OIDCLogin(c *gin.Context) {
	authURL := h.services.Auth.GetOIDCAuthURL()
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *AuthHandler) OIDCCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No authorization code provided"})
		return
	}

	tokens, err := h.services.Auth.OIDCCallback(c.Request.Context(), code)
	if err != nil {
		h.logger.Error("OIDC callback failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
