package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
)

const HomeLat float64 = 43.593583
const HomeLng float64 = 1.448228
const WorkLat float64 = 43.564060
const WorkLng float64 = 1.389155
const RadiusKM float64 = 0.1

type Activity struct {
	Id           int64
	Type         string
	Start_latlng [2]float64
	End_latlng   [2]float64
}

type ActivityUpdate struct {
	Commute      bool
	HideFromHome bool
}

func main() {

	godotenv.Load("local.env")

	log.Println(os.Getenv("AWS_ACCESS_KEY_ID"))

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1")},
	)

	svc := dynamodb.New(sess)

	http.HandleFunc("/", handlerHttp)
	http.ListenAndServe(":8080", nil)
}

func handlerHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("In http handler")
	log.Println(r.URL)

	if r.URL.String() == "/activity" {
		var a Activity
		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ProcessActivity(a)
		return
	}
}

func ProcessActivity(a Activity) (err error) {
	if strings.EqualFold(a.Type, "ride") {
		log.Println("Received Ride activity")
		if isCommute(a.Start_latlng[0], a.Start_latlng[1], a.End_latlng[0], a.End_latlng[1]) {
			log.Println("is ride and commute")
			toSend := ActivityUpdate{Commute: true, HideFromHome: true}
			fmt.Printf("To send: %v", toSend)
		} else {
			log.Println("ride not between home and work locations")
		}
	}
	return nil
}

func isCommute(startLat, startLng, endLat, endLng float64) bool {
	var isCommute bool = false
	isHomeStart := isWithinRadius(HomeLat, HomeLng, startLat, startLng)
	if isHomeStart {
		// if home is start, is commute if end is work
		isCommute = isWithinRadius(WorkLat, WorkLng, endLat, endLng)
	} else {
		isHomeEnd := isWithinRadius(HomeLat, HomeLng, endLat, endLng)
		if isHomeEnd {
			// if home is end, is commute if start is work
			isCommute = isWithinRadius(WorkLat, WorkLng, startLat, startLng)
		}
	}
	return isCommute
}

func isWithinRadius(lat1, lng1, lat2, lng2 float64) bool {
	distKm := hsDist(
		degToRad(lat1),
		degToRad(lng1),
		degToRad(lat2),
		degToRad(lng2),
	)
	log.Println(distKm)
	return distKm < RadiusKM
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func haversine(theta float64) float64 {
	return .5 * (1 - math.Cos(theta))
}

const rEarth = 6371 // km

func hsDist(phi1, psi1, phi2, psi2 float64) float64 {
	return 2 * rEarth * math.Asin(math.Sqrt(haversine(phi2-phi1)+
		math.Cos(phi1)*math.Cos(phi2)*haversine(psi2-psi1)))
}
