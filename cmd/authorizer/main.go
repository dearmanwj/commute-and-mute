package main

import (
	"context"
	"errors"
	"log"
	"strconv"
	"willd/commute-and-mute/internal/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	lambda.Start(HandleAuth)
}

func HandleAuth(context context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

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
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	} else {
		log.Println("token valid")
		return events.APIGatewayCustomAuthorizerResponse{
			PrincipalID:    strconv.Itoa(id),
			PolicyDocument: generatePolicy(request.MethodArn, true),
		}, nil
	}

}

func generatePolicy(resource string, allow bool) events.APIGatewayCustomAuthorizerPolicy {
	var effect string
	if allow {
		effect = "Allow"
	} else {
		effect = "Deny"
	}
	statement := events.IAMPolicyStatement{
		Action:   []string{"execute-api:Invoke"},
		Effect:   effect,
		Resource: []string{resource},
	}
	return events.APIGatewayCustomAuthorizerPolicy{
		Statement: []events.IAMPolicyStatement{statement},
	}
}
