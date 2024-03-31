package main

import (
	"context"
	"log"
	"willd/commute-and-mute/internal/strava"
	"willd/commute-and-mute/internal/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	users.GetDbConnection(context)
	stravaClient := strava.NewStravaClient(strava.STRAVA_BASE_URL)
	auth, err := stravaClient.ExchangeToken(code)

	if err != nil {
		return ExchangeTokenResponse{}, err
	}

	user := auth.ToUser()

	log.Printf("user: %v", user)

	users.UpdateUser(context, user)
	return ExchangeTokenResponse{
		user.ID,
		user.HomeLat,
		user.HomeLng,
		user.WorkLat,
		user.WorkLng,
		"",
	}, nil
}
