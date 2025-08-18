package webrtc_service

import (
	"fmt"

	"OpsMastery.v5/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/pion/webrtc/v3"
)

func StartSignalingServer() {
	// Create a new peer connection configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new peer connection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		fmt.Println("Failed to create peer connection:", err)
		return
	}

	defer peerConnection.Close()

	app := fiber.New()

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
		fmt.Println("WebRTC signaling WebSocket connected for user:", userEmail)
		// TODO: Implement signaling logic here
		// You can use userEmail for user identification
	}))

	app.Listen(":4000")

	// Set up signaling handlers (to be implemented)
	fmt.Println("Signaling server started")
}
