package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
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
		var userRequestBody UserResource
		err := json.Unmarshal([]byte(request.Body), &userRequestBody)
		if err != nil {
			return UserResource{}, fmt.Errorf("request body %+v could not be parsed", request.Body)
		}
		return HandleUserUpdate(context, userId, userRequestBody)
	case "DELETE":
		log.Printf("DELETE request with context: %+v", request.RequestContext)
		return UserResource{}, DeleteUser(context, userId)
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
	db, err := users.GetDbConnection(context)
	if err != nil {
		panic(err)
	}
	user, err := db.GetUser(context, userId)
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

func DeleteUser(context context.Context, userId int) error {
	db, err := users.GetDbConnection(context)
	if err != nil {
		return err
	}
	return db.DeleteUser(context, userId)
}

func HandleUserUpdate(context context.Context, userId int, userDetails UserResource) (UserResource, error) {
	users.GetDbConnection(context)
	updatedUser, err := HandleUserSubmitDetails(context, userId, userDetails.HomeLat, userDetails.HomeLng, userDetails.WorkLat, userDetails.WorkLng)
	return updatedUser, err
}

func HandleUserSubmitDetails(ctx context.Context, id int, hlat float64, hlng float64, wlat float64, wlng float64) (UserResource, error) {
	db, err := users.GetDbConnection(ctx)
	if err != nil {
		panic(err)
	}
	user, err := db.GetUser(ctx, id)
	if err != nil {
		return UserResource{}, err
	}

	user.HomeLat = hlat
	user.HomeLng = hlng
	user.WorkLat = wlat
	user.WorkLng = wlng

	err = db.UpdateUser(ctx, user)
	if err != nil {
		return UserResource{}, err
	}

	log.Printf("updated user: %v, %v, %v, %v, %v", user.ID, hlat, hlng, wlat, wlng)
	return UserResource{
		HomeLat: hlat,
		HomeLng: hlng,
		WorkLat: wlat,
		WorkLng: wlng,
	}, nil
}
