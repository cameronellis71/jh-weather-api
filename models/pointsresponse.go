package models

// Struct to parse JSON response from /points endpoint
type PointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}
