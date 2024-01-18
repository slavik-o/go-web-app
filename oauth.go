package main

import (
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func SetupOAuth() {
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CALLBACK_URL"),
			"profile", "email", "openid",
		),
	)
}
