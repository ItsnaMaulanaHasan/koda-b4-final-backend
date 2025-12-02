package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserPayload struct {
	Id        int `json:"id"`
	SessionId int `json:"sessionId"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId, sessionId int) (string, error) {
	secretKey := []byte(os.Getenv("APP_SECRET"))
	claims := UserPayload{
		Id:        userId,
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshToken(userId int) (string, time.Time, error) {
	secretKey := []byte(os.Getenv("REFRESH_SECRET"))
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	claims := UserPayload{
		Id: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	return signedToken, expiresAt, err
}

func VerifyAccessToken(tokenString string) (*UserPayload, error) {
	secretKey := []byte(os.Getenv("APP_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &UserPayload{}, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func VerifyRefreshToken(tokenString string) (*UserPayload, error) {
	secretKey := []byte(os.Getenv("REFRESH_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &UserPayload{}, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
