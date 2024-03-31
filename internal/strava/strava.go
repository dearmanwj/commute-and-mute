package strava

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

type Activity struct {
	Id           int64
	Type         string
	Start_latlng [2]float64
	End_latlng   [2]float64
	Athlete      Athlete
	Map          Map
}

type ActivityUpdate struct {
	Commute        bool `json:"commute"`
	Hide_From_Home bool `json:"hide_from_home"`
}

type Map struct {
	Polyline string
}

type StravaClient struct {
	baseUrl string
}

var STRAVA_BASE_URL = "https://www.strava.com"
var STRAVA_EXCHANGE_PATH = "/api/v3/oauth/token"
var STRAVA_ACTIVITY_PATH = "/api/v3/activities/%v"

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
	if refreshToken == "" {
		return AuthorizationResponse{}, errors.New("refresh token is empty")
	}
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
	if grantType == "authorization_code" {
		queryParams.Add("code", token)
	} else if grantType == "refresh_token" {
		queryParams.Add("refresh_token", token)
	}
	queryParams.Add("grant_type", grantType)

	resp, err := http.Post(client.baseUrl+STRAVA_EXCHANGE_PATH+"?"+queryParams.Encode(), "text/plain", nil)
	var auth AuthorizationResponse
	if err != nil {
		log.Printf("Error making token exchange request")
		return auth, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return auth, fmt.Errorf("error authorizing user: %v", resp.Status)
	}
	json.NewDecoder(resp.Body).Decode(&auth)
	log.Println("Successfully obtained token")
	return auth, nil
}

func (client StravaClient) GetActivity(activityId int64, bearerToken string) (Activity, error) {
	httpClient := http.DefaultClient
	url := fmt.Sprintf(client.baseUrl+STRAVA_ACTIVITY_PATH, activityId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic("cannot build activity get request")
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Panic("could not execute get activity request")
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return Activity{}, errors.New("unauthorized")
	}
	defer resp.Body.Close()
	var activity Activity
	err = json.NewDecoder(resp.Body).Decode(&activity)
	if err != nil {
		log.Panicf("could not decode activity response: %v", err)
	}
	return activity, nil
}

func (client StravaClient) UpdateActivity(activityId int64, bearerToken string) error {
	httpClient := http.DefaultClient
	url := fmt.Sprintf(client.baseUrl+STRAVA_ACTIVITY_PATH, activityId)
	body, err := json.Marshal(ActivityUpdate{Commute: true, Hide_From_Home: true})
	if err != nil {
		log.Panicf("cannot build activity update request body, %v", err)
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		log.Panicf("cannot build activity update request, %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error sending update request: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error updating activity, status: %v", resp.StatusCode)
	}

	log.Printf("Successfully updated commute activity: %v, status code: %v\n", activityId, resp.Status)
	return nil
}
