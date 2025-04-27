package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
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

type ChallengeResponse struct {
	HubChallenge string `json:"hub.challenge"`
}

func main() {
	lambda.Start(handleNewActivity)
}

func handleNewActivity(ctx context.Context, request *events.LambdaFunctionURLRequest) (*string, error) {
	var message string
	var err error
	if request.RequestContext.HTTP.Method == "POST" {
		update, err := DecodeUpdateEvent(request.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decode event: %v", err)
		}
		log.Printf("handling update event: %v", update)
		if update.ObjectType == "activity" && update.AspectType == "create" {
			err = ProcessActivity(ctx, update)
			if err != nil {
				message = "error processing activity"
			} else {
				message = "activity handled successfully"
			}
		} else {
			message = "update event not a create activity, no action required"
		}
	} else if request.RequestContext.HTTP.Method == "GET" {
		log.Printf("confirming webhook subscription")
		verifyToken, present := os.LookupEnv("WEBHOOK_VERIFY_TOKEN")
		if !present {
			log.Panicln("Webhook verification not configured")
		}
		if request.QueryStringParameters["hub.mode"] == "subscribe" &&
			request.QueryStringParameters["hub.verify_token"] == verifyToken {
			challengeResponse := ChallengeResponse{HubChallenge: request.QueryStringParameters["hub.challenge"]}
			responseBytes, _ := json.Marshal(challengeResponse)
			message = string(responseBytes)
		} else {
			message = "webhook verification failed"
			err = errors.New("invalid parameters")
		}
	} else {
		message = "request not handled"
		err = errors.New("unrecognized request")
	}
	return &message, err
}

func DecodeUpdateEvent(rawEvent string) (StravaEvent, error) {
	var update StravaEvent
	err := json.Unmarshal([]byte(rawEvent), &update)
	if err != nil {
		return StravaEvent{}, err
	}
	return update, nil
}

func GetBearerToken(db users.TableBasics, ctx context.Context, user users.User, stravaClient *strava.StravaClient) string {
	if user.ExpiresAt < time.Now().Unix() {
		log.Println("Token expired, refreshing")
		authResponse, err := stravaClient.RefreshToken(user.RefreshToken)
		if err != nil {
			log.Panicf("could not refresh user token: %v", err)
		}
		user.AccessToken = authResponse.Access_Token
		user.ExpiresAt = authResponse.Expires_At
		db.UpdateUser(ctx, user)
	}

	return user.AccessToken
}

func ProcessActivity(ctx context.Context, update StravaEvent) (err error) {
	db, err := users.GetDbConnection(ctx)
	if err != nil {
		panic(err)
	}
	user, err := db.GetUser(ctx, int(update.OwnerId))
	if err != nil || user.ID == 0 {
		return fmt.Errorf("error retrieving user with id: %v, %v", update.OwnerId, err)
	}

	stravaClient := strava.NewStravaClient(strava.STRAVA_BASE_URL)

	token := GetBearerToken(db, ctx, user, &stravaClient)

	newActivity, err := stravaClient.GetActivity(update.ObjectId, token)
	if err != nil {
		return err
	}

	if strings.EqualFold(newActivity.Type, "ride") {
		log.Println("Received Ride activity")
		route, _ := maps.DecodePolyline(newActivity.Map.Polyline)
		startLat := strconv.FormatFloat(route[0].Lat, 'E', -1, 64)
		startLng := strconv.FormatFloat(route[0].Lng, 'E', -1, 64)
		endLat := strconv.FormatFloat(route[len(route)-1].Lat, 'E', -1, 64)
		endLng := strconv.FormatFloat(route[len(route)-1].Lng, 'E', -1, 64)
		if isCommute(startLat, startLng, endLat, endLng, user) {
			log.Println("is ride and commute")
			err := stravaClient.UpdateActivity(newActivity.Id, token)
			if err != nil {
				return err
			}
			user.Counter++
			db.UpdateUser(ctx, user)
			return nil
		} else {
			log.Println("ride not between home and work locations")
		}
	}
	return nil
}

func isCommute(startLat, startLng, endLat, endLng string, user users.User) bool {
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
