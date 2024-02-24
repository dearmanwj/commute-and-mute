package main

import "testing"

func TestDecodeEvent(t *testing.T) {
	// Given
	rawEvent := `{
		"aspect_type": "update",
		"event_time": 1516126040,
		"object_id": 1360128428,
		"object_type": "activity",
		"owner_id": 134815,
		"subscription_id": 120475,
		"updates": {
			"title": "Messy"
		}
	}`

	// When
	result, err := DecodeUpdateEvent(rawEvent)
	if err != nil {
		t.Errorf("failed to decode event")
	}
	if result.AspectType != "update" {
		t.Errorf("did not get correct values")
	}
}
