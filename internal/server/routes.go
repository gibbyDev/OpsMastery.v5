package server

import (
	"OpsMastery.v5/internal/handlers"
	"OpsMastery.v5/internal/middleware"
	"OpsMastery.v5/internal/webrtc_service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (s *FiberServer) RegisterFiberRoutes(db *gorm.DB) {
	api := s.Group("/api/v1")

	// Public routes (no authentication required)
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running!")
	})
	// Authentication routes
	api.Post("/signup", handlers.SignUp)
	api.Post("/signin", handlers.SignIn)
	api.Get("/verify/:token", handlers.VerifyEmail)
	api.Post("/forgot-password", handlers.RequestPasswordReset)
	api.Post("/reset-password", handlers.ResetPassword)

	// Public route to get a user's profile photo by ID
	api.Get("/users/:id/profile_photo", handlers.GetUserProfilePhoto)

	// OAuth routes
	s.Get("/auth/:provider", handlers.OAuthLogin)
	s.Get("/auth/:provider/callback", handlers.OAuthCallback)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.JWTMiddleware)

	// Allow all authenticated users to list users
	protected.Get("/users", handlers.ListUsers)
	protected.Get("/users/search", handlers.SearchUsers)

	// Admin routes
	protected.Delete("/users/:id", middleware.OnlyAdmin(db, handlers.DeleteUserByID))
	protected.Put("/users/:id/role", middleware.OnlyAdmin(db, handlers.SetUserRole))

	// Moderator routes
	protected.Get("/users/:id", middleware.OnlyModerator(db, handlers.GetUserByID))
	protected.Put("/users/:id", middleware.OnlyModerator(db, handlers.UpdateUserByID))

	// User routes
	protected.Get("/users/me", middleware.OnlyUser(db, handlers.GetCurrentUser))
	protected.Put("/users/me", middleware.OnlyUser(db, handlers.UpdateCurrentUser))

	// Other protected routes
	protected.Post("/signout", handlers.SignOut)
	protected.Post("/auth/refresh", handlers.RefreshToken)

	// Protected routes for tickets
	protected.Post("/ticket", handlers.CreateTicket)
	protected.Get("/tickets", handlers.ListTickets)
	protected.Get("/ticket/:id", handlers.GetTicketByID)
	protected.Put("/ticket/:id", handlers.UpdateTicketByID)
	protected.Delete("/ticket/:id", handlers.DeleteTicketByID)

	// Protected routes for chat
	protected.Get("/chats", handlers.GetChatHistory)
	protected.Post("/chats", handlers.CreateChat)
	protected.Post("/chats/:chatId/users", handlers.AddUsersToChat)
	protected.Get("/chats/user/:userId", handlers.GetChatsForUser)
	protected.Get("/chats/:chatId/messages", handlers.GetChatMessages)
	protected.Delete("/chats/:chatId", handlers.DeleteChat)

	// WebRTC routes
	protected.Post("/webrtc/start", func(c *fiber.Ctx) error {
		go webrtc_service.StartSignalingServer()
		return c.JSON(fiber.Map{"status": "signaling server started"})
	})

	// Chat WebSocket route
	// protected.Get("/chat/ws", websocket.New(func(c *websocket.Conn) {
	// 	token := c.Query("token")
	// 	if token == "" {
	// 		fmt.Println("WebSocket closed: missing token")
	// 		c.Close()
	// 		return
	// 	}
	// 	claims, err := utils.ValidateJWT(token, false)
	// 	if err != nil {
	// 		fmt.Println("WebSocket closed: invalid token:", err)
	// 		c.Close()
	// 		return
	// 	}
	// 	// Set user info as locals for use in HandleChat

	// 	handlers.HandleChat(c, claims)
	// }))

	api.Options("/chat/ws", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		return c.SendStatus(fiber.StatusNoContent)
	})
}
