package strava

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
		if r.URL.Query().Get("grant_type") != "code" {
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

func TestExchangeTokenEmptyCode(t *testing.T) {
	// Given
	var code string
	params := map[string]string{}
	code = params["code"]

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer server.Close()

	client := NewStravaClient(server.URL)

	// When
	_, err := client.ExchangeToken(code)

	// Then
	if err == nil {
		t.Errorf("error calling client: %v", err)
	}
}

func TestRefreshToken(t *testing.T) {
	// Given
	code := "my_auth_token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("refresh_token") != code {
			t.Error("code not in query params")
		}
		if r.URL.Query().Get("grant_type") != "refresh_token" {
			t.Error("grant type not correct")
		}

		responseBody := mockResponse()
		responseBytes, _ := json.Marshal(responseBody)
		w.Write(responseBytes)
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

func TestUpdateActivity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/api/v3/activities/10583609809" {
			t.Errorf("incorrect url: %v", r.URL.String())
		}
		if r.Method != http.MethodPut {
			t.Error("correct url called with wrong method")
		}
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Error("token not used")
		}

		var updateBody ActivityUpdate
		err := json.NewDecoder(r.Body).Decode(&updateBody)
		if err != nil {
			t.Errorf("error parsing request body: %v", err)
		}
		if (updateBody != ActivityUpdate{Commute: true, Hide_From_Home: true}) {
			t.Errorf("incorrect update request sent: %v", updateBody)
		}
	}))
	defer server.Close()

	stravaClient := NewStravaClient(server.URL)

	// When
	err := stravaClient.UpdateActivity(10583609809, "token")

	if err != nil {
		t.Errorf("request failed: %v", err)
	}

}

func TestGetActivityOK(t *testing.T) {
	// Given
	toReturn := Activity{
		Id:           10583609809,
		Type:         "run",
		Start_latlng: [2]float64{12.34, 56.78},
		End_latlng:   [2]float64{12.34, 56.78},
		Athlete:      Athlete{UserName: "user", ID: 12345},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/api/v3/activities/10583609809" {
			t.Errorf("incorrect url: %v", r.URL.String())
		}
		if r.Method != http.MethodGet {
			t.Error("correct url called with wrong method")
		}
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Error("token not used")
		}

		json.NewEncoder(w).Encode(toReturn)
	}))
	defer server.Close()

	stravaClient := NewStravaClient(server.URL)

	// When
	result, err := stravaClient.GetActivity(10583609809, "token")

	// Then
	if err != nil {
		t.Error("error getting activity")
	}
	if result != toReturn {
		t.Error("incorrect activity returned")
	}
}

func TestGetActivityHandleUnauthorized(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	stravaClient := NewStravaClient(server.URL)

	// When
	_, err := stravaClient.GetActivity(10583609809, "token")

	// Then
	if err == nil {
		t.Error("did not pass on unauthorized message")
	}
}
