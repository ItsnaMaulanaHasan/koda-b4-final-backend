package handlers

import (
	"backend-koda-shortlink/internal/services"
	"backend-koda-shortlink/pkg/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// Stats godoc
// @Summary      Get dashboard statistics
// @Description  Retrieve user-specific statistics for dashboard overview
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.ResponseSuccess{data=services.DashboardStats}
// @Failure      500  {object}  response.ResponseError
// @Router       /dashboard/stats [get]
func (h *DashboardHandler) Stats(c *gin.Context) {
	userId := c.GetInt("userId")

	data, err := h.dashboardService.Stats(c.Request.Context(), userId)
	if err != nil {
		c.JSON(500, response.ResponseError{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, response.ResponseSuccess{
		Success: true,
		Message: "Success get data statistic",
		Data:    data,
	})
}
