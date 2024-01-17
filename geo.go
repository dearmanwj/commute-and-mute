package main

import (
	"log"
	"math"
)

func IsWithinRadius(lat1, lng1, lat2, lng2 float64) bool {
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
