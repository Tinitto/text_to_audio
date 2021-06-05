package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Example handler for logging in if one was to then generate their own token
//  func (cfg config) loginHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()

// 	// parse the GoogleJWT that was POSTed from the front-end
// 	type parameters struct {
// 		GoogleJWT *string
// 	}
// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, 500, "Couldn't decode parameters")
// 		return
// 	}

// 	// Validate the JWT is valid
// 	claims, err := auth.ValidateGoogleJWT(*params.GoogleJWT)
// 	if err != nil {
// 		respondWithError(w, 403, "Invalid google auth")
// 		return
// 	}
// 	if claims.Email != user.Email {
// 		respondWithError(w, 403, "Emails don't match")
// 		return
// 	}

// 	// create a JWT for OUR app and give it back to the client for future requests
// 	tokenString, err := auth.MakeJWT(claims.Email, cfg.JWTSecret)
// 	if err != nil {
// 		respondWithError(w, 500, "Couldn't make authentication token")
// 		return
// 	}

// 	respondWithJSON(w, 200, struct {
// 		Token string `json:"token"`
// 	}{
// 		Token: tokenString,
// 	})
// }

// Verifies that the JWT token is valid
func VerifyJWTToken(r *http.Request, jwtSecret string, allowedEmails map[string]struct{}) (bool, error) {
	tokenString := extractToken(r)

	token, err := jwt.Parse(
		tokenString, 
		func(token *jwt.Token) (interface{}, error) {
	   //Make sure that the token method conform to "SigningMethodHMAC"
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		  return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	   }
	   return []byte(jwtSecret), nil
	})
	if err != nil {
	   return false, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false, err
	}

	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		email := mapClaims["email"].(string)
		_, authorized := allowedEmails[email]
		return authorized, nil
	}

	return false, nil
}

// Extracts token from the authorization header
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
	   return strings.Trim(strArr[1], "")
	}
	return ""
}


// Generates a JWT specific to this site
func MakeJWT(email string, secret string, allowedEmails map[string]struct{}) (string, error) {
	if _, authorized := allowedEmails[email]; !authorized {
		return "", errors.New("email not allowed")
	}

	claims := jwt.MapClaims{
		"email": email,
		"exp": time.Now().UTC().Add(time.Hour*2).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validates the JWT token passed by the client to ensure it is still valid according to Google
func ValidateGoogleJWT(tokenString string, clientId string) (GoogleClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, 
		&GoogleClaims{}, 
		func(t *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(t.Header["kid"].(string))
			if err != nil {
				return nil, err
			}

			key, err := jwt.ParseECPrivateKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err 
			}

			return key, nil
	})
	if err != nil {
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{}, errors.New("invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	if claims.Audience != clientId {
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return GoogleClaims{}, errors.New("JWT is expired")
	}

	return GoogleClaims{}, nil
}

// The structure of the claims Google is to add to its JWT
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

// Google hosts their public key over HTTPS. 
// Each time we need to verify a request we can go grab their public key as follows
func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}
