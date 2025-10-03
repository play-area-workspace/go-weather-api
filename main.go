package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Temperature struct {
	Temp float64 `json:"temp"`
}

type WeatherResponse struct {
	Main Temperature `json:"main"`
}

func fetchWeatherData(lat, lon, apiKey string, logger *slog.Logger) (*Temperature, error) {

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", lat, lon, apiKey)

	weatherResponse := WeatherResponse{}
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Error making HTTP request", "error", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		logger.Error("Error decoding JSON response", "error", err.Error())
		return nil, err
	}
	return &weatherResponse.Main, nil
}

func main() {

	startNow := time.Now()
	defer func() {
		fmt.Printf("Execution time: %s\n", time.Since(startNow))
	}()

	var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file", "error", err.Error())
	} else {
		logger.Info(".env file loaded successfully")
	}

	logger.Info("Application started")
	defer logger.Info("Application stopped")

	apiKey := os.Getenv("API_KEY")

	cities := make(map[string][2]string)
	// Latitude, Longitude
	cities["Colombo"] = [2]string{"6.9271", "79.8612"}
	cities["New York"] = [2]string{"40.7128", "-74.0060"}
	cities["London"] = [2]string{"51.5074", "-0.1278"}
	cities["Tokyo"] = [2]string{"35.6895", "139.6917"}

	for cityName, city := range cities {
		temperature, err := fetchWeatherData(city[0], city[1], apiKey, logger)
		if err != nil {
			logger.Error("Failed to fetch weather data", "error", err.Error(), "city", cityName)
			continue
		}

		tempC := temperature.Temp - 273.15
		logger.Info("Fetched weather data successfully", "city", cityName,
			"temperature(°C)", round(tempC, 2, "round"))
	}
}

// Round temperature to 2 decimal points, using ceil or floor as needed
func round(val float64, decimals int, mode string) float64 {
	mult := math.Pow(10, float64(decimals))
	switch mode {
	case "ceil":
		return math.Ceil(val*mult) / mult
	case "floor":
		return math.Floor(val*mult) / mult
	default: // "round"
		return math.Round(val*mult) / mult
	}
}
