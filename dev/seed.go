// File: dev/seed.go
package dev

import (
	"fmt"
	"math/rand"
	"time"

	"OpsMastery.v5/internal/models"
	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

// Helper to find a user by ID from a slice
func findUserByID(users []models.User, id uint) models.User {
	for _, u := range users {
		if u.ID == id {
			return u
		}
	}
	return models.User{}
}

// Helper to get n random users from a slice
func randomUsers(users []models.User, n int) []models.User {
	if len(users) < n {
		return users
	}
	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	return users[:n]
}

func Seed(db *gorm.DB) {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Seeding database...")

	roles := []string{"User", "Admin", "Manager"}
	statuses := []string{"Open", "In Progress", "Resolved", "Closed"}
	priorities := []string{"Low", "Normal", "High", "Urgent"}

	// Seed Clients
	var clients []models.Client
	for i := 0; i < 5; i++ {
		client := models.Client{
			Name:  faker.Name(),
			Email: faker.Email(),
			Phone: faker.Phonenumber(),
		}
		db.Create(&client)
		clients = append(clients, client)
	}

	// Seed Users
	var users []models.User
	for i := 0; i < 20; i++ {
		user := models.User{
			Name:     faker.Name(),
			Email:    faker.Email(),
			Password: "hashedpassword123", // Replace with actual hash in production
			Role:     roles[rand.Intn(len(roles))],
			Active:   rand.Intn(2) == 1,
			Username: faker.Username(),
			ClientID: &clients[rand.Intn(len(clients))].ID,
		}
		db.Create(&user)
		users = append(users, user)
	}

	// Seed Tickets
	var tickets []models.Ticket
	for i := 0; i < 10; i++ {
		reporter := users[rand.Intn(len(users))]
		assignee := users[rand.Intn(len(users))]
		client := clients[rand.Intn(len(clients))]
		ticket := models.Ticket{
			Title:       faker.Sentence(),
			Description: faker.Paragraph(),
			ReporterID:  reporter.ID,
			AssigneeID:  assignee.ID,
			ClientID:    client.ID,
			Status:      statuses[rand.Intn(len(statuses))],
			Priority:    priorities[rand.Intn(len(priorities))],
		}
		db.Create(&ticket)
		tickets = append(tickets, ticket)
	}

	// Seed General/Group Chats
	for i := 0; i < 5; i++ {
		chatUsers := randomUsers(users, rand.Intn(5)+3) // 3-7 users
		chat := models.Chat{
			Name:      faker.Word() + " Group",
			IsPrivate: false,
			Users:     chatUsers,
		}
		db.Create(&chat)

		// Seed messages for group chat
		for j := 0; j < rand.Intn(10)+5; j++ {
			sender := chatUsers[rand.Intn(len(chatUsers))]
			msg := models.ChatMessage{
				ChatID:   chat.ID,
				SenderID: sender.ID,
				Content:  faker.Sentence(),
				SentAt:   time.Now().Add(time.Duration(-rand.Intn(1000)) * time.Minute),
			}
			db.Create(&msg)
		}
	}

	// Seed Ticket Chats
	for _, ticket := range tickets {
		// Reporter, Assignee, and 1-2 random watchers
		watchers := randomUsers(users, rand.Intn(2)+1)
		chatUsers := []models.User{
			findUserByID(users, ticket.ReporterID),
			findUserByID(users, ticket.AssigneeID),
		}
		chatUsers = append(chatUsers, watchers...)

		chat := models.Chat{
			Name:      "Ticket Chat",
			TicketID:  &ticket.ID,
			IsPrivate: false,
			Users:     chatUsers,
		}
		db.Create(&chat)

		// Seed messages for ticket chat
		for j := 0; j < rand.Intn(10)+5; j++ {
			sender := chatUsers[rand.Intn(len(chatUsers))]
			msg := models.ChatMessage{
				ChatID:   chat.ID,
				SenderID: sender.ID,
				Content:  faker.Sentence(),
				SentAt:   time.Now().Add(time.Duration(-rand.Intn(1000)) * time.Minute),
			}
			db.Create(&msg)
		}
	}

	// Seed Private 1:1 Chats
	for i := 0; i < 10; i++ {
		usersPair := randomUsers(users, 2)
		chat := models.Chat{
			Name:      usersPair[0].Name + " & " + usersPair[1].Name,
			IsPrivate: true,
			Users:     usersPair,
		}
		db.Create(&chat)

		// Seed messages for private chat
		for j := 0; j < rand.Intn(10)+5; j++ {
			sender := usersPair[rand.Intn(2)]
			msg := models.ChatMessage{
				ChatID:   chat.ID,
				SenderID: sender.ID,
				Content:  faker.Sentence(),
				SentAt:   time.Now().Add(time.Duration(-rand.Intn(1000)) * time.Minute),
			}
			db.Create(&msg)
		}
	}

	fmt.Println("Seeding complete!")
}
