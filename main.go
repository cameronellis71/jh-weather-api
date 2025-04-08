package main

import (
	"fmt"
	"log"
	"net/http"

	"jh-weather-api/router"
)

func main() {
	r := router.NewRouter()

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
