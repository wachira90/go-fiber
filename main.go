package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("hello world")

	// fiber instance
	app := fiber.New()

	// routes
	//  app.Get("/", func(c *fiber.Ctx) error {
	//   return c.SendString("hello world 🌈")
	//  })

	app.Get("/", func(c *fiber.Ctx) error { // STRING
		return c.SendString("hello world 🌈")
	})

	app.Get("/info", func(c *fiber.Ctx) error { // JSON
		return c.JSON(fiber.Map{
			"msg":     "hello world 🚀",
			"go":      "fiber 🥦",
			"boolean": true,
			"number":  1234,
		})
	})

	// app listening at PORT: 3000
	app.Listen(":3000")
}
