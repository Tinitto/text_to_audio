package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var templates = template.Must(template.ParseFiles(
	"../templates/base.html", 
	"../templates/error.html", 
	"../templates/home.html",
	"../templates/login.html",
	))

type Page struct {
	Body string
}


type App struct {
	AllowedEmails map[string]struct{}
	JWTSecret string
	GoogleClientId string
}

// Handles the logging in of users
// verify the google token 
// return to client a JWT token if google token is valid
func (app *App) LoginPage(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		respondWithTemplate(w, "../templates/login.html", "")
		return
	}

	defer r.Body.Close()

	// parse the GoogleJWT that was POSTed from the front-end
	type parameters struct {
		GoogleJWT *string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, "error retrieving your token", 500)
		return
	}

	// Validate the JWT is valid
	claims, err := ValidateGoogleJWT(*params.GoogleJWT, app.GoogleClientId)
	if err != nil {
		respondWithError(w, "invalid google auth", 403)
		return
	}

	// create a JWT for OUR app and give it back to the client for future requests
	tokenString, err := MakeJWT(claims.Email, app.JWTSecret, app.AllowedEmails)
	if err != nil {
		respondWithError(w, "couldn't make authentication token", 500)
		return
	}

	respondWithJSON(w, map[string]string {"token": tokenString}, 200)
}

// Handles the showing of the error page
// Extract the message of from the query param and display it in error page
func (app *App) ErrorPage(w http.ResponseWriter, r *http.Request)  {
	message := r.URL.Query().Get("msg")
	respondWithTemplate(w, "../templates/error.html", message)
}

// Handles the home page where the actual conversion from text to audio occurs
// receive text from JSON body, call Google texttospeech API, return audio as a blob
func (app *App) HomePage(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		respondWithTemplate(w, "../templates/home.html", "")
		return
	}
	defer r.Body.Close()

	isAuthorized, err := VerifyJWTToken(r, app.JWTSecret, app.AllowedEmails)
	if err != nil {
		respondWithTemplate(w, "../templates/error.html", "failed to authenticate")
		return
	}

	if !isAuthorized {
		respondWithTemplate(w, "../templates/error.html", "Sorry. You are not known.")
	}
	
	type parameters struct {
		Text string `json:"text"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithTemplate(w, "../templates/error.html", "error extracting your text.")
		return
	}

	ctx := r.Context()
	audioFile, err := ConvertTextToAudio(ctx, params.Text, "en-US")
	if err != nil {
		respondWithTemplate(w, "../templates/error.html", "error converting your text.")
		return
	}

	filename := fmt.Sprintf("sop_audio_%s.mp3", time.Now().UTC().Format("200601002150405"))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(audioFile)))
	
	_, err = w.Write(audioFile)
	if err != nil {
		respondWithTemplate(w, "../templates/error.html", "error downloading the audio file.")
	}
}

// Returns error in json
func respondWithError(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// Responds to the client in JSON
func respondWithJSON(w http.ResponseWriter, data interface{}, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(data)
}

// Responds with a given template and body
func respondWithTemplate(w http.ResponseWriter, templatePath string, body string)  {
	err := templates.ExecuteTemplate(w, templatePath, &Page{Body: body})	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
