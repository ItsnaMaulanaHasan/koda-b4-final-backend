package controllers

import (
	"backend-koda-shortlink/internal/models"
	"backend-koda-shortlink/internal/utils"
	"backend-koda-shortlink/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/matthewhartstonge/argon2"
)

// Login         godoc
// @Summary      Login user
// @Description  Log in with existing email data
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        email     formData  string  true  "User email"
// @Param        password  formData  string  true  "User password" format(password)
// @Success      200  {object}  response.ResponseSuccess{data=object{token=string}}  "Login successful"
// @Failure      400  {object}  response.ResponseError  "Invalid request body"
// @Failure      401  {object}  response.ResponseError  "Wrong email or password"
// @Failure      500  {object}  response.ResponseError  "Internal server error"
// @Router       /auth/login [post]
func Login(ctx *gin.Context) {
	var bodyLogin models.Login
	err := ctx.ShouldBindWith(&bodyLogin, binding.Form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Message: "Please provide valid email and password",
			Error:   err.Error(),
		})
		return
	}

	// check input email
	user, message, err := models.GetUserByEmail(&bodyLogin)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if message == "Wrong email or password" {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, response.ResponseError{
			Success: false,
			Message: message,
			Error:   err.Error(),
		})
		return
	}

	// check input password
	isPasswordValid, err := argon2.VerifyEncoded(
		[]byte(bodyLogin.Password),
		[]byte(user.Password),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Message: "Failed to verify password. Please try again",
			Error:   err.Error(),
		})
		return
	}

	if !isPasswordValid {
		ctx.JSON(http.StatusUnauthorized, response.ResponseError{
			Success: false,
			Message: "Wrong email or password",
		})
		return
	}

	// generate refresh token
	refreshToken, expiresAt, err := utils.GenerateRefreshToken(user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Message: "Failed to generate refresh token. Please try again",
			Error:   err.Error(),
		})
		return
	}

	// save refresh token to db
	ipAddress := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()
	sessionId, err := models.CreateSession(user.Id, refreshToken, expiresAt, ipAddress, userAgent)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Message: "Failed to create session. Please try again",
			Error:   err.Error(),
		})
		return
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(user.Id, sessionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Message: "Failed to generate access token. Please try again",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.ResponseSuccess{
		Success: true,
		Message: "Login successful!",
		Data: gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// Register      godoc
// @Summary      Register new user
// @Description  Create a new user with a unique email
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        fullname  formData  string  true  "Full name user"
// @Param        email     formData  string  true  "Email user"
// @Param        password  formData  string  true  "Password user" format(password)
// @Success      201  {object}  response.ResponseSuccess{data=models.Register}  "User created successfully."
// @Failure      400  {object}  response.ResponseError  "Invalid request body or failed to hash password."
// @Failure      409  {object}  response.ResponseError  "Email already registered."
// @Failure      500  {object}  response.ResponseError  "Internal server error while creating user."
// @Router       /auth/register [post]
func Register(ctx *gin.Context) {
	var bodyRegister models.Register
	err := ctx.ShouldBindWith(&bodyRegister, binding.Form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Message: "Please provide valid registration information",
			Error:   err.Error(),
		})
		return
	}

	// check user email
	exists, err := models.CheckUserEmail(bodyRegister.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Message: "Unable to process registration. Please try again",
			Error:   err.Error(),
		})
		return
	}

	if exists {
		ctx.JSON(http.StatusConflict, response.ResponseError{
			Success: false,
			Message: "Email is already registered",
		})
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(bodyRegister.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
		return
	}
	bodyRegister.Password = string(hashedPassword)

	isSuccess, message, err := models.RegisterUser(&bodyRegister)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: isSuccess,
			Message: message,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.ResponseSuccess{
		Success: isSuccess,
		Message: message,
		Data: models.Register{
			Id:       bodyRegister.Id,
			FullName: bodyRegister.FullName,
			Email:    bodyRegister.Email,
		},
	})
}

// RefreshToken  godoc
// @Summary      Refresh access token
// @Description  Get new access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  body  object{refresh_token=string}  true  "Refresh Token"
// @Success      200  {object}  response.ResponseSuccess{data=object{access_token=string,expires_in=int}}  "Token refreshed"
// @Failure      400  {object}  response.ResponseError  "Invalid request"
// @Failure      401  {object}  response.ResponseError  "Invalid or expired refresh token"
// @Failure      500  {object}  response.ResponseError  "Internal server error"
// @Router       /auth/refresh [post]
func RefreshToken(ctx *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Message: "Please provide refresh token",
			Error:   err.Error(),
		})
		return
	}

	claims, err := utils.VerifyRefreshToken(body.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ResponseError{
			Success: false,
			Message: "Invalid or expired refresh token",
			Error:   err.Error(),
		})
		return
	}

	// check refresh token on db
	session, err := models.GetSessionByRefreshToken(body.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ResponseError{
			Success: false,
			Message: "Invalid or expired refresh token",
			Error:   err.Error(),
		})
		return
	}

	// generate new access token
	accessToken, err := utils.GenerateAccessToken(claims.Id, session.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Message: "Failed to generate access token",
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

// Logout        godoc
// @Summary      Logout user
// @Description  Invalidate refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  body  object{refresh_token=string}  true  "Refresh Token"
// @Success      200  {object}  response.ResponseSuccess  "Logout successful"
// @Failure      400  {object}  response.ResponseError  "Invalid request"
// @Failure      500  {object}  response.ResponseError  "Internal server error"
// @Router       /auth/logout [post]
func Logout(ctx *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ResponseError{
			Success: false,
			Message: "Please provide refresh token",
			Error:   err.Error(),
		})
		return
	}

	// Invalidate session
	err := models.InvalidateSession(1, body.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ResponseError{
			Success: false,
			Message: "Failed to logout",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.ResponseSuccess{
		Success: true,
		Message: "Logout successful",
	})
}
