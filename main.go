package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Book struct represents a book model
type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

var DB *gorm.DB

func main() {
	// Connect to PostgreSQL database
	dsn := "host=192.168.1.15 user=postgres password=postgres dbname=testgo port=5432 TimeZone=Asia/Bangkok sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	DB = db

	// Migrate the Book model
	DB.AutoMigrate(&Book{})

	// Create a new Fiber instance
	app := fiber.New()

	// API routes
	app.Get("/books", getAllBooks)
	app.Post("/books", createBook)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

// getAllBooks handler retrieves all books from the database
func getAllBooks(c *fiber.Ctx) error {
	var books []Book
	DB.Find(&books)
	return c.JSON(books)
}

// createBook handler creates a new book in the database
func createBook(c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := DB.Create(book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(book)
}
