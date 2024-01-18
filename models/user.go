package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         string `gorm:"primaryKey"`
	Provider   string
	ProviderID string
	Email      string
	FirstName  string
	LastName   string
}
