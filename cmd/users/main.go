package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"willd/commute-and-mute/internal/auth"
	"willd/commute-and-mute/internal/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type UserResource struct {
	HomeLat float64 `json:"hlat"`
	HomeLng float64 `json:"hlng"`
	WorkLat float64 `json:"wlat"`
	WorkLng float64 `json:"wlng"`
}

func main() {
	lambda.Start(HandleUserRequest)
}

func HandleUserRequest(context context.Context, request *events.APIGatewayV2HTTPRequest) (UserResource, error) {
	method := request.RequestContext.HTTP.Method
	userId := GetUserIdFromContext(request.RequestContext)
	switch method {
	case "GET":
		log.Printf("GET request with context: %+v", request.RequestContext)
		return GetUser(context, userId)
	case "PUT":
		log.Printf("PUT request with context: %+v", request.RequestContext)
		log.Printf("PUT request with context: %+v", request.Body)
		return UserResource{
			HomeLat: 2.1,
			HomeLng: 2.2,
			WorkLat: 2.2,
			WorkLng: 2.3,
		}, nil
	default:
		return UserResource{}, errors.New("unsupported method")
	}
}

func GetUserIdFromContext(requestContext events.APIGatewayV2HTTPRequestContext) int {
	userIdString := requestContext.Authorizer.Lambda["user"].(string)
	userIdInt, err := strconv.Atoi(userIdString)
	if err != nil {
		log.Panicf("Unparseable id %v received", userIdString)
		return -1
	} else {
		return userIdInt
	}
}

func GetUser(context context.Context, userId int) (UserResource, error) {
	users.GetDbConnection(context)
	user, err := users.GetUser(context, userId)
	if err != nil {
		return UserResource{}, err
	}
	return UserResource{
		HomeLat: user.HomeLat,
		HomeLng: user.HomeLng,
		WorkLat: user.WorkLat,
		WorkLng: user.WorkLng,
	}, nil
}

func HandleUserUpdate(context context.Context, request *events.LambdaFunctionURLRequest) error {
	users.GetDbConnection(context)
	log.Println("Updating user with home and work locations")
	authHeader := request.Headers["Authorization"]
	token := strings.Split(authHeader, "Bearer ")[1]
	id, err := auth.GetConnectedUserId(token)
	if err != nil {
		log.Panicln(err)
	}
	var formData UserResource
	json.Unmarshal([]byte(request.Body), &formData)
	HandleUserSubmitDetails(context, id, formData.HomeLat, formData.HomeLng, formData.WorkLat, formData.WorkLng)
	return nil
}

func HandleUserSubmitDetails(ctx context.Context, id int, hlat float64, hlng float64, wlat float64, wlng float64) error {

	user, err := users.GetUser(ctx, id)
	if err != nil {
		return err
	}

	user.HomeLat = hlat
	user.HomeLng = hlng
	user.WorkLat = wlat
	user.WorkLng = wlng

	err = users.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	log.Printf("updated user: %v, %v, %v, %v, %v", user.ID, hlat, hlng, wlat, wlng)
	return nil
}
