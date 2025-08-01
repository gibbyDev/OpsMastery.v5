package server

import (
	"github.com/gofiber/fiber/v2"

	"OpsMastery.v5/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "OpsMastery.v5",
			AppName:      "OpsMastery.v5",
		}),

		db: database.New(),
	}

	return server
}
