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
	"strconv"
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
	portString := os.Getenv("PORT")
	port, err := strconv.ParseInt(portString, 10, 64)
	if err != nil {
		log.Fatal("The PORT has not been set in your environment")
	}

	app := utils.App{
		AllowedEmails: allowedEmailsMap,
		JWTSecret: os.Getenv("JWT_SECRET"),
		GoogleClientId: os.Getenv("GOOGLE_CLIENT_ID"),
		Port: port,
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

	fmt.Printf("Starting the server on 0.0.0.0:%d. Hit Ctrl-C to stop.\n", app.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", app.Port), mux))
	
}
