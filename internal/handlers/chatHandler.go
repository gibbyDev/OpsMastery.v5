// // filepath: /home/cody/OpsMastery.v5/internal/handlers/chatHandler.go
package handlers

import (
	"fmt"
	"sync"
	"time"

	"OpsMastery.v5/internal/database"
	"OpsMastery.v5/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	clients     = make(map[*websocket.Conn]uint) // map connection to user ID
	broadcast   = make(chan Message)
	mu          sync.Mutex
	chatClients = make(map[uint]map[*websocket.Conn]bool) // chatID -> set of connections
)

type Message struct {
	SenderID    uint   `json:"sender_id"`
	RecipientID uint   `json:"recipient_id"`
	Content     string `json:"content"`
}

type Chat struct {
	ID        uint                 `json:"id"`
	Name      string               `json:"name"`
	TicketID  *uint                `json:"ticket_id,omitempty"`
	IsPrivate bool                 `json:"is_private"`
	Users     []models.User        `json:"users"`
	Messages  []models.ChatMessage `json:"messages"`
}

// WebSocket chat endpoint
func ChatWebSocket(c *fiber.Ctx) error {
	// Upgrade to WebSocket
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func HandleChat(c *websocket.Conn, claims map[string]interface{}) {
	senderID, _ := claims["userID"].(uint)

	// Register connection
	if chatClients[0] == nil {
		chatClients[0] = make(map[*websocket.Conn]bool)
	}
	chatClients[0][c] = true
	defer func() {
		delete(chatClients[0], c)
		c.Close()
	}()

	for {
		var msg struct {
			Content string `json:"content"`
		}
		if err := c.ReadJSON(&msg); err != nil {
			break
		}

		// Save message to DB
		chatMsg := models.ChatMessage{
			ChatID:   0,
			SenderID: senderID,
			Content:  msg.Content,
			SentAt:   time.Now(),
		}
		database.DB().Create(&chatMsg)

		// Broadcast to all clients in the chat
		for client := range chatClients[0] {
			client.WriteJSON(chatMsg)
		}
	}
}

// Broadcast messages to all clients
func StartBroadcast() {
	for {
		msg := <-broadcast
		mu.Lock()
		for conn := range clients {
			if err := conn.WriteJSON(msg); err != nil {
				fmt.Println("Error broadcasting:", err)
			}
		}
		mu.Unlock()
	}
}

// GET /api/v1/chats?user1=username1&user2=username2
func GetChatHistory(c *fiber.Ctx) error {
	user1 := c.Query("user1")
	user2 := c.Query("user2")
	var messages []models.ChatMessage
	database.DB().Where(
		"(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
		user1, user2, user2, user1,
	).Order("created_at asc").Find(&messages)
	return c.JSON(messages) // Returns [] if no messages exist
}

// GET /api/v1/chat_partners?user_id=some_user_id
func GetChatPartners(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	var partners []models.User
	database.DB().Raw(`
        SELECT DISTINCT u.*
        FROM users u
        JOIN chat_messages cm
          ON (cm.sender_id = u.id OR cm.recipient_id = u.id)
        WHERE (cm.sender_id = ? OR cm.recipient_id = ?) AND u.id != ?
    `, userID, userID, userID).Scan(&partners)
	return c.JSON(partners)
}

// DELETE /api/v1/chats/delete?user1=username1&user2=username2
func DeleteChatBetweenUsers(c *fiber.Ctx) error {
	user1 := c.Query("user1")
	user2 := c.Query("user2")
	result := database.DB().Where(
		"(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
		user1, user2, user2, user1,
	).Delete(&models.ChatMessage{})
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.JSON(fiber.Map{"deleted": result.RowsAffected})
}

// Create a new chat (group or direct)
func CreateChat(c *fiber.Ctx) error {
	var input struct {
		Name      string `json:"name"`
		UserIDs   []uint `json:"user_ids"`
		IsPrivate bool   `json:"is_private"`
		TicketID  *uint  `json:"ticket_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	chat := models.Chat{
		Name:      input.Name,
		IsPrivate: input.IsPrivate,
		TicketID:  input.TicketID,
	}

	if err := database.DB().Create(&chat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create chat"})
	}

	if len(input.UserIDs) > 0 {
		var users []models.User
		if err := database.DB().Where("id IN ?", input.UserIDs).Find(&users).Error; err == nil {
			database.DB().Model(&chat).Association("Users").Append(users)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(chat)
}

// Add users to an existing chat
func AddUsersToChat(c *fiber.Ctx) error {
	chatID := c.Params("chatId")
	var input struct {
		UserIDs []uint `json:"user_ids"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	var chat models.Chat
	if err := database.DB().First(&chat, chatID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Chat not found"})
	}
	var users []models.User
	if err := database.DB().Where("id IN ?", input.UserIDs).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Users not found"})
	}
	database.DB().Model(&chat).Association("Users").Append(users)
	return c.JSON(fiber.Map{"added": len(users)})
}

// Get all chats for a user (for sidebar)
func GetChatsForUser(c *fiber.Ctx) error {
	userID := c.Params("userId")
	var chats []models.Chat
	if err := database.DB().Joins("JOIN chat_users ON chat_users.chat_id = chats.id").
		Where("chat_users.user_id = ?", userID).
		Preload("Users").
		Preload("Messages").
		Find(&chats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch chats"})
	}
	return c.JSON(chats)
}

// Get all messages for a chat
func GetChatMessages(c *fiber.Ctx) error {
	chatID := c.Params("chatId")
	var messages []models.ChatMessage
	if err := database.DB().Where("chat_id = ?", chatID).Order("sent_at asc").Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch messages"})
	}
	return c.JSON(messages)
}

// Delete a chat (and its messages)
func DeleteChat(c *fiber.Ctx) error {
	chatID := c.Params("chatId")
	if err := database.DB().Where("chat_id = ?", chatID).Delete(&models.ChatMessage{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete messages"})
	}
	if err := database.DB().Delete(&models.Chat{}, chatID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete chat"})
	}
	return c.JSON(fiber.Map{"deleted": chatID})
}
