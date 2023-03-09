package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"example.com/weather-scraper/validator"
	"example.com/weather-scraper/weather"
)

func HandleWeatherRequest(w http.ResponseWriter, r *http.Request) {
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

	if !validator.IsValidInput(requestBody.City, requestBody.State) {
		http.Error(w, "Invalid input. Please enter a valid city and state.", http.StatusBadRequest)
		return
	}

	locationKey, err := weather.GetLocationKey(requestBody.City, requestBody.State)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	weather, err := weather.GetCurrentConditions(locationKey)
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