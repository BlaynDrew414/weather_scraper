package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	apiKey      = "9p3ELusMaMlbU4fJF79K97c4iQjp7Zq6"
	locationURL = "http://dataservice.accuweather.com/locations/v1/cities/search"
	currentURL  = "http://dataservice.accuweather.com/currentconditions/v1/%s?apikey=%s&details=true&imperial=true"
)

type Location struct {
	Key                string             `json:"Key"`
	LocalizedName      string             `json:"LocalizedName"`
	AdministrativeArea AdministrativeArea `json:"AdministrativeArea"`
}

type AdministrativeArea struct {
	ID            string `json:"ID"`
	LocalizedName string `json:"LocalizedName"`
}

type CurrentConditions struct {
    Temperature struct {
        Metric struct {
            Value float32 `json:"Value"`
        } `json:"Metric"`
        Imperial struct {
            Value float32 `json:"Value"`
        } `json:"Imperial"`
    } `json:"Temperature"`
    WeatherText string `json:"WeatherText"`
}


func main() {
	http.HandleFunc("/weather", handleWeatherRequest)
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)
}

func isValidInput(city, state string) bool {
	// Ensure that city and state are not empty
	if strings.TrimSpace(city) == "" || strings.TrimSpace(state) == "" {
		return false
	}

	// Ensure that city and state contain only alphabetic characters and spaces
	if !regexp.MustCompile(`^[A-Za-z ]+$`).MatchString(city) || !regexp.MustCompile(`^[A-Za-z ]+$`).MatchString(state) {
		return false
	}

	return true
}

func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		City  string `json:"city"`
		State string `json:"state"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.City == "" || requestBody.State == "" {
		http.Error(w, "Invalid input. Please enter a valid city and state.", http.StatusBadRequest)
		return
	}

	if !isValidInput(requestBody.City, requestBody.State) {
		http.Error(w, "Invalid input. Please enter a valid city and state.", http.StatusBadRequest)
		return
	}

	locationKey, err := getLocationKey(requestBody.City, requestBody.State)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	weather, err := getCurrentConditions(locationKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"city":      requestBody.City,
		"state":     requestBody.State,
		"temp":      weather.Temperature.Imperial.Value,
		"condition": weather.WeatherText,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}



func getLocationKey(city, state string) (string, error) {
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
        return "", fmt.Errorf("failed to decode JSON response: %v", err)
    }

    if len(locationResp) < 1 {
        return "", fmt.Errorf("no city/state found for %s, %s", city, state)
    }

    return locationResp[0].Key, nil
}

func getCurrentConditions(locationKey string) (*CurrentConditions, error) {
    url := fmt.Sprintf(currentURL, locationKey, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("HTTP request failed: %v", err)
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get current conditions: %s", resp.Status)
    }

    var weatherResp []*CurrentConditions

    if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
        return nil, fmt.Errorf("failed to decode JSON response: %v", err)
    }

    if len(weatherResp) < 1 {
        return nil, fmt.Errorf("no current conditions found for location key %s", locationKey)
    }

    return weatherResp[0], nil
}
