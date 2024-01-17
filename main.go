package main

import (
	"log"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/slavik-o/go-web-app/views"
)

// User model
type User struct {
	gorm.Model
	ID         string `gorm:"primaryKey"`
	Provider   string
	ProviderID string
	Email      string
	FirstName  string
	LastName   string
}

func UpsertAuthUser(db *gorm.DB, authUser goth.User) (*User, error) {
	user := new(User)

	result := db.Where("provider = ? AND provider_id = ?", authUser.Provider, authUser.UserID).Find(user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected < 1 {
		uuid, err := uuid.NewRandom()
		if err != nil {
			log.Fatalf("Failed to generate UUID: %v", err)
		}

		user.ID = uuid.String()
	}

	user.Provider = authUser.Provider
	user.ProviderID = authUser.UserID
	user.Email = authUser.Email
	user.FirstName = authUser.FirstName
	user.LastName = authUser.LastName

	result = db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FindUserByID(db *gorm.DB, ID string) (*User, error) {
	user := new(User)

	result := db.First(user, "id = ?", ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// Auth middleware
type AuthMiddleware struct {
	DB *gorm.DB
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		DB: db,
	}
}

func (m *AuthMiddleware) AuthRequired(c *fiber.Ctx) error {
	if c.Cookies("user") == "" {
		return c.Redirect("/login")
	}

	user, err := FindUserByID(m.DB, c.Cookies("user"))
	if err != nil {
		return c.Redirect("/login")
	}

	c.Locals("user", user)

	return c.Next()
}

func Render(ctx *fiber.Ctx, cmp templ.Component) error {
	ctx.Set("Content-Type", "text/html")
	return cmp.Render(ctx.Context(), ctx.Response().BodyWriter())
}

// Main
func main() {
	// Env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Database
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	db.AutoMigrate(&User{}) // Don't use in production

	// OAuth
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CALLBACK_URL"),
			"profile", "email", "openid",
		),
	)

	// App
	app := fiber.New()

	// Assets
	app.Static("/", "./assets")

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_SECRET"),
	}))

	auth := NewAuthMiddleware(db)

	// Routes
	app.Get("/", auth.AuthRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		return Render(c, views.Index(user.FirstName))
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return Render(c, views.Login())
	})

	app.Get("/auth/:provider", func(c *fiber.Ctx) error {
		return goth_fiber.BeginAuthHandler(c)
	})

	app.Get("/auth/:provider/callback", func(c *fiber.Ctx) error {
		authUser, err := goth_fiber.CompleteUserAuth(c)
		if err != nil {
			return c.Redirect("/login")
		}

		user, err := UpsertAuthUser(db, authUser)
		if err != nil {
			return c.Redirect("/login")
		}

		c.Cookie(&fiber.Cookie{
			Name:    "user",
			Value:   user.ID,
			Expires: time.Now().Add(7 * 24 * time.Hour),
		})

		return c.Redirect("/")
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("user")

		return c.Redirect("/login")
	})

	// Run
	log.Fatal(app.Listen(os.Getenv("BIND_URL")))
}
