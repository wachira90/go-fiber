package main

import (
	"fmt"
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
	dsn := "host=localhost user=postgres password=your_password dbname=your_database_name port=5432 sslmode=disable"
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
	app.Get("/books/:id", getBookByID)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

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

// getBookByID handler retrieves a book by its ID from the database
func getBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("Book with ID %s not found", id),
		})
	}
	return c.JSON(book)
}

// updateBook handler updates a book in the database
func updateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("Book with ID %s not found", id),
		})
	}

	updatedBook := new(Book)
	if err := c.BodyParser(updatedBook); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	book.Title = updatedBook.Title
	book.Author = updatedBook.Author
	book.Rating = updatedBook.Rating

	if err := DB.Save(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(book)
}

// deleteBook handler deletes a book from the database
func deleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("Book with ID %s not found", id),
		})
	}

	if err := DB.Delete(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Book with ID %s deleted successfully", id),
	})
}
