package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"willd/commute-and-mute/internal/auth"
	"willd/commute-and-mute/internal/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type UserFormData struct {
	HomeLat float64 `json:"hlat"`
	HomeLng float64 `json:"hlng"`
	WorkLat float64 `json:"wlat"`
	WorkLng float64 `json:"wlng"`
}

func main() {
	lambda.Start(HandleUserUpdate)
}

func HandleUserUpdate(context context.Context, request *events.LambdaFunctionURLRequest) error {
	users.GetDbConnection()
	log.Println("Updating user with home and work locations")
	authHeader := request.Headers["Authorization"]
	token := strings.Split(authHeader, "Bearer ")[1]
	id, err := auth.GetConnectedUserId(token)
	if err != nil {
		log.Panicln(err)
	}
	var formData UserFormData
	json.Unmarshal([]byte(request.Body), &formData)
	HandleUserSubmitDetails(id, formData.HomeLat, formData.HomeLng, formData.WorkLat, formData.WorkLng)
	return nil
}

func HandleUserSubmitDetails(id int, hlat float64, hlng float64, wlat float64, wlng float64) error {

	user, err := users.GetUser(id)
	if err != nil {
		return err
	}

	user.HomeLat = hlat
	user.HomeLng = hlng
	user.WorkLat = wlat
	user.WorkLng = wlng

	err = users.UpdateUser(user)
	if err != nil {
		return err
	}

	log.Printf("updated user: %v, %v, %v, %v, %v", user.ID, hlat, hlng, wlat, wlng)
	return nil
}
