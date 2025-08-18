package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"OpsMastery.v5/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var (
	AccessTokenSecret  = []byte(os.Getenv("JWT_ACCESS_SECRET"))
	RefreshTokenSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
)

func ValidateAccessToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return AccessTokenSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	var token string

	if authHeader != "" {
		// Extract the token from the "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			token = tokenParts[1]
		}
	}

	// If no token in header, check cookie
	if token == "" {
		token = c.Cookies("access_token")
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	claims, err := utils.ValidateJWT(token, false)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	c.Locals("userID", claims["sub"])
	c.Locals("userRole", claims["role"])
	c.Locals("userEmail", claims["email"])

	return c.Next()
}
