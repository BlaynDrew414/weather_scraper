package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"example.com/weather-scraper/models"
)

const (
	apiKey      = "APIKey"
	locationURL = "http://dataservice.accuweather.com/locations/v1/cities/search"
	currentURL  = "http://dataservice.accuweather.com/currentconditions/v1/%s?apikey=%s&details=true&imperial=true"
)

func GetLocationKey(city, state string) (string, error) {
	location := fmt.Sprintf("%s, %s", strings.TrimSpace(city), strings.TrimSpace(state))
	query := url.Values{
		"apikey": {apiKey},
		"q":      {location},
	}

	resp, err := http.Get(fmt.Sprintf("%s?%s", locationURL, query.Encode()))
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get location key: %s", resp.Status)
	}

	var locationResp []struct {
		Key string `json:"Key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&locationResp); err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(locationResp) < 1 {
		return "", fmt.Errorf("no city/state found for %s, %s", city, state)
	}

	return locationResp[0].Key, nil
}

func GetCurrentConditions(locationKey string) (models.CurrentConditions, error) {
	url := fmt.Sprintf(currentURL, locationKey, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return models.CurrentConditions{}, fmt.Errorf("HTTP request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.CurrentConditions{}, fmt.Errorf("failed to get current conditions: %s", resp.Status)
	}

	var weatherResp []models.CurrentConditions

	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return models.CurrentConditions{}, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(weatherResp) < 1 {
		return models.CurrentConditions{}, fmt.Errorf("no current conditions found for location key %s", locationKey)
	}

	return weatherResp[0], nil
}

