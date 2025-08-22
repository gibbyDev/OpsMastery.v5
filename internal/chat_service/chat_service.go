package chat_service

import (
	"fmt"

	"OpsMastery.v5/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Message struct definition
// This struct represents a chat message sent/received via WebSocket.
// It contains the message content and the sender's identifier (email).
type Message struct {
	Content string `json:"content"` // The actual message text
	Sender  string `json:"sender"`  // The sender's email or username
	// Add more fields as needed (e.g., timestamp, chatId, etc.)
}

// clients is a map of active WebSocket connections.
// The key is the pointer to the websocket.Conn, the value is a boolean (always true).
// This allows us to keep track of all connected clients for broadcasting messages.
var clients = make(map[*websocket.Conn]bool)

// broadcastMessage sends a given Message to all connected clients.
// If a client connection fails, it is removed from the clients map.
func broadcastMessage(msg Message) {
	for client := range clients {
		// Send the message as JSON to the client
		err := client.WriteJSON(msg)
		if err != nil {
			// If sending fails, print error, close connection, and remove client
			fmt.Println("Error sending message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// Start initializes and runs the chat WebSocket server.
func Start() {
	// Create a new Fiber app instance
	app := fiber.New()

	// Serve static files from the ./public directory at /public
	app.Static("/public", "./public")

	// Middleware for /ws route to validate JWT token before allowing WebSocket upgrade
	app.Use("/ws", func(c *fiber.Ctx) error {
		// Extract token from query parameters
		token := c.Query("token")
		if token == "" {
			// If no token, reject the request
			return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
		}
		// Validate the JWT token
		claims, err := utils.ValidateJWT(token, false)
		if err != nil {
			// If token is invalid, reject the request
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}
		// Store the user's email in Fiber's context locals for later use
		c.Locals("userEmail", claims["email"])
		// Proceed to the next handler (WebSocket upgrade)
		return c.Next()
	})

	// WebSocket endpoint for chat communication
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Retrieve the user's email from context locals
		userEmail, _ := c.Locals("userEmail").(string)
		// Register the new client connection
		clients[c] = true
		// Ensure client is removed and connection closed when handler exits
		defer func() {
			delete(clients, c)
			c.Close()
		}()
		// Main loop: read messages from this client and broadcast to all clients
		for {
			var msg Message
			// Try to read a JSON message from the client
			err := c.ReadJSON(&msg)
			if err != nil {
				// If JSON parsing fails, try to read a raw message (e.g., plain text)
				_, rawMsg, readErr := c.ReadMessage()
				if readErr != nil {
					// If reading fails, log and break the loop (disconnect client)
					fmt.Println("WebSocket closed:", readErr)
					break
				}
				// Construct a Message from the raw message and set sender
				msg = Message{Content: string(rawMsg), Sender: userEmail}
			} else {
				// If JSON was parsed, set the sender field
				msg.Sender = userEmail
			}
			// Broadcast the received message to all connected clients
			broadcastMessage(msg)
		}
	}))

	// Start the Fiber server on port 5000
	app.Listen(":5000")
}
