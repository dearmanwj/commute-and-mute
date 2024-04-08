package main

import (
	"context"
	"fmt"
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

	config, err := config.LoadDefaultConfig(context, config.WithRegion("eu-north-1"))
	if err != nil {
		log.Panicf("Error getting aws config: %v\n", err)
	}

	token := request.AuthorizationToken

	generator := auth.NewTokenGenerator(config)

	id, err := generator.GetIdIfValid(context, token)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("unauthorized, %w", err)
	} else {
		return events.APIGatewayCustomAuthorizerResponse{
			PrincipalID:    strconv.Itoa(id),
			PolicyDocument: generatePolicy(request.MethodArn),
		}, nil
	}

}

func generatePolicy(resource string) events.APIGatewayCustomAuthorizerPolicy {
	statement := events.IAMPolicyStatement{
		Action:   []string{"execute-api:Invoke"},
		Effect:   "Allow",
		Resource: []string{resource},
	}
	return events.APIGatewayCustomAuthorizerPolicy{
		Statement: []events.IAMPolicyStatement{statement},
	}
}
