package transportauth

import (
	domainauth "go-bootstrap/internal/domain/auth"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/http/server/ginx"
	"github.com/gin-gonic/gin"
)

type AuthRestAPIHandler struct {
	authService domainauth.AuthService
	helper      *ginx.GinHelper
}

func NewRestAPIHandler(authService domainauth.AuthService, helper *ginx.GinHelper) *AuthRestAPIHandler {
	return &AuthRestAPIHandler{
		authService: authService,
		helper:      helper,
	}
}

// User login
// (POST /api/v1/auth/login)
func (h *AuthRestAPIHandler) ApiV1PostAuthLogin(c *gin.Context) {
	// TODO: Implement login handler
}

// User logout
// (POST /api/v1/auth/logout)
func (h *AuthRestAPIHandler) ApiV1PostAuthLogout(c *gin.Context) {
	// TODO: Implement logout handler
}

// Refresh access token
// (POST /api/v1/auth/refresh)
func (h *AuthRestAPIHandler) ApiV1PostAuthRefresh(c *gin.Context) {
	// TODO: Implement refresh token handler
}
