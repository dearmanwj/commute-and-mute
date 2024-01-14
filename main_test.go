package main

import (
	"log"
	"math"
	"testing"
)

func TestIsWithinRadius(t *testing.T) {
	withinRadius := isWithinRadius(43.593583, 1.448228, 43.564088, 1.389139)
	log.Println(withinRadius)
	if withinRadius {
		t.Fail()
	}
}

func TestCalculateDistance(t *testing.T) {
	distance := hsDist(1.2232, 0.1234, 1.2245, 0.1235)
	log.Println(distance)
	if math.Round(distance*100)/100 != 8.290 {
		t.Fail()
	}
}
