/*
* A simple application allowing users to send text to an endpoint and receive back
* an audio they can play
 */
package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"sopherapps.com/text_to_audio/utils"
)

func main()  {
	allowedEmails := strings.Split(os.Getenv("ALLOWED_EMAILS"), ",")
	allowedEmailsMap := map[string]struct{}{}

	for _, email := range allowedEmails {
		allowedEmailsMap[email] = struct{}{}
	}
	app := utils.App{
		AllowedEmails: allowedEmailsMap,
		JWTSecret: os.Getenv("JWT_SECRET"),
		GoogleClientId: os.Getenv("GOOGLE_CLIENT_ID"),
	}

	http.HandleFunc("/", app.HomePage)
	http.HandleFunc("/login", app.LoginPage)
	http.HandleFunc("/error", app.ErrorPage)

	log.Fatal(http.ListenAndServe(":8000", nil))
	
}
