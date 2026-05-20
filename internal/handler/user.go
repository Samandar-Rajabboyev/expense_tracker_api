package handler

import (
	"expense-tracker-api/internal/jwt"
	"expense-tracker-api/internal/response"
	"expense-tracker-api/internal/service"

	// "net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service    *service.UserService
	jwtService *jwt.JWTService
}

func NewUserHandler(s *service.UserService, jwtService *jwt.JWTService) *UserHandler {
	return &UserHandler{service: s, jwtService: jwtService}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	user, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 201, user)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	user, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, 200, map[string]any{"token": token, "refresh_token": refreshToken, "user": user})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *UserHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	claims, err := h.jwtService.ValidationToken(req.RefreshToken)
	if err != nil {
		response.Error(c, 401, "invalid refresh token")
		return
	}

	newToken, err := h.jwtService.GenerateToken(claims.UserID, claims.Email)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, 200, gin.H{"token": newToken})
}
