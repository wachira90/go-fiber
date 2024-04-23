# golang fiber


## run 

```
go run main.go
```
## Test data

- Create a new book:
  ```
  POST http://localhost:3000/books
  {
    "title": "The Catcher in the Rye",
    "author": "J.D. Salinger",
    "rating": 4
  }
  ```

- Retrieve a book by ID:
  ```
  GET http://localhost:3000/books/1
  ```

- Update a book:
  ```
  PUT http://localhost:3000/books/1
  {
    "title": "The Catcher in the Rye (Updated)",
    "author": "J.D. Salinger",
    "rating": 5
  }
  ```

- Delete a book:
  ```
  DELETE http://localhost:3000/books/1
  ```


## build 

```
go build -o main.exe main.go

```

## authen with jsonwebtoken

```go
package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User struct represents a user model
type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

// Book struct represents a book model
type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

var DB *gorm.DB
var secretKey = "your_secret_key"

func main() {
	// Connect to PostgreSQL database
	dsn := "host=localhost user=postgres password=your_password dbname=your_database_name port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	DB = db

	// Migrate the User and Book models
	DB.AutoMigrate(&User{}, &Book{})

	// Create a new Fiber instance
	app := fiber.New()

	// Authentication routes
	app.Post("/register", registerUser)
	app.Post("/login", loginUser)

	// Authenticated routes
	app.Use(authMiddleware)
	app.Get("/books", getAllBooks)
	app.Post("/books", createBook)
	app.Get("/books/:id", getBookByID)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

// registerUser handler registers a new user
func registerUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	user.Password = string(hashedPassword)

	if err := DB.Create(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(user)
}

// loginUser handler authenticates a user and generates a JWT token
func loginUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	var existingUser User
	if err := DB.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid username or password",
		})
	}

	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid username or password",
		})
	}

	// Generate a JWT token
	token, err := generateJWTToken(strconv.Itoa(int(existingUser.ID)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

// authMiddleware middleware checks if the user is authenticated
func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Missing Authorization header",
		})
	}

	token, err := validateJWTToken(authHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userID, err := strconv.Atoi(token.Claims.(jwt.MapClaims)["sub"].(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	c.Locals("userID", uint(userID))
	return c.Next()
}

// generateJWTToken generates a new JWT token
func generateJWTToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// validateJWTToken validates a JWT token
func validateJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}

// Existing getAllBooks, createBook, getBookByID, updateBook, and deleteBook handlers
// ...
```


In this example, we've added the following new features:

1. **User Registration**: The `registerUser` handler registers a new user by hashing the password using `bcrypt` and storing the user information in the database.

2. **User Login**: The `loginUser` handler authenticates a user by checking the provided username and password against the database. If the credentials are valid, it generates a JWT token using the `generateJWTToken` function and returns it to the client.

3. **Authentication Middleware**: The `authMiddleware` function checks if the incoming request has a valid JWT token in the `Authorization` header. If the token is valid, it extracts the user ID from the token claims and stores it in the request context using `c.Locals`. If the token is missing or invalid, it returns a 401 Unauthorized response.

4. **Authenticated Routes**: The authenticated routes (`/books`, `/books/:id`, etc.) are protected by the `authMiddleware`. Only authenticated users with a valid JWT token can access these routes.

5. **JWT Token Generation and Validation**: The `generateJWTToken` function generates a new JWT token with the user ID as the subject claim and a 24-hour expiration time. The `validateJWTToken` function validates a JWT token and returns the token claims if the token is valid.

The user registration and login work as follows:

1. To register a new user, send a POST request to `http://localhost:3000/register` with a JSON payload containing the `username` and `password`.

2. To log in, send a POST request to `http://localhost:3000/login` with a JSON payload containing the

