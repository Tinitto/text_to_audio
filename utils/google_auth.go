package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
)

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

// Gets the token info from google
func ValidateGoogleJWT(ctx context.Context, idToken string, clientId string) (*oauth2.Tokeninfo, error) {
	httpClient := &http.Client{}
    oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
    tokenInfoCall := oauth2Service.Tokeninfo()    
    tokenInfo, err := tokenInfoCall.IdToken(idToken).Do()
    if err != nil {
        return nil, err
    }

	if tokenInfo.Issuer != "accounts.google.com" && tokenInfo.Issuer != "https://accounts.google.com" {
		return nil, errors.New("iss is invalid")
	}

	if tokenInfo.Audience != clientId {
		return nil, errors.New("aud is invalid")
	}

	if tokenInfo.ExpiresIn < 2 {
		return nil, errors.New("JWT is expired")
	}

    return tokenInfo, nil
}
