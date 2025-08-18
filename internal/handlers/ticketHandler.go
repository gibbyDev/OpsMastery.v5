package handlers

import (
	"net/http"
	"strconv"

	"OpsMastery.v5/internal/database"
	"OpsMastery.v5/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateTicket(c *fiber.Ctx) error {
	var ticket models.Ticket
	if err := c.BodyParser(&ticket); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := database.DB().Create(&ticket).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(ticket)
}

func ListTickets(c *fiber.Ctx) error {
	var tickets []models.Ticket
	if err := database.DB().Find(&tickets).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tickets)
}

func GetTicketByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var ticket models.Ticket

	ticketID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	if err := database.DB().First(&ticket, ticketID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ticket)
}

func DeleteTicketByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ticketID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	result := database.DB().Unscoped().Delete(&models.Ticket{}, ticketID)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
	}
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Ticket deleted successfully"})
}

func UpdateTicketByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ticketID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	var ticket models.Ticket
	if err := c.BodyParser(&ticket); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := database.DB().Model(&models.Ticket{}).Where("id = ?", ticketID).Updates(ticket).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Ticket updated successfully"})
}
