package auth

import (
	"log"
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	// Given
	id := 1234

	// When
	result := GenerateUserToken(id)

	// Then
	token, err := jwt.Parse(result, func(token *jwt.Token) (interface{}, error) {
		return GetPublicKey(), nil
	})
	if err != nil {
		log.Print("error parsing token")
		t.Fail()
	}
	sub, _ := token.Claims.GetSubject()
	if sub != strconv.FormatInt(int64(id), 10) {
		log.Print("id does not match")
		t.Fail()
	}
	if !token.Valid {
		log.Print("token not valid")
		t.Fail()
	}
}

func TestValidate(t *testing.T) {
	// Given
	id := 1234

	// When
	result := "eyJhbGciOiJFUzI1NiIsInR5cCI6Imp3dCJ9.eyJzdWIiOjEwNTAzODEyLCJpc3MiOiJjb21tdXRlLWFuZC1tdXRlIiwiZXhwIjoxNzExOTA4MjA1fQ.MEUCIDvhm4jJNUQNf5_K-Uhes-ibRPW8kHQtwUsjzvkE537tAiEAgpzAqYyvdlhckm8k5e7dBa1cH02Lfpl2SuzIFc4-Nj4"

	// Then
	token, err := jwt.Parse(result, func(token *jwt.Token) (interface{}, error) {
		return GetPublicKey(), nil
	})
	if err != nil {
		log.Print("error parsing token")
		t.Fail()
	}
	sub, _ := token.Claims.GetSubject()
	if sub != strconv.FormatInt(int64(id), 10) {
		log.Print("id does not match")
		t.Fail()
	}
	if !token.Valid {
		log.Print("token not valid")
		t.Fail()
	}
}
