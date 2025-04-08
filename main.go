package main

import (
	"fmt"
	"log"
	"net/http"

	"jh-weather-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/getWeather/{latitude}/{longitude}", handlers.GetWeatherHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))

	fmt.Println("Server is running on http://localhost:8080")
}
