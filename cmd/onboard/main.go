package main

import (
	"log"
	"willd/commute-and-mute/internal/auth"
	"willd/commute-and-mute/internal/users"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleTokenExchange)
}

func HandleTokenExchange(code string) (users.User, error) {
	users.GetDbConnection()
	auth, err := auth.ExchangeToken(code)

	if err != nil {
		return users.User{}, err
	}

	user := auth.ToUser()

	log.Printf("user: %v", user)

	users.UpdateUser(user)
	return user, nil
}
