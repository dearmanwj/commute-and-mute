package main

import (
	"context"
	"log"
	"willd/commute-and-mute/internal/auth"
	"willd/commute-and-mute/internal/strava"
	"willd/commute-and-mute/internal/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type ExchangeTokenResponse = struct {
	ID      int
	HomeLat float64
	HomeLng float64
	WorkLat float64
	WorkLng float64
	Token   string
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
	users.GetDbConnection(context)
	stravaClient := strava.NewStravaClient(strava.STRAVA_BASE_URL)
	authResponse, err := stravaClient.ExchangeToken(code)

	if err != nil {
		return ExchangeTokenResponse{}, err
	}

	user := authResponse.ToUser()

	log.Printf("user: %v", user)

	users.UpdateUser(context, user)

	generator := auth.NewTokenGenerator(config)
	token := generator.GenerateForId(context, user.ID)

	return ExchangeTokenResponse{
		user.ID,
		user.HomeLat,
		user.HomeLng,
		user.WorkLat,
		user.WorkLng,
		token,
	}, nil
}
