// // filepath: /home/cody/OpsMastery.v5/internal/handlers/chatHandler.go
package handlers

import (
	"fmt"
	"sync"

	"OpsMastery.v5/internal/database"
	"OpsMastery.v5/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	clients   = make(map[*websocket.Conn]uint) // map connection to user ID
	broadcast = make(chan Message)
	mu        sync.Mutex
)

type Message struct {
	SenderID    uint   `json:"sender_id"`
	RecipientID uint   `json:"recipient_id"`
	Content     string `json:"content"`
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

	mu.Lock()
	clients[c] = 0 // Not used anymore, but kept for compatibility
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, c)
		mu.Unlock()
		c.Close()
	}()

	for {
		var msg Message
		if err := c.ReadJSON(&msg); err != nil {
			break
		}
		msg.SenderID = senderID

		// Persist to DB
		chatMsg := models.ChatMessage{
			SenderID:    msg.SenderID,
			RecipientID: msg.RecipientID,
			Content:     msg.Content,
		}
		database.DB().Create(&chatMsg)

		broadcast <- msg
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
