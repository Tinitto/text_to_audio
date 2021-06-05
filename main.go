/*
* A simple application allowing users to send text to an endpoint and receive back
* an audio they can play
 */
package main

import (
	"fmt"
	"html/template"
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
		Templates : template.Must(template.ParseFiles(
			"templates/head.html",
			"templates/footer.html",
			"templates/error.html", 
			"templates/login.html", 
			"templates/home.html", 
			)),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login", app.LoginPage)
	mux.HandleFunc("/error", app.ErrorPage)
	mux.HandleFunc("/", app.HomePage)

	fmt.Println("Starting the serer on 0.0.0.0:8000. Hit Ctrl-C to stop.")
	log.Fatal(http.ListenAndServe(":8000", mux))
	
}
