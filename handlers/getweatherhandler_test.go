package handlers

import (
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
	// Your test code here
}

func TestGetTempCharicaterization(t *testing.T) {
	// Your test code here
}

func TestGetWeather(t *testing.T) {
	// Your test code here
}
