package auth

import (
	"log"
	"testing"
)

func TestGenerateUnsigned(t *testing.T) {
	// Given
	id := 123
	utils := KmsUtils{}

	// When
	token := utils.sign(id)

	// Then
	log.Printf("token: %v", token)
}
