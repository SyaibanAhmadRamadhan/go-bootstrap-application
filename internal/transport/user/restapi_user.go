package transportuser

import (
	domainuser "go-bootstrap/internal/domain/user"
	"go-bootstrap/internal/gen/restapigen"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/http/server/ginx"
	"github.com/gin-gonic/gin"
)

type UserRestAPIHandler struct {
	userService domainuser.UserService
	helper      *ginx.GinHelper
}

func NewRestAPIHandler(userService domainuser.UserService, helper *ginx.GinHelper) *UserRestAPIHandler {
	return &UserRestAPIHandler{
		userService: userService,
		helper:      helper,
	}
}

// Get list of users
// (GET /api/v1/users)
func (h *UserRestAPIHandler) ApiV1GetUsers(c *gin.Context, params restapigen.ApiV1GetUsersParams) {
	// TODO: Implement get users list handler
}

// Change password
// (POST /api/v1/users/change-password)
func (h *UserRestAPIHandler) ApiV1PostUsersChangePassword(c *gin.Context) {
	// TODO: Implement change password handler
}

// Get user profile
// (GET /api/v1/users/profile)
func (h *UserRestAPIHandler) ApiV1GetUsersProfile(c *gin.Context) {
	// TODO: Implement get user profile handler
}

// Update user profile
// (PUT /api/v1/users/profile)
func (h *UserRestAPIHandler) ApiV1PutUsersProfile(c *gin.Context) {
	// TODO: Implement update user profile handler
}

// Register new user
// (POST /api/v1/users/register)
func (h *UserRestAPIHandler) ApiV1PostUsersRegister(c *gin.Context) {
	// TODO: Implement user registration handler
}

// Update user status
// (PUT /api/v1/users/{user_id}/status)
func (h *UserRestAPIHandler) ApiV1PutUsersStatus(c *gin.Context, userId string) {
	// TODO: Implement update user status handler
}
