package auth

import (
	"app-auth/config"

	"fmt"
	"log"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
)

func Goth() string {
	redirectUrl := os.Getenv("ENDPOINT"); if redirectUrl == "" {
		redirectUrl = "https://id.scaratec.com"
	}

	// just print the current redirect URL from the process environment
	log.Print(redirectUrl)

	facebookRedirectUrl := fmt.Sprintf(`%s/auth/%s/callback?provider=%s`, redirectUrl, "facebook", "facebook")
	googleRedirectUrl := fmt.Sprintf(`%s/auth/%s/callback?provider=%s`, redirectUrl, "google", "google")

	goth.UseProviders(
		google.New(config.GoogleClientKey, config.GoogleSecret, googleRedirectUrl),
		facebook.New(config.FacebookClientKey, config.FacebookSecret, facebookRedirectUrl),
	)

	return ""
}