package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/slavik-o/go-web-app/managers"
	"github.com/slavik-o/go-web-app/models"
)

type Managers struct {
	*managers.UserManager
}

func Connect() *Managers {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	db.AutoMigrate(&models.User{}) // Don't use in production

	return &Managers{
		&managers.UserManager{
			DB: db,
		},
	}
}
