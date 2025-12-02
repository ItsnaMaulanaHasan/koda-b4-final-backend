package handlers

import (
	"backend-koda-shortlink/internal/services"
	"backend-koda-shortlink/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserDetail godoc
// @Summary      Get user detail
// @Description  Get specific user detail by ID
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.ResponseSuccess{data=models.User}
// @Failure      404  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /users [get]
func (h *UserHandler) GetUserDetail(c *gin.Context) {
	userId := c.GetInt("userId")

	user, err := h.userService.GetById(c.Request.Context(), userId)
	if err != nil {
		c.JSON(404, response.ResponseError{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(200, response.ResponseSuccess{
		Success: true,
		Message: "Success get user detail",
		Data:    user,
	})
}
