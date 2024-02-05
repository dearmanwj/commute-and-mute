package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
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
	Map          Map
}

type Map struct {
	Polyline string
}

type ActivityUpdate struct {
	Commute        bool `json:"commute"`
	Hide_From_Home bool `json:"hide_from_home"`
}

type UserFormData struct {
	HomeLat float64 `json:"hlat"`
	HomeLng float64 `json:"hlng"`
	WorkLat float64 `json:"wlat"`
	WorkLng float64 `json:"wlng"`
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
		user, err := HandleTokenExchange(url.Query().Get("code"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		token := GenerateUserToken(user.ID)
		cookie := http.Cookie{
			Name:  "user-jwt",
			Value: token,
			Path:  "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/static/success.html", http.StatusSeeOther)
	} else if url.Path == "/app/user" {

		authHeader := r.Header.Get("Authorization")
		token := strings.Split(authHeader, "Bearer ")[1]
		id, err := GetConnectedUserId(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		r.ParseMultipartForm(512)
		var (
			hlat     float64
			hlng     float64
			wlat     float64
			wlng     float64
			parseErr error
		)
		hlat, parseErr = strconv.ParseFloat(r.Form.Get("hlat"), 64)
		hlng, parseErr = strconv.ParseFloat(r.Form.Get("hlng"), 64)
		wlat, parseErr = strconv.ParseFloat(r.Form.Get("wlat"), 64)
		wlng, parseErr = strconv.ParseFloat(r.Form.Get("wlng"), 64)
		if parseErr != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}
		HandleUserSubmitDetails(id, hlat, hlng, wlat, wlng)
		return
	} else {
		http.NotFound(w, r)
	}
}

func HandleTokenExchange(code string) (User, error) {
	auth, err := ExchangeToken(code)

	if err != nil {
		return User{}, err
	}

	user := auth.toUser()

	log.Printf("user: %v", user)

	UpdateUser(user)
	return user, nil
}

func ProcessActivity(a Activity) (err error) {
	if strings.EqualFold(a.Type, "ride") {
		log.Println("Received Ride activity")
		route, _ := maps.DecodePolyline(a.Map.Polyline)
		if isCommute(route[0].Lat, route[0].Lng, route[len(route)-1].Lat, route[len(route)-1].Lng) {
			log.Println("is ride and commute")
			sendCommuteAndMuteRequest(a)
			return nil
		} else {
			log.Println("ride not between home and work locations")
		}
	}
	return nil
}

func HandleUserSubmitDetails(id int, hlat float64, hlng float64, wlat float64, wlng float64) error {

	user, err := GetUser(id)
	if err != nil {
		return err
	}

	log.Printf("connected user: %v, %v, %v, %v, %v", user.ID, hlat, hlng, wlat, wlng)
	return nil
}

func sendCommuteAndMuteRequest(activity Activity) error {
	user, err := GetUser(activity.Athlete.ID)

	if err != nil {
		log.Printf("")
	}

	if user.ExpiresAt < time.Now().Unix() {
		log.Println("Token expired, refreshing")
		auth, err := RefreshToken(user.RefreshToken)
		if err != nil {
			return err
		}
		user = auth.toUser()
		UpdateUser(user)
	}

	toSend := ActivityUpdate{Commute: true, Hide_From_Home: true}

	client := &http.Client{}
	data, err := json.Marshal(toSend)
	if err != nil {
		log.Printf("Could not generate request body: %v\n", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPut,
		"https://www.strava.com/api/v3/activities/"+strconv.FormatInt(activity.Id, 10),
		bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error building update request: %v\n", err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending update request: %v\n", err)
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error updating activity, status: %v", resp.StatusCode)
	}

	log.Printf("Successfully updated activity: %v\n", activity.Id)
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
