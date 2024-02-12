package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"willd/commute-and-mute/internal/auth"
	"willd/commute-and-mute/internal/users"

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

func HandleUserUpdate(ctx context.Context, event *UserFormData) {
	users.GetDbConnection()
	log.Println("Updating user with home and work locations")
	authHeader := r.Header.Get("Authorization")
	token := strings.Split(authHeader, "Bearer ")[1]
	id, err := auth.GetConnectedUserId(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	r.ParseMultipartForm(512)
	hlat := extractFloatParameterFromReq(r.Form.Get("hlat"), w)
	hlng := extractFloatParameterFromReq(r.Form.Get("hlng"), w)
	wlat := extractFloatParameterFromReq(r.Form.Get("wlat"), w)
	wlng := extractFloatParameterFromReq(r.Form.Get("wlng"), w)
	HandleUserSubmitDetails(id, hlat, hlng, wlat, wlng)
	return
}

func extractFloatParameterFromReq(param string, w http.ResponseWriter) float64 {
	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		http.Error(w, "non-numeric lat/long provided", http.StatusBadRequest)
	}
	return val
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
