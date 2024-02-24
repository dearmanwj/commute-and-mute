package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"willd/commute-and-mute/internal/strava"
	"willd/commute-and-mute/internal/users"

	"github.com/aws/aws-lambda-go/lambda"
	"googlemaps.github.io/maps"
)

func main() {
	lambda.Start(handleNewActivity)
}

func handleNewActivity(ctx context.Context, event *strava.Activity) (*string, error) {
	ProcessActivity(*event)
	return nil, nil
}

func ProcessActivity(a strava.Activity) (err error) {
	users.GetDbConnection()
	user, err := users.GetUser(a.Athlete.ID)
	if err != nil {
		log.Panicf("error retrieving user with id: %v", a)
	}
	if strings.EqualFold(a.Type, "ride") {
		log.Println("Received Ride activity")
		route, _ := maps.DecodePolyline(a.Map.Polyline)
		if isCommute(route[0].Lat, route[0].Lng, route[len(route)-1].Lat, route[len(route)-1].Lng, user) {
			log.Println("is ride and commute")
			sendCommuteAndMuteRequest(a)
			return nil
		} else {
			log.Println("ride not between home and work locations")
		}
	}
	return nil
}

func sendCommuteAndMuteRequest(activity strava.Activity) error {
	user, err := users.GetUser(activity.Athlete.ID)

	if err != nil {
		log.Printf("")
	}

	if user.ExpiresAt < time.Now().Unix() {
		log.Println("Token expired, refreshing")
		stravaClient := strava.NewStravaClient(strava.STRAVA_BASE_URL)
		authResponse, err := stravaClient.RefreshToken(user.RefreshToken)
		if err != nil {
			return err
		}
		user = authResponse.ToUser()
		users.UpdateUser(user)
	}

	toSend := strava.ActivityUpdate{Commute: true, Hide_From_Home: true}

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

func GetUser(i int) {
	panic("unimplemented")
}

func isCommute(startLat, startLng, endLat, endLng float64, user users.User) bool {
	var isCommute bool = false
	isHomeStart := IsWithinRadius(user.HomeLat, user.HomeLng, startLat, startLng)
	if isHomeStart {
		// if home is start, is commute if end is work
		isCommute = IsWithinRadius(user.WorkLat, user.WorkLng, endLat, endLng)
	} else {
		isHomeEnd := IsWithinRadius(user.HomeLat, user.HomeLng, endLat, endLng)
		if isHomeEnd {
			// if home is end, is commute if start is work
			isCommute = IsWithinRadius(user.WorkLat, user.WorkLng, startLat, startLng)
		}
	}
	return isCommute
}
