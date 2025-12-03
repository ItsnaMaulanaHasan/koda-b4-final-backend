package handlers

import (
	"backend-koda-shortlink/internal/models"
	"backend-koda-shortlink/internal/services"
	"backend-koda-shortlink/pkg/response"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShortLinkHandler struct {
	service *services.ShortLinkService
}

func NewShortLinkHandler(service *services.ShortLinkService) *ShortLinkHandler {
	return &ShortLinkHandler{service: service}
}

// CreateShortLink godoc
// @Summary      Create short link
// @Description  Create a new short link with auto-generated code
// @Tags         links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  models.CreateShortLinkRequest  true  "Short link details"
// @Success      201  {object}  response.ResponseSuccess
// @Failure      400  {object}  response.ResponseError
// @Failure      401  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /links [post]
func (h *ShortLinkHandler) CreateShortLink(c *gin.Context) {
	userId := c.GetInt("userId")

	var req models.CreateShortLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	link, err := h.service.CreateShortLink(c.Request.Context(), userId, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.ResponseSuccess{
		Success: true,
		Message: "Short link created successfully",
		Data: models.ShortLinkResponse{
			ShortCode:   link.ShortCode,
			OriginalUrl: link.OriginalURL,
			ShortUrl:    os.Getenv("APP_URL") + link.ShortCode,
		},
	})
}

// GetAllLinks godoc
// @Summary      Get all short links
// @Description  Get all short links created by authenticated user with filters
// @Tags         links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page     query  int     false  "Page number" default(1)
// @Param        limit    query  int     false  "Items per page" default(10)
// @Param        search   query  string  false  "Search query"
// @Param        status   query  string  false  "Filter by status (active/inactive)"
// @Success      200  {object}  response.ResponseSuccess
// @Failure      401  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /links [get]
func (h *ShortLinkHandler) GetAllLinks(c *gin.Context) {
	userId := c.GetInt("userId")

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	search := c.Query("search")
	status := c.Query("status")

	links, total, err := h.service.GetUserLinksWithFilter(c.Request.Context(), userId, page, limit, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Error:   "Failed to fetch links",
		})
		return
	}

	appURL := os.Getenv("APP_URL")
	linkResponses := make([]map[string]interface{}, len(links))
	for i, link := range links {
		linkResponses[i] = map[string]any{
			"id":             link.ID,
			"userId":         link.UserID,
			"shortCode":      link.ShortCode,
			"shortUrl":       appURL + link.ShortCode,
			"originalUrl":    link.OriginalURL,
			"isActive":       link.IsActive,
			"clickCount":     link.ClickCount,
			"lastClicked_at": link.LastClickedAt,
			"createdAt":      link.CreatedAt,
			"updatedAt":      link.UpdatedAt,
			"createdBy":      link.CreatedBy,
			"updatedBy":      link.UpdatedBy,
		}
	}

	c.JSON(http.StatusOK, response.ResponseSuccess{
		Success: true,
		Message: "Links retrieved successfully",
		Data: gin.H{
			"links": linkResponses,
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"totalPages": (total + limit - 1) / limit,
			},
		},
	})
}

// GetLinkByShortCode godoc
// @Summary      Get short link by code
// @Description  Get specific short link details by short code
// @Tags         links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        shortCode  path  string  true  "Short code"
// @Success      200  {object}  response.ResponseSuccess
// @Failure      401  {object}  response.ResponseError
// @Failure      403  {object}  response.ResponseError
// @Failure      404  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /links/{shortCode} [get]
func (h *ShortLinkHandler) GetLinkByShortCode(c *gin.Context) {
	userId := c.GetInt("userId")
	shortCode := c.Param("shortCode")

	link, err := h.service.GetLinkByShortCode(c.Request.Context(), shortCode, userId)
	if err != nil {
		if err.Error() == "short link not found" {
			c.JSON(http.StatusNotFound, response.ResponseError{
				Success: false,
				Error:   "Short link not found",
			})
			return
		}
		if err.Error() == "unauthorized access" {
			c.JSON(http.StatusForbidden, response.ResponseError{
				Success: false,
				Error:   "Access denied",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Error:   "Failed to fetch link",
		})
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess{
		Success: true,
		Message: "Link retrieved successfully",
		Data:    link,
	})
}

// UpdateShortLink godoc
// @Summary      Update short link
// @Description  Update short link details (original URL and/or active status)
// @Tags         links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        shortCode  path  string  true  "Short code"
// @Param        request  body  models.UpdateShortLinkRequest  true  "Update details"
// @Success      200  {object}  response.ResponseSuccess
// @Failure      400  {object}  response.ResponseError
// @Failure      401  {object}  response.ResponseError
// @Failure      403  {object}  response.ResponseError
// @Failure      404  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /links/{shortCode} [put]
func (h *ShortLinkHandler) UpdateShortLink(c *gin.Context) {
	userId := c.GetInt("userId")
	shortCode := c.Param("shortCode")

	var req models.UpdateShortLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	link, err := h.service.UpdateShortLink(c.Request.Context(), shortCode, userId, &req)
	if err != nil {
		if err.Error() == "short link not found" {
			c.JSON(http.StatusNotFound, response.ResponseError{
				Success: false,
				Error:   "Short link not found",
			})
			return
		}
		if err.Error() == "unauthorized access" {
			c.JSON(http.StatusForbidden, response.ResponseError{
				Success: false,
				Error:   "Access denied",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Error:   "Failed to update link",
		})
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess{
		Success: true,
		Message: "Link updated successfully",
		Data:    link,
	})
}

// DeleteShortLink godoc
// @Summary      Delete short link
// @Description  Delete short link by short code
// @Tags         links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        shortCode  path  string  true  "Short code"
// @Success      200  {object}  response.ResponseSuccess
// @Failure      401  {object}  response.ResponseError
// @Failure      403  {object}  response.ResponseError
// @Failure      404  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /links/{shortCode} [delete]
func (h *ShortLinkHandler) DeleteShortLink(c *gin.Context) {
	userId := c.GetInt("userId")
	shortCode := c.Param("shortCode")

	err := h.service.DeleteShortLink(c.Request.Context(), shortCode, userId)
	if err != nil {
		if err.Error() == "short link not found" || err.Error() == "short link not found or unauthorized" {
			c.JSON(http.StatusNotFound, response.ResponseError{
				Success: false,
				Error:   "Short link not found",
			})
			return
		}
		if err.Error() == "unauthorized access" {
			c.JSON(http.StatusForbidden, response.ResponseError{
				Success: false,
				Error:   "Access denied",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Error:   "Failed to delete link",
		})
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess{
		Success: true,
		Message: "Link deleted successfully",
		Data:    nil,
	})
}

func (h *ShortLinkHandler) Redirect(c *gin.Context) {
	code := c.Param("shortCode")

	link, err := h.service.ResolveShortCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ResponseError{
			Success: false,
			Error:   "Short link not found",
		})
		return
	}

	go h.service.LogClick(code)

	h.service.SaveClickAnalytics(c.Request, link)

	c.Redirect(http.StatusTemporaryRedirect, link.OriginalURL)
}
