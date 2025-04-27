package main

import (
	"log"
	"math"
	"strconv"
)

const RadiusKM float64 = 0.1

func IsWithinRadius(lat1, lng1, lat2, lng2 string) bool {
	distKm := hsDist(
		degToRad(lat1),
		degToRad(lng1),
		degToRad(lat2),
		degToRad(lng2),
	)
	log.Println(distKm)
	return distKm < RadiusKM
}

func degToRad(deg string) float64 {
	degFloat, err := strconv.ParseFloat(deg, 64)
	if err != nil {
		log.Panicf("error parsing db info %v", err)
	}
	return degFloat * math.Pi / 180
}

func haversine(theta float64) float64 {
	return .5 * (1 - math.Cos(theta))
}

const rEarth = 6371 // km

func hsDist(phi1, psi1, phi2, psi2 float64) float64 {
	return 2 * rEarth * math.Asin(math.Sqrt(haversine(phi2-phi1)+
		math.Cos(phi1)*math.Cos(phi2)*haversine(psi2-psi1)))
}
