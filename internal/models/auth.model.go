package models

type User struct {
	Id           int     `json:"id" db:"id"`
	ProfilePhoto *string `json:"profilePhoto" db:"profile_photo"`
	FullName     string  `json:"fullName" db:"fullname"`
	Email        string  `json:"email" db:"email"`
	Password     string  `json:"-" db:"-"`
}

type RegisterRequest struct {
	FullName string `form:"fullname" json:"fullName" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
