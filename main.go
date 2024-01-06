package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("Hello world")
	response := events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       "This is my function",
	}
	return response, nil
}
