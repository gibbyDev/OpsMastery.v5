package chat_service

import (
	"fmt"

	"OpsMastery.v5/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Message struct definition
type Message struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
	// Add more fields as needed
}

// Place this at the package level:
var clients = make(map[*websocket.Conn]bool)

// Place this at the package level:
func broadcastMessage(msg Message) {
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func Start() {
	app := fiber.New()

	// Static files
	app.Static("/public", "./public")

	app.Use("/ws", func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
		}
		claims, err := utils.ValidateJWT(token, false)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}
		c.Locals("userEmail", claims["email"])
		return c.Next()
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		userEmail, _ := c.Locals("userEmail").(string)
		clients[c] = true
		defer func() {
			delete(clients, c)
			c.Close()
		}()
		for {
			var msg Message
			err := c.ReadJSON(&msg)
			if err != nil {
				_, rawMsg, readErr := c.ReadMessage()
				if readErr != nil {
					fmt.Println("WebSocket closed:", readErr)
					break
				}
				msg = Message{Content: string(rawMsg), Sender: userEmail}
			} else {
				msg.Sender = userEmail
			}
			broadcastMessage(msg)
		}
	}))

	// Start the server
	app.Listen(":5000")
}
