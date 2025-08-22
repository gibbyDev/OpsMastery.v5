package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"OpsMastery.v5/internal/database"
	"OpsMastery.v5/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ListUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := database.DB().Unscoped().Find(&users).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	if err := database.DB().First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func DeleteUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	userIDParsed, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	result := database.DB().Unscoped().Delete(&models.User{}, userIDParsed)
	if result.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}

func UpdateUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	if err := database.DB().First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Only update fields if present in input
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Address != "" {
		user.Address = input.Address
	}
	if input.PhoneNumber != "" {
		user.PhoneNumber = input.PhoneNumber
	}
	if input.Role != "" {
		user.Role = input.Role
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		user.Password = input.Password // Consider hashing here!
	}
	// ...handle other fields as needed...

	if err := database.DB().Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
	}

	return c.JSON(user)
}

func SetUserRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB().First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	var input struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user.Role = input.Role
	if err := database.DB().Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user role"})
	}
	return c.JSON(user)
}

func GetCurrentUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found"})
	}

	var user models.User
	if err := database.DB().First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func UpdateCurrentUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found"})
	}

	var user models.User
	if err := database.DB().First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse form fields using Fiber's API
	name := c.FormValue("name")
	email := c.FormValue("email")
	// Handle file upload (if needed)
	fileHeader, err := c.FormFile("avatar")
	if err == nil {
		// Process the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to open uploaded file"})
		}
		defer file.Close()
		// ...handle file...
	}

	// Update user fields
	user.Name = name
	user.Email = email

	if err := database.DB().Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
	}
	return c.JSON(user)
}

func GetUserProfilePhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := database.DB().First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if len(user.ProfilePhoto) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Profile photo not found"})
	}
	// You may want to detect the image type; here we default to jpeg
	c.Set("Content-Type", "image/jpeg")
	return c.Send(user.ProfilePhoto)
}

func SearchUsers(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return c.JSON([]models.User{})
	}
	normalized := strings.ToLower(q)
	var users []models.User
	database.DB().Debug().Where(
		"LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(name) LIKE ?",
		"%"+normalized+"%", "%"+normalized+"%", "%"+normalized+"%",
	).Find(&users)
	return c.JSON(users)
}
