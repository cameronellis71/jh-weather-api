package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Struct to parse JSON response from /points endpoint
type PointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// Struct to parse JSON forecast data
type ForecastResponse struct {
	Properties struct {
		Periods []struct {
			Name            string `json:"name"`
			Temperature     int    `json:"temperature"`
			TemperatureUnit string `json:"temperatureUnit"`
			ShortForecast   string `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

func main() {
	lat := 37.7749
	lon := -122.4194

	getWeather(lat, lon)
}

func getForecastUrl(latitude, longitude float64) string {
	pointsURL := fmt.Sprintf("https://api.weather.gov/points/%f,%f", latitude, longitude)

	req, _ := http.NewRequest("GET", pointsURL, nil)
	req.Header.Set("User-Agent", "cameronellis71@gmail.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var pointsData PointsResponse
	json.Unmarshal(body, &pointsData)

	forecastURL := pointsData.Properties.Forecast
	return forecastURL
}

func getForecastData(forecastURL string) ForecastResponse {
	req2, _ := http.NewRequest("GET", forecastURL, nil)
	req2.Header.Set("User-Agent", "your-email@example.com")

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)

	var forecastData ForecastResponse
	json.Unmarshal(body2, &forecastData)

	return forecastData
}

func getWeather(latitude, longitude float64) {
	// Get forecast URL from /points/{lat},{lon}
	forecastURL := getForecastUrl(latitude, longitude)

	// Get forecast data
	forecastData := getForecastData(forecastURL)
	periods := forecastData.Properties.Periods

	// Print the short forecast if one exists
	if len(periods) <= 0 {
		fmt.Println("No forecasts available to display")
	} else {
		fmt.Println(periods[0].ShortForecast)
	}
}
