package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetWeatherHandlerReturnsNonEmptyStringValidInput(t *testing.T) {
	// Setup the request
	req, err := http.NewRequest("GET", "/weather/40.7128/-74.0060", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and register the handler
	r := mux.NewRouter()
	r.HandleFunc("/weather/{latitude}/{longitude}", GetWeatherHandler)

	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	// Check if the status code is what you expect
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 but got %v", rr.Code)
	}

	// Check if the response body is not empty
	if rr.Body.String() == "" {
		t.Errorf("expected body %v but got %v", "", rr.Body.String())
	}
}

func TestGetWeatherHandlerBadLatitudeValue(t *testing.T) {
	// Setup the request
	req, err := http.NewRequest("GET", "/weather/40.7128/abcde", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and register the handler
	r := mux.NewRouter()
	r.HandleFunc("/weather/{latitude}/{longitude}", GetWeatherHandler)

	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	// Check if the status code is what you expect
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 200 but got %v", rr.Code)
	}

}

func TestGetWeatherHandlerBadLongitudeValue(t *testing.T) {
	// Setup the request
	req, err := http.NewRequest("GET", "/weather/abcde/-74.0060", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new router and register the handler
	r := mux.NewRouter()
	r.HandleFunc("/weather/{latitude}/{longitude}", GetWeatherHandler)

	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	// Check if the status code is what you expect
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 200 but got %v", rr.Code)
	}

}

func TestGetForecastUrlHappyPath(t *testing.T) {
	// Initialize the mock HTTP client
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register a mock responder for the weather API endpoint
	httpmock.RegisterResponder("GET", "https://api.weather.gov/points/40.712800,-74.006000",
		httpmock.NewStringResponder(200, `{
            "properties": {
                "forecast": "https://api.weather.gov/gridpoints/XYZ/123,456/forecast"
            }
        }`))

	// Call the function to test
	forecastURL := getForecastUrl(40.712800, -74.006000)

	// Check that the URL returned matches the mock response
	assert.Equal(t, "https://api.weather.gov/gridpoints/XYZ/123,456/forecast", forecastURL)
}

func TestGetForecastData(t *testing.T) {
	// Initialize httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Set up the mock response for the forecast URL
	httpmock.RegisterResponder("GET", "https://api.weather.gov/gridpoints/ABC/123,456/forecast",
		httpmock.NewStringResponder(200, `{
			"properties": {
				"periods": [
					{
						"shortForecast": "Partly Cloudy",
						"temperature": 72
					}
				]
			}
		}`))

	// Call the function being tested
	forecastData := getForecastData("https://api.weather.gov/gridpoints/ABC/123,456/forecast")

	// Check that the forecast data is correct
	assert.Equal(t, "Partly Cloudy", forecastData.Properties.Periods[0].ShortForecast)
	assert.Equal(t, 72, forecastData.Properties.Periods[0].Temperature)
}

func TestGetTempCharicaterization(t *testing.T) {
	tests := []struct {
		temperature int
		expected    string
	}{
		{50, "cold"},
		{70, "moderate"},
		{90, "hot"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Temperature %d", test.temperature), func(t *testing.T) {
			result := getTempCharicaterization(test.temperature)
			if result != test.expected {
				t.Errorf("expected %s but got %s", test.expected, result)
			}
		})
	}
}

func TestGetWeather_EmptyPeriods(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Match any /points/{lat},{lon} URL
	httpmock.RegisterResponder("GET", `=~^https://api.weather.gov/points/.*`,
		httpmock.NewStringResponder(200, `{
			"properties": {
				"forecast": "https://api.weather.gov/gridpoints/XYZ/123,456/forecast"
			}
		}`))

	// Mock the forecast response with empty periods
	httpmock.RegisterResponder("GET", "https://api.weather.gov/gridpoints/XYZ/123,456/forecast",
		httpmock.NewStringResponder(200, `{
			"properties": {
				"periods": []
			}
		}`))

	result, err := getWeather(40.7128, -74.0060)

	assert.NotNil(t, err)
	assert.Equal(t, "", result)
	assert.Equal(t, "no forecasts available to display", err.Error())
}
