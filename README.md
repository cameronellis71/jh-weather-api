# JH Weather API Takehome Assignment

## About
* JH Weather API is an application that takes in latitude (lat) and longitude (long) coordinates (in decimal form) & returns the short forecast for the area and a characterization of the temperature.

## Getting Started
To get a local copy instance and running follow these steps:

### Prerequisites

* Install Golang
  ```sh
  brew install go
  ```

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/github_username/repo_name.git
   ```
2. Install program dependencies
   ```sh
   go get -u github.com/gorilla/mux
   go get github.com/stretchr/testify
   go get github.com/jarcoal/httpmock
   go get github.com/stretchr/testify
   ```
3. Run the program on localhost
   ```sh
   go run main.go
   ```

<!-- USAGE EXAMPLES -->
## Usage
Type in the following into your web broser:
```
http://localhost:8080/getWeather/37.7749/-122.4195
```
Response in the browser:
```
Light Rain: cold
```
