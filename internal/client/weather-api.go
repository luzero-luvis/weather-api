package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type WeatherResponse struct {
	ResolvedAddress string `json:"resolvedAddress"`
	Description     string `json:"description"`
	Timezone        string `json:"timezone"`
	Days            []Day  `json:"days"`
}

type Day struct {
	Temp       float32 `json:"temp"`
	Windspeed  float32 `json:"windspeed"`
	Humidity   float32 `json:"humidity"`
	Conditions string  `json:"Conditions"`
}

// the inputs for giving url to key and city name to func so it can get it from third party api

func GetWeather(baseUrl, apiKey, city string) (*WeatherResponse, error) {
	// removing white spaces from city when we enter in url
	escapedCity := url.QueryEscape(city)

	// complete api with url and apiKey and city name
	repsUrl := fmt.Sprintf("%s/%s?key=%s&unitGroup=metric",
		baseUrl,
		escapedCity,
		apiKey,
	)

	// loggin that we are calling weather api
	slog.Debug("calling weather api",
		"city", city,
	)

	// getting json data from the weather party api
	resp, err := http.Get(repsUrl)
	if err != nil {

		slog.Error("failed to call weather api",
			"city", city,
			"erorr", err.(*url.Error).Err,
		)

		return nil, fmt.Errorf("error fetching from api: %w", err)
	}

	// closing the thing when im done with (saves memory)

	defer resp.Body.Close()

	// if api key is wrong or url is wrong it will log that

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		slog.Warn("weather api returnd error",
			"city", city,
			"status_code", resp.StatusCode,
			"reponse", string(body),
		)

		return nil, fmt.Errorf("ERROR %d : %s", resp.StatusCode, string(body))

	}

	slog.Debug("weather api responded successfully",
		"city", city,
		"status_code", resp.StatusCode,
	)
	// we are reading all json data
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		slog.Error("failed to read body",
			"city", city,
			"erorr", err,
		)

		return nil, fmt.Errorf("erorr getting data %w", err)
	}

	// convertig json data to go data
	var weather WeatherResponse

	if err := json.Unmarshal(body, &weather); err != nil {

		slog.Error("failed to parse weather data",
			"city", city,
			"error", err,
		)
		return nil, fmt.Errorf("error convertig json to go data %w", err)
	}

	slog.Info("weather data fetched successfully",
		"city", city,
		"resolved_address", weather.ResolvedAddress,
		"days_count", len(weather.Days),
	)

	return &weather, nil
}
