package main

import (
	"log"

	"example.com/e/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/receipts/process", handlers.ProcessReceipt)
	app.Get("/receipts/:id/points", handlers.GetReceiptPoints)

	log.Println("Server started on :8080")
	log.Fatal(app.Listen(":8080"))
}
