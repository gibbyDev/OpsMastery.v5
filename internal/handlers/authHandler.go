package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"OpsMastery.v5/internal/database"
	"OpsMastery.v5/internal/models"
	"OpsMastery.v5/internal/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	var user models.User
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user.Email = input.Email
	user.Password = input.Password
	user.Name = input.Name
	user.Role = input.Role
	user.Active = false
	user.VerificationToken = utils.GenerateRandomToken()

	storedHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = string(storedHash)

	log.Printf("Hashed password for user %s: %s\n", user.Email, user.Password)

	if err := database.DB().Create(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := utils.SendVerificationEmail(user.Email, user.VerificationToken); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send verification email"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

func SignIn(c *fiber.Ctx) error {
	var userInput struct {
		EmailOrUsername string `json:"emailOrUsername"`
		Password        string `json:"password"`
	}

	if err := c.BodyParser(&userInput); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := database.DB().Where("email = ? OR username = ?", userInput.EmailOrUsername, userInput.EmailOrUsername).First(&user).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email/username or password"})
	}

	if !user.Active {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Please verify your email before signing in",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email/username or password"})
	}

	// Generate JWTs
	accessToken, refreshToken, err := utils.GenerateJWT(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate tokens"})
	}

	// Set Refresh Token in HttpOnly Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/api/v1/auth/refresh", // restrict usage
	})

	// Return Access Token in response (frontend stores in memory)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":      "Sign in successful",
		"access_token": accessToken,
		"user": fiber.Map{
			"id":          user.ID,
			"email":       user.Email,
			"name":        user.Name,
			"role":        user.Role,
			"username":    user.Username,
			"address":     user.Address,
			"phoneNumber": user.PhoneNumber,
			"active":      user.Active,
		},
	})
}

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token not found"})
	}

	claims, err := utils.ValidateJWT(refreshToken, true)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	// Get user from claims
	var user models.User
	if err := database.DB().First(&user, claims["sub"]).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := utils.GenerateJWT(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate tokens"})
	}

	// Rotate refresh token
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/api/auth/refresh",
	})

	// Return new Access Token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func SignOut(c *fiber.Ctx) error {
	// Invalidate the refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/api/auth/refresh",
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully signed out"})
}

// func SignIn(c *fiber.Ctx) error {
// 	var userInput struct {
// 		EmailOrUsername string `json:"email"` // Accepts either email or username
// 		Password        string `json:"password"`
// 	}

// 	if err := c.BodyParser(&userInput); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
// 	}

// 	var user models.User
// 	if err := database.DB().Where("email = ? OR username = ?", userInput.EmailOrUsername, userInput.EmailOrUsername).First(&user).Error; err != nil {
// 		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email/username or password"})
// 	}

// 	if !user.Active {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Please verify your email before signing in",
// 		})
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
// 		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email/username or password"})
// 	}

// 	accessToken, refreshToken, err := utils.GenerateJWT(user)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate tokens"})
// 	}

// 	log.Printf("Setting access token: %s", accessToken)
// 	log.Printf("Setting refresh token: %s", refreshToken)

// 	c.Cookie(&fiber.Cookie{
// 		Name:     "access_token",
// 		Value:    accessToken,
// 		Expires:  time.Now().Add(15 * time.Minute),
// 		HTTPOnly: true, // <-- set to true
// 		Secure:   true,
// 		SameSite: "strict",
// 		Domain:   "",
// 		Path:     "/",
// 	})

// 	c.Cookie(&fiber.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		Expires:  time.Now().Add(7 * 24 * time.Hour),
// 		HTTPOnly: true, // <-- set to true
// 		Secure:   true,
// 		SameSite: "strict",
// 		Domain:   "",
// 		Path:     "/",
// 	})

// 	// Return user data (excluding sensitive fields)
// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"message":      "Sign in successful",
// 		"user_id":      user.ID,
// 		"email":        user.Email,
// 		"name":         user.Name,
// 		"role":         user.Role,
// 		"username":     user.Username,
// 		"address":      user.Address,
// 		"phone_number": user.PhoneNumber,
// 		"active":       user.Active,
// 	})
// }

// func SignOut(c *fiber.Ctx) error {
// 	c.Cookie(&fiber.Cookie{
// 		Name:     "jwt",
// 		Value:    "",
// 		Expires:  time.Now().Add(-1 * time.Hour),
// 		HTTPOnly: true,
// 	})
// 	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully signed out"})
// }

// func RefreshToken(c *fiber.Ctx) error {
// 	refreshToken := c.Cookies("refresh_token")
// 	if refreshToken == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Refresh token not found",
// 		})
// 	}

// 	claims, err := utils.ValidateJWT(refreshToken, true)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Invalid refresh token",
// 		})
// 	}

// 	// Get user from claims
// 	var user models.User
// 	if err := database.DB().First(&user, claims["sub"]).Error; err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "User not found",
// 		})
// 	}

// 	// Generate new tokens
// 	accessToken, refreshToken, err := utils.GenerateJWT(user)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Could not generate tokens",
// 		})
// 	}

// 	// Set new cookies
// 	c.Cookie(&fiber.Cookie{
// 		Name:     "access_token",
// 		Value:    accessToken,
// 		Expires:  time.Now().Add(15 * time.Minute),
// 		HTTPOnly: true,
// 	})

// 	c.Cookie(&fiber.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		Expires:  time.Now().Add(7 * 24 * time.Hour),
// 		HTTPOnly: true,
// 	})

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Tokens refreshed successfully",
// 	})
// }

// func VerifyEmail(c *fiber.Ctx) error {
// 	token := c.Params("token")

// 	var user models.User
// 	if err := database.DB().Where("verification_token = ?", token).First(&user).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": "Invalid verification token",
// 		})
// 	}

// 	user.Active = true
// 	user.VerificationToken = "" // Clear the token after verification

// 	if err := database.DB().Save(&user).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to verify email",
// 		})
// 	}

//		return c.Status(fiber.StatusOK).JSON(fiber.Map{
//			"message": "Email verified successfully",
//		})
//	}
func VerifyEmail(c *fiber.Ctx) error {
	var body struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&body); err != nil || body.Token == "" {
		return c.Redirect("http://localhost:3000/auth/verification-failed", fiber.StatusSeeOther)
	}

	var user models.User
	if err := database.DB().Where("verification_token = ?", body.Token).First(&user).Error; err != nil {
		return c.Redirect("http://localhost:3000/auth/verification-failed", fiber.StatusSeeOther)
	}

	user.Active = true
	user.VerificationToken = ""

	if err := database.DB().Save(&user).Error; err != nil {
		return c.Redirect("http://localhost:3000/auth/verification-failed", fiber.StatusSeeOther)
	}

	return c.Redirect("http://localhost:3000/auth/verification-success", fiber.StatusSeeOther)
}

func RequestPasswordReset(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := database.DB().Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Don't reveal if email exists or not
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If your email is registered, you will receive a password reset link",
		})
	}

	resetToken := utils.GenerateRandomToken()
	user.ResetToken = resetToken
	user.ResetTokenExpiry = time.Now().Add(1 * time.Hour)

	// Log the token being set
	log.Printf("Setting reset token for user %s: %s", user.Email, resetToken)

	if err := database.DB().Save(&user).Error; err != nil {
		log.Printf("Error saving user with reset token: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process password reset",
		})
	}

	// Add this log right before SendPasswordResetEmail
	log.Printf("About to send reset email with token: %s", resetToken)

	if err := utils.SendPasswordResetEmail(user.Email, resetToken); err != nil {
		log.Printf("Error sending reset email: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send reset email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If your email is registered, you will receive a password reset link",
	})
}

func ResetPassword(c *fiber.Ctx) error {
	var body struct {
		ResetToken  string `json:"reset_token"`
		NewPassword string `json:"new_password"`
	}

	// Log the raw request body
	rawBody := string(c.Body())
	log.Printf("Received reset password request body: %s", rawBody)

	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("Reset token: %s", body.ResetToken)

	// Validate the token
	var user models.User
	if err := database.DB().Where("reset_token = ?", body.ResetToken).First(&user).Error; err != nil {
		log.Printf("No user found with reset token: %s, error: %v", body.ResetToken, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reset token",
		})
	}

	// Check if token has expired
	if user.ResetTokenExpiry.Before(time.Now()) {
		log.Printf("Token expired. Expiry: %v, Current time: %v", user.ResetTokenExpiry, time.Now())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Reset token has expired",
		})
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Update user's password and clear reset token
	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}

	if err := database.DB().Save(&user).Error; err != nil {
		log.Printf("Error saving user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	log.Printf("Password successfully reset for user: %s", user.Email)
	return c.JSON(fiber.Map{
		"message": "Password successfully reset",
	})
}

func validateResetToken(token string) (bool, error) {
	var user models.User
	if err := database.DB().Where("reset_token = ?", token).First(&user).Error; err != nil {
		log.Printf("No user found with reset token: %s, error: %v", token, err)
		return false, err
	}

	// Check if token has expired
	if user.ResetTokenExpiry.Before(time.Now()) {
		log.Printf("Token expired. Expiry: %v, Current time: %v", user.ResetTokenExpiry, time.Now())
		return false, fmt.Errorf("reset token has expired")
	}

	log.Printf("Token validated successfully for user: %s", user.Email)
	return true, nil
}
