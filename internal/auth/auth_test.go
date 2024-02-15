package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExchangeToken(t *testing.T) {
	// Given
	code := "my_auth_token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("code") != code {
			t.Error("code not in query params")
		}
		if r.URL.Query().Get("grant_type") != "authorization_code" {
			t.Error("grant type not correct")
		}

		responseBody := mockResponse()
		repsonseBytes, _ := json.Marshal(responseBody)
		w.Write(repsonseBytes)
	}))
	defer server.Close()

	client := NewStravaClient(server.URL)

	// When
	authResponse, err := client.ExchangeToken(code)

	// Then
	if err != nil {
		t.Errorf("error calling client: %v", err)
	}
	if authResponse.Access_Token != "Access" {
		t.Error("response token null")
	}
}

func TestRefreshToken(t *testing.T) {
	// Given
	code := "my_auth_token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("code") != code {
			t.Error("code not in query params")
		}
		if r.URL.Query().Get("grant_type") != "refresh_token" {
			t.Error("grant type not correct")
		}

		responseBody := mockResponse()
		repsonseBytes, _ := json.Marshal(responseBody)
		w.Write(repsonseBytes)
	}))
	defer server.Close()

	client := NewStravaClient(server.URL)

	// When
	authResponse, err := client.RefreshToken(code)

	// Then
	if err != nil {
		t.Errorf("error calling client: %v", err)
	}
	if authResponse.Access_Token != "Access" {
		t.Error("response token null")
	}
}

func mockResponse() AuthorizationResponse {
	return AuthorizationResponse{
		Token_Type:    "Access",
		Expires_At:    123,
		Expires_In:    123,
		Refresh_Token: "Refresh",
		Access_Token:  "Access",
		Athlete: Athlete{
			UserName: "user",
			ID:       1234,
		},
	}
}
