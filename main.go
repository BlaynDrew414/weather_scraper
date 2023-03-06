package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	apiKey      = "9p3ELusMaMlbU4fJF79K97c4iQjp7Zq6"
	locationURL = "http://dataservice.accuweather.com/locations/v1/cities/search"
	currentURL  = "http://dataservice.accuweather.com/currentconditions/v1/%s?apikey=%s&details=true"
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

type currentConditions []struct {
	Temperature struct {
		Metric struct {
			Value float32 `json:"Value"`
		} `json:"Metric"`
	} `json:"Temperature"`
	WeatherText string `json:"WeatherText"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter City: ")
	scanner.Scan()
	city := scanner.Text()

	fmt.Print("Enter State: ")
	scanner.Scan()
	state := scanner.Text()

	// Check if the user input is valid
	if !isValidInput(city, state) {
		fmt.Println("Error: Invalid input. Please enter a valid city and state.")
		return
	}

	locationKey, err := getLocationKey(city, state)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	weather, err := getCurrentConditions(locationKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("The temperature in %s, %s is %.1f°C and the current weather condition is %s.\n", city, state, weather[0].Temperature.Metric.Value, weather[0].WeatherText)
}

// Check if the user input is valid
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

func getLocationKey(city, state string) (string, error) {
	location := fmt.Sprintf("%s, %s", strings.TrimSpace(city), strings.TrimSpace(state))

	query := url.Values{}
	query.Set("apikey", apiKey)
	query.Set("q", location)

	resp, err := http.Get(fmt.Sprintf("%s?%s", locationURL, query.Encode()))
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get location key: %s", resp.Status)
	}

	var locations []Location
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %v", err)
	}

	if len(locations) < 1 {
		return "", fmt.Errorf("no city/state found for %s, %s", city, state)
	}

	return locations[0].Key, nil
}

func getCurrentConditions(locationKey string) (currentConditions, error) {
	resp, err := http.Get(fmt.Sprintf(currentURL, locationKey, apiKey))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get current conditions: %s", resp.Status)
	}

	var conditions currentConditions
	err = json.NewDecoder(resp.Body).Decode(&conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return conditions, nil
}
