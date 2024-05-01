package main

import (
	"context"
	"errors"
	"log"
	"willd/commute-and-mute/internal/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	lambda.Start(HandleAuth)
}

func HandleAuth(context context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {

	log.Printf("received authorization request for %v", request.MethodArn)

	config, err := config.LoadDefaultConfig(context, config.WithRegion("eu-north-1"))
	if err != nil {
		log.Panicf("Error getting aws config: %v\n", err)
	}

	token := request.AuthorizationToken

	generator := auth.NewTokenGenerator(config)

	id, err := generator.GetIdIfValid(context, token)
	if err != nil {
		log.Printf("error authorizing user: %v", err)
		unauthorizedResponse := events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
			Context:      make(map[string]interface{}),
		}
		return unauthorizedResponse, errors.New("Unauthorized")
	} else {
		log.Println("token valid")
		authorizedContext := map[string]interface{}{"user": id}
		authorizedResponse := events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: true,
			Context:      authorizedContext,
		}
		return authorizedResponse, nil
	}

}
