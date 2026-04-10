package handler

import (
	"errors"
	"net/http"
	"strings"

	authSvc "kbmanage/backend/internal/service/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	svc *authSvc.LoginService
}

func NewAuthHandler(svc *authSvc.LoginService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type authUserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName,omitempty"`
}

type loginResponse struct {
	AccessToken  string           `json:"accessToken"`
	RefreshToken string           `json:"refreshToken"`
	ExpiresIn    int64            `json:"expiresIn"`
	User         authUserResponse `json:"user"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Login(c.Request.Context(), authSvc.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		writeAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, toLoginResponse(result))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Refresh(c.Request.Context(), authSvc.RefreshInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		writeAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, toLoginResponse(result))
}

func toLoginResponse(result *authSvc.LoginResult) loginResponse {
	if result == nil {
		return loginResponse{}
	}

	return loginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		User: authUserResponse{
			ID:          result.User.ID,
			Username:    result.User.Username,
			DisplayName: result.User.DisplayName,
		},
	}
}

func writeAuthError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, authSvc.ErrInvalidCredentials),
		errors.Is(err, authSvc.ErrInvalidRefresh),
		errors.Is(err, authSvc.ErrUserDisabled),
		errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusUnauthorized
	case errors.Is(err, gorm.ErrInvalidDB):
		status = http.StatusNotImplemented
	case strings.Contains(strings.ToLower(err.Error()), "required"):
		status = http.StatusBadRequest
	}
	c.JSON(status, gin.H{"error": err.Error()})
}
