package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ExchangeResponse struct {
	Result string             `json:"result"`
	Base   string             `json:"base_code"`
	Rates  map[string]float64 `json:"rates"`
}

// گرفتن نرخ IRR برای هر 1 USD
func FetchDollarRate() (float64, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	url := "https://open.er-api.com/v6/latest/USD"

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("api error status: %d", resp.StatusCode)
	}

	var data ExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	if data.Result != "success" {
		return 0, errors.New("api returned unsuccessful result")
	}

	irrRate, ok := data.Rates["IRR"]
	if !ok {
		return 0, errors.New("IRR rate not found")
	}

	if irrRate <= 0 {
		return 0, errors.New("invalid IRR rate")
	}

	return irrRate, nil
}

// تبدیل تومان به دلار
func ConvertTomanToDollarWithAPI(toman float64) (float64, error) {

	if toman < 0 {
		return 0, errors.New("toman cannot be negative")
	}

	rate, err := FetchDollarRate()
	if err != nil {
		return 0, err
	}

	return toman / rate, nil
}
