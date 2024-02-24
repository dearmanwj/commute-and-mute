package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"willd/commute-and-mute/internal/users"
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

type StravaClient struct {
	baseUrl string
}

var STRAVA_BASE_URL = "https://www.strava.com"
var STRAVA_EXCHANGE_PATH = "/oauth/token"

func NewStravaClient(baseUrl string) StravaClient {
	return StravaClient{
		baseUrl: baseUrl,
	}
}

func (client StravaClient) ExchangeToken(code string) (AuthorizationResponse, error) {
	log.Println("Getting exchange token...")
	if code == "" {
		return AuthorizationResponse{}, errors.New("code is empty")
	}
	return client.makeTokenRequest(code, "authorization_code")
}

func (client StravaClient) RefreshToken(refreshToken string) (AuthorizationResponse, error) {
	log.Println("Refreshing access token...")
	return client.makeTokenRequest(refreshToken, "refresh_token")
}

func (auth AuthorizationResponse) ToUser() users.User {
	return users.User{
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

func (client StravaClient) makeTokenRequest(token string, grantType string) (AuthorizationResponse, error) {
	const clientId string = "116416"
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	queryParams := url.Values{}
	queryParams.Add("client_id", clientId)
	queryParams.Add("client_secret", clientSecret)
	queryParams.Add("code", token)
	queryParams.Add("grant_type", grantType)

	resp, err := http.Post(client.baseUrl+STRAVA_EXCHANGE_PATH+"?"+queryParams.Encode(), "text/plain", nil)
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
