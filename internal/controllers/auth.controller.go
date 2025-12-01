package controllers

import (
	"backend-koda-shortlink/internal/models"
	"backend-koda-shortlink/internal/utils"
	"backend-koda-shortlink/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Register      godoc
// @Summary      Register new user
// @Description  Create a new user with a unique email
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        fullName  formData  string  true  "Full name user"
// @Param        email     formData  string  true  "Email user"
// @Param        password  formData  string  true  "Password user" format(password)
// @Success      201  {object}  lib.ResponseSuccess{data=models.Register}  "User created successfully."
// @Failure      400  {object}  lib.ResponseError  "Invalid request body or failed to hash password."
// @Failure      409  {object}  lib.ResponseError  "Email already registered."
// @Failure      500  {object}  lib.ResponseError  "Internal server error while creating user."
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
