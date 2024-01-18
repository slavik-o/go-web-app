package main

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
}
