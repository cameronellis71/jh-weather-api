package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Home Page!\n")
}

func getWeatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	latitude := vars["latitude"]
	longitude := vars["longitude"]
	fmt.Fprintf(w, "Latitude %s!\n", latitude)
	fmt.Fprintf(w, "Longitude %s!\n", longitude)
}

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Define your routes and associate them with handler functions
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/getWeather/{latitude}/{longitude}", getWeatherHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))

	fmt.Println("Server is running on http://localhost:8080")
}

// getForecastUrl returns the URL from which to get the weather data.
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

// getForecastData returns the forecast data for the given latitude &
// longitude.
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

// getTemperatureCharicaterization returns the charicaterization of the
// temperature.
func getTempCharicaterization(temperature int) string {
	if temperature < 60 {
		return "cold"
	} else if temperature >= 60 && temperature <= 80 {
		return "moderate"
	}
	return "hot"
}

// getWeather returns the short forecast and charicaterizaion of the weather.
func getWeather(latitude, longitude float64) (string, error) {
	// Get the forecast URL from /points/{lat},{lon}
	forecastURL := getForecastUrl(latitude, longitude)

	// Get the forecast data
	forecastData := getForecastData(forecastURL)

	periods := forecastData.Properties.Periods
	if len(periods) == 0 {
		return "", errors.New("no forecasts available to display")
	}

	period := periods[0]
	shortForecast := period.ShortForecast
	tempCharicaterization := getTempCharicaterization(period.Temperature)

	return shortForecast + ": " + tempCharicaterization, nil
}
