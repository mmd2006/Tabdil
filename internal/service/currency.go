package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type DollarRateResponse struct {
	USD float64 `json:"usd"`
}

func FetchDollarRate(apiURL string, apiKey string) (float64, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to fetch dollar rate")
	}

	var data DollarRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	if data.USD <= 0 {
		return 0, errors.New("invalid dollar rate received")
	}

	return data.USD, nil
}

func TomanToDollar(toman float64, usdRate float64) (float64, error) {
	if toman < 0 {
		return 0, errors.New("toman cannot be negative")
	}
	if usdRate <= 0 {
		return 0, errors.New("usd rate must be greater than zero")
	}

	return toman / usdRate, nil
}

func ConvertTomanToDollarWithAPI(
	toman float64,
	apiURL string,
	apiKey string,
) (float64, error) {

	rate, err := FetchDollarRate(apiURL, apiKey)
	if err != nil {
		return 0, err
	}

	return TomanToDollar(toman, rate)
}
