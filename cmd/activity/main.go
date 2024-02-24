package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"willd/commute-and-mute/internal/strava"
	"willd/commute-and-mute/internal/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"googlemaps.github.io/maps"
)

type StravaEvent struct {
	AspectType     string  `json:"aspect_type"`
	EventTime      int64   `json:"event_time"`
	ObjectId       int64   `json:"object_id"`
	ObjectType     string  `json:"object_type"`
	OwnerId        int64   `json:"owner_id"`
	SubscriptionId int64   `json:"subscription_id"`
	Updates        Updates `json:"updates"`
}

type Updates struct {
	Title      string `json:"title"`
	UpdateType string `json:"type"`
	Private    bool   `json:"private"`
}

func main() {
	lambda.Start(handleNewActivity)
}

func handleNewActivity(ctx context.Context, request *events.LambdaFunctionURLRequest) (*string, error) {
	if request.RequestContext.HTTP.Method == "POST" {
		update, err := DecodeUpdateEvent(request.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decode event: %v", err)
		}

	}
	return nil, nil
}

func DecodeUpdateEvent(rawEvent string) (StravaEvent, error) {
	var update StravaEvent
	err := json.Unmarshal([]byte(rawEvent), &update)
	if err != nil {
		return StravaEvent{}, err
	}
	return update, nil
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
		return err
	}

	stravaClient := strava.NewStravaClient(strava.STRAVA_BASE_URL)

	if user.ExpiresAt < time.Now().Unix() {
		log.Println("Token expired, refreshing")
		authResponse, err := stravaClient.RefreshToken(user.RefreshToken)
		if err != nil {
			return err
		}
		user = authResponse.ToUser()
		users.UpdateUser(user)
	}

	return stravaClient.MakeActivityUpdateRequest(activity.Id, user.AccessToken)
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
