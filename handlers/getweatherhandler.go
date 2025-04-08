package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"jh-weather-api/models"

	"github.com/gorilla/mux"
)

func GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	latitude := vars["latitude"]
	longitude := vars["longitude"]

	latitudeAsFloat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		fmt.Println("An error occurred: ", err.Error())
	}

	longitudeAsFloat, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		fmt.Println("An error occurred")
	}

	resp, err := getWeather(latitudeAsFloat, longitudeAsFloat)
	if err != nil {
		fmt.Println("An error occurred")
	}

	fmt.Fprintf(w, "%s\n", resp)
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

	var pointsData models.PointsResponse
	json.Unmarshal(body, &pointsData)

	forecastURL := pointsData.Properties.Forecast
	return forecastURL
}

// getForecastData returns the forecast data for the given latitude &
// longitude.
func getForecastData(forecastURL string) models.ForecastResponse {
	req2, _ := http.NewRequest("GET", forecastURL, nil)
	req2.Header.Set("User-Agent", "your-email@example.com")

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)

	var forecastData models.ForecastResponse
	json.Unmarshal(body2, &forecastData)

	return forecastData
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
