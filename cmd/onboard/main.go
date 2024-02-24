package main

import (
	"context"
	"log"
	"willd/commute-and-mute/internal/strava"
	"willd/commute-and-mute/internal/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleTokenExchange)
}

func HandleTokenExchange(context context.Context, request *events.LambdaFunctionURLRequest) (users.User, error) {
	code := request.QueryStringParameters["code"]
	users.GetDbConnection()
	stravaClient := strava.NewStravaClient(strava.STRAVA_BASE_URL)
	auth, err := stravaClient.ExchangeToken(code)

	if err != nil {
		return users.User{}, err
	}

	user := auth.ToUser()

	log.Printf("user: %v", user)

	users.UpdateUser(user)
	return user, nil
}
