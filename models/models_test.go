package models

import (
	"encoding/json"
	"testing"
)

func TestLocationJSONUnmarshal(t *testing.T) {
	jsonStr := `{
        "Key": "12345",
        "LocalizedName": "New York",
        "AdministrativeArea": {
            "ID": "NY",
            "LocalizedName": "New York"
        }
    }`

	var location Location
	if err := json.Unmarshal([]byte(jsonStr), &location); err != nil {
		t.Errorf("failed to unmarshal JSON: %v", err)
	}

	if location.Key != "12345" {
		t.Errorf("expected Key to be %q, but got %q", "12345", location.Key)
	}
	if location.LocalizedName != "New York" {
		t.Errorf("expected LocalizedName to be %q, but got %q", "New York", location.LocalizedName)
	}
	if location.AdministrativeArea.ID != "NY" {
		t.Errorf("expected AdministrativeArea.ID to be %q, but got %q", "NY", location.AdministrativeArea.ID)
	}
	if location.AdministrativeArea.LocalizedName != "New York" {
		t.Errorf("expected AdministrativeArea.LocalizedName to be %q, but got %q", "New York", location.AdministrativeArea.LocalizedName)
	}
}

func TestCurrentConditionsJSONUnmarshal(t *testing.T) {
	jsonStr := `{
        "Temperature": {
            "Metric": {
                "Value": 10
            },
            "Imperial": {
                "Value": 50
            }
        },
        "WeatherText": "Sunny"
    }`

	var conditions CurrentConditions
	if err := json.Unmarshal([]byte(jsonStr), &conditions); err != nil {
		t.Errorf("failed to unmarshal JSON: %v", err)
	}

	if conditions.Temperature.Metric.Value != 10 {
		t.Errorf("expected Metric Temperature Value to be %v, but got %v", 10, conditions.Temperature.Metric.Value)
	}
	if conditions.Temperature.Imperial.Value != 50 {
		t.Errorf("expected Imperial Temperature Value to be %v, but got %v", 50, conditions.Temperature.Imperial.Value)
	}
	if conditions.WeatherText != "Sunny" {
		t.Errorf("expected WeatherText to be %q, but got %q", "Sunny", conditions.WeatherText)
	}
}
