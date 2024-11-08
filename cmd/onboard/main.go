package main

import (
	"context"
	"log"
	"willd/commute-and-mute/internal/auth"
	"willd/commute-and-mute/internal/serialization"
	"willd/commute-and-mute/internal/strava"
	"willd/commute-and-mute/internal/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type ExchangeTokenResponse = struct {
	ID      int      `json:"id"`
	HomeLat *float64 `json:"hlat"`
	HomeLng *float64 `json:"hlng"`
	WorkLat *float64 `json:"wlat"`
	WorkLng *float64 `json:"wlng"`
	Token   string   `json:"token"`
}

func main() {
	lambda.Start(HandleTokenExchange)
}

func HandleTokenExchange(context context.Context, request *events.LambdaFunctionURLRequest) (ExchangeTokenResponse, error) {
	code := request.QueryStringParameters["code"]
	config, err := config.LoadDefaultConfig(context, config.WithRegion("eu-north-1"))
	if err != nil {
		log.Panicf("Error getting aws config: %v\n", err)

	}
	db, err := users.GetDbConnection(context)
	if err != nil {
		panic(err)
	}
	stravaClient := strava.NewStravaClient(strava.STRAVA_BASE_URL)
	authResponse, err := stravaClient.ExchangeToken(code)

	if err != nil {
		return ExchangeTokenResponse{}, err
	}

	userInDb, err := db.GetUser(context, authResponse.Athlete.ID)
	if err != nil {
		return ExchangeTokenResponse{}, err
	}

	userInDb = authResponse.AddToUser(userInDb)

	log.Printf("user: %v", userInDb)

	db.UpdateUser(context, userInDb)

	generator := auth.NewTokenGenerator(config)
	token := generator.GenerateForId(context, userInDb.ID)

	return ExchangeTokenResponse{
		userInDb.ID,
		serialization.GetNilFromNegative1(userInDb.HomeLat),
		serialization.GetNilFromNegative1(userInDb.HomeLng),
		serialization.GetNilFromNegative1(userInDb.WorkLat),
		serialization.GetNilFromNegative1(userInDb.WorkLng),
		token,
	}, nil
}
