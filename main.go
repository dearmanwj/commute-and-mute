package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	if len(os.Getenv("_LAMBDA_SERVER_PORT")) > 0 {
		lambda.Start(handler)
	} else {
		http.HandleFunc("/", handlerHttp)
		http.ListenAndServe(":8080", nil)
	}
}

func handler(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("Hello world")
	log.Println(request.QueryStringParameters)
	if request.RequestContext.HTTP.Method == "GET" {
		response := events.LambdaFunctionURLResponse{
			StatusCode: 200,
			Body:       request.QueryStringParameters["hub.challenge"],
		}
		return response, nil
	}
	response := events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       "Not a get request",
	}
	return response, nil
}

func handlerHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("In http handler")
	if r.Method == "GET" {
		challengeString := r.URL.Query().Get("hub.challenge")
		fmt.Fprintf(w, challengeString)
	}
}
