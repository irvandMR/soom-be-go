package handler

import (
	"net/http"
	"strings"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/repository"
	"soom-be-go/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	usecase *usecase.AuthUsecase
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	repo := repository.NewAuthRepository(db)
	uc := usecase.NewAuthUsecase(repo)
	return &AuthHandler{usecase: uc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	if err := h.usecase.Register(req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse("Register success", nil))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	data, err := h.usecase.Login(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Login success", data))
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	data, err := h.usecase.RefreshToken(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Token refreshed", data))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}

	authHeader := c.GetHeader("Authorization")
	accessToken := ""
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			accessToken = parts[1]
		}
	}

	if err := h.usecase.Logout(accessToken, req.RefreshToken); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Logout success", nil))
}

func (h *AuthHandler) me(c *gin.Context) {
	user := c.MustGet("username").(string)
	c.JSON(http.StatusOK, SuccessResponse("User info", user))
}
