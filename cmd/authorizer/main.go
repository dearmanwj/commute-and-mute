package authorizer

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func main() {

}

func handleAuth(context context.Context, request events.APIGatewayCustomAuthorizerRequest) events.APIGatewayCustomAuthorizerResponse {

	token := request.AuthorizationToken

}
