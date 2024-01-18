package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const HomeLat float64 = 43.593583
const HomeLng float64 = 1.448228
const WorkLat float64 = 43.564060
const WorkLng float64 = 1.389155
const RadiusKM float64 = 0.1

type Activity struct {
	Id           int64
	Type         string
	Start_latlng [2]float64
	End_latlng   [2]float64
}

type ActivityUpdate struct {
	Commute      bool
	HideFromHome bool
}

type AuthorizationResponse struct {
	Token_Type    string
	Expires_At    int64
	Expires_In    int64
	Refresh_Token string
	Access_Token  string
}

func main() {

	godotenv.Load(".env")

	// sess, _ := session.NewSession(&aws.Config{
	// 	Region: aws.String("eu-north-1")},
	// )

	// svc := dynamodb.New(sess)

	http.HandleFunc("/app/", handlerHttp)
	fs := http.FileServer(http.Dir("./static/"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

func handlerHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("In http handler")
	log.Println(r.URL)

	url := r.URL
	if url.Path == "/app/activity" {
		var a Activity
		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ProcessActivity(a)
		return
	} else if url.Path == "/app/exchange_token" {
		log.Println("Exchanging token")
		HandleTokenExchange(url.Query().Get("code"))
	} else {
		http.NotFound(w, r)
	}
}

func HandleTokenExchange(code string) {
	const clientId int = 116416
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")
	log.Println("Getting exchange token...")

	queryParams := url.Values{}
	queryParams.Add("client_id", strconv.FormatInt(int64(clientId), 10))
	queryParams.Add("client_secret", clientSecret)
	queryParams.Add("code", code)
	queryParams.Add("grant_type", "authorization_code")

	req, _ := http.NewRequest("POST", "https://www.strava.com/oauth/token?"+queryParams.Encode(), nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	var auth AuthorizationResponse
	json.NewDecoder(resp.Body).Decode(&auth)
	log.Printf("Authorization: %v", auth)
}

func ProcessActivity(a Activity) (err error) {
	if strings.EqualFold(a.Type, "ride") {
		log.Println("Received Ride activity")
		if isCommute(a.Start_latlng[0], a.Start_latlng[1], a.End_latlng[0], a.End_latlng[1]) {
			log.Println("is ride and commute")
			toSend := ActivityUpdate{Commute: true, HideFromHome: true}
			fmt.Printf("To send: %v", toSend)
		} else {
			log.Println("ride not between home and work locations")
		}
	}
	return nil
}

func isCommute(startLat, startLng, endLat, endLng float64) bool {
	var isCommute bool = false
	isHomeStart := IsWithinRadius(HomeLat, HomeLng, startLat, startLng)
	if isHomeStart {
		// if home is start, is commute if end is work
		isCommute = IsWithinRadius(WorkLat, WorkLng, endLat, endLng)
	} else {
		isHomeEnd := IsWithinRadius(HomeLat, HomeLng, endLat, endLng)
		if isHomeEnd {
			// if home is end, is commute if start is work
			isCommute = IsWithinRadius(WorkLat, WorkLng, startLat, startLng)
		}
	}
	return isCommute
}
