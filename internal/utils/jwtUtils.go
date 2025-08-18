package utils

import (
	"log"
	"os"
	"time"

	"OpsMastery.v5/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

var (
	accessTokenSecret  = []byte(os.Getenv("JWT_ACCESS_SECRET"))
	refreshTokenSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
)

func GenerateJWT(user models.User) (string, string, error) {
	accessClaims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString(accessTokenSecret)
	if err != nil {
		log.Println("Error generating access token:", err)
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(refreshTokenSecret)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func ValidateJWT(tokenString string, isRefreshToken bool) (jwt.MapClaims, error) {
	secret := accessTokenSecret
	if isRefreshToken {
		secret = refreshTokenSecret
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorClaimsInvalid)
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
