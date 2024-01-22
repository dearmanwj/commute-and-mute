package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	Athlete      Athlete
}

type ActivityUpdate struct {
	Commute      bool
	HideFromHome bool
}

func main() {

	godotenv.Load(".env")

	http.HandleFunc("/app/", handlerHttp)
	fs := http.FileServer(http.Dir("./static/"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

func handlerHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("In http handler")
	log.Println(r.URL)

	url := r.URL
	getDbConnection()
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
	auth := ExchangeToken(code)
	// auth := AuthorizationResponse{
	// 	Refresh_Token: "token",
	// 	Token_Type:    "type",
	// 	Expires_At:    12345,
	// 	Expires_In:    123,
	// 	Access_Token:  "token2",
	// 	Athlete: Athlete{
	// 		UserName: "wdearman",
	// 		ID:       10503812,
	// 	},
	// }

	user := User{
		ID:           auth.Athlete.ID,
		AccessToken:  auth.Access_Token,
		RefreshToken: auth.Refresh_Token,
		HomeLat:      HomeLat,
		HomeLng:      HomeLng,
		WorkLat:      WorkLat,
		WorkLng:      WorkLng,
	}

	log.Printf("user: %v", user)

	UpdateUser(user)
}

func ProcessActivity(a Activity) (err error) {
	if strings.EqualFold(a.Type, "ride") {
		log.Println("Received Ride activity")
		if isCommute(a.Start_latlng[0], a.Start_latlng[1], a.End_latlng[0], a.End_latlng[1]) {
			log.Println("is ride and commute")
			toSend := ActivityUpdate{Commute: true, HideFromHome: true}
			fmt.Printf("To send: %v", toSend)
			user, _ := GetUser(a.Athlete.ID)
			log.Printf("Retrieved user: %v\n", user.ID)

			client := &http.Client{}
			data, err := json.Marshal(toSend)
			if err != nil {
				log.Printf("Could not generate request body: %v\n", err)
				return err
			}
			req, err := http.NewRequest(http.MethodPut,
				"https://www.strava.com/api/v3/activities/"+strconv.FormatInt(a.Id, 10),
				bytes.NewBuffer(data))
			if err != nil {
				log.Printf("Error building update request: %v\n", err)
				return err
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error sending update request: %v\n", err)
				return err
			}
			log.Printf("Successfully updated activity: %v\n", resp.StatusCode)
		} else {
			log.Println("ride not between home and work locations")
		}
	}
	return nil
}

func updateActivity(activityId string, user User)

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
