package router

import (
	"github.com/gorilla/mux"

	"jh-weather-api/handlers"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/getWeather/{latitude}/{longitude}", handlers.GetWeatherHandler).Methods("GET")

	return r
}
