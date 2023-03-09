package main

import (
	
	"net/http"
	"example.com/weather-scraper/handlers"
)




func main() {
	http.HandleFunc("/weather", handlers.HandleWeatherRequest)
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)
}

