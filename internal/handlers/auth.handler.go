package handlers

import (
	"backend-koda-shortlink/internal/models"
	"backend-koda-shortlink/internal/services"
	"backend-koda-shortlink/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new user with a unique email
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        fullname  formData  string  true  "Full name user"
// @Param        email     formData  string  true  "Email user"
// @Param        password  formData  string  true  "Password user" format(password)
// @Success      201  {object}  response.ResponseSuccess{data=models.User}
// @Failure      400  {object}  response.ResponseError
// @Failure      409  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindWith(&req, binding.Form); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Error:   "Please provide valid registration information",
		})
		return
	}

	user, err := h.authService.Register(ctx.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already registered" {
			statusCode = http.StatusConflict
		}

		ctx.JSON(statusCode, response.ResponseError{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.ResponseSuccess{
		Success: true,
		Message: "User registered successfully",
		Data: gin.H{
			"id":       user.Id,
			"fullName": user.FullName,
			"email":    user.Email,
		},
	})
}

// Login godoc
// @Summary      Login user
// @Description  Log in with existing email data
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        email     formData  string  true  "User email"
// @Param        password  formData  string  true  "User password" format(password)
// @Success      200  {object}  response.ResponseSuccess{data=models.LoginResponse}
// @Failure      400  {object}  response.ResponseError
// @Failure      401  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindWith(&req, binding.Form); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Error:   "Please provide valid email and password",
		})
		return
	}

	ipAddress := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	loginResp, err := h.authService.Login(ctx.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "wrong email or password" {
			statusCode = http.StatusUnauthorized
		}

		ctx.JSON(statusCode, response.ResponseError{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.ResponseSuccess{
		Success: true,
		Message: "Login successful!",
		Data:    loginResp,
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Get new access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refreshToken  body  models.RefreshTokenRequest  true  "Refresh Token"
// @Success      200  {object}  response.ResponseSuccess{data=object{access_token=string}}
// @Failure      400  {object}  response.ResponseError
// @Failure      401  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var req models.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Error:   "Please provide refresh token",
		})
		return
	}

	accessToken, err := h.authService.RefreshToken(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ResponseError{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.ResponseSuccess{
		Success: true,
		Message: "Token refreshed successfully",
		Data: gin.H{
			"access_token": accessToken,
		},
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refreshToken  body  models.LogoutRequest  true  "Refresh Token"
// @Success      200  {object}  response.ResponseSuccess
// @Failure      400  {object}  response.ResponseError
// @Failure      500  {object}  response.ResponseError
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	var req models.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Error:   "Please provide refresh token",
		})
		return
	}

	err := h.authService.Logout(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Error:   "Failed to logout",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.ResponseSuccess{
		Success: true,
		Message: "Logout successful",
	})
}
