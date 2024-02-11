package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
)

type AuthorizationResponse struct {
	Token_Type    string
	Expires_At    int64
	Expires_In    int64
	Refresh_Token string
	Access_Token  string
	Athlete       Athlete
}

type Athlete struct {
	UserName string
	ID       int
}

func ExchangeToken(code string) (AuthorizationResponse, error) {
	log.Println("Getting exchange token...")
	return makeTokenRequest(code, "authorization_code")
}

func RefreshToken(refreshToken string) (AuthorizationResponse, error) {
	log.Println("Refreshing access token...")
	return makeTokenRequest(refreshToken, "refresh_token")
}

func (auth AuthorizationResponse) toUser() User {
	return User{
		ID:           auth.Athlete.ID,
		AccessToken:  auth.Access_Token,
		RefreshToken: auth.Refresh_Token,
		HomeLat:      -1,
		HomeLng:      -1,
		WorkLat:      -1,
		WorkLng:      -1,
		ExpiresAt:    auth.Expires_At,
	}
}

func makeTokenRequest(token string, grantType string) (AuthorizationResponse, error) {
	const clientId string = "116416"
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	queryParams := url.Values{}
	queryParams.Add("client_id", clientId)
	queryParams.Add("client_secret", clientSecret)
	queryParams.Add("code", token)
	queryParams.Add("grant_type", grantType)

	resp, err := http.Post("https://www.strava.com/oauth/token?"+queryParams.Encode(), "text/plain", nil)
	var auth AuthorizationResponse
	if err != nil {
		log.Printf("Error making token exchange request")
		return auth, err
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&auth)
	log.Println("Successfully obtained token")
	return auth, nil
}
