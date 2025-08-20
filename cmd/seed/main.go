package main

import (
	"OpsMastery.v5/dev"
	"OpsMastery.v5/internal/database"
)

func main() {
	database.Init()
	dev.Seed(database.DB())
}
