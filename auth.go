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
}

func ExchangeToken(code string) AuthorizationResponse {
	const clientId string = "116416"
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")
	log.Println("Getting exchange token...")

	queryParams := url.Values{}
	queryParams.Add("client_id", clientId)
	queryParams.Add("client_secret", clientSecret)
	queryParams.Add("code", code)
	queryParams.Add("grant_type", "authorization_code")

	resp, err := http.Post("https://www.strava.com/oauth/token?"+queryParams.Encode(), "text/plain", nil)
	if err != nil {
		log.Fatal("Error making token exchange request")
	}
	defer resp.Body.Close()
	var auth AuthorizationResponse
	json.NewDecoder(resp.Body).Decode(&auth)
	log.Println("Successfully exchanged token")
	return auth
}
