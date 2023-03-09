package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var testLocationURL = "http://dataservice.accuweather.com/locations/v1/cities/search"


func TestGetLocationKey(t *testing.T) {
	// Create a fake HTTP server to mimic the AccuWeather API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"Key": "12345",
				"LocalizedName": "San Francisco",
				"AdministrativeArea": {
					"ID": "CA",
					"LocalizedName": "California"
				}
			}
		]`))
	}))
	defer server.Close()

	// Override the AccuWeather API URL with the fake server's URL
	testLocationURL = server.URL

	// Test with valid city and state inputs
	city := "San Francisco"
	state := "CA"
	locationKey, err := GetLocationKey(city, state)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if locationKey != "12345" {
		t.Errorf("Expected locationKey to be %s, but got %s", "12345", locationKey)
	}

	// Test with invalid city and state inputs
	city = "Invalid City"
	state = "Invalid State"
	locationKey, err = GetLocationKey(city, state)
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
	if locationKey != "" {
		t.Errorf("Expected locationKey to be empty, but got %s", locationKey)
	}
}
