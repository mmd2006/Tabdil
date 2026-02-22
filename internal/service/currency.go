package service

import "errors"

func TomanToDollar(toman float64, usdRate float64) (float64, error) {
	if toman < 0 {
		return 0, errors.New("toman cannot be negative")
	}
	if usdRate <= 0 {
		return 0, errors.New("usd rate must be greater than zero")
	}

	return toman / usdRate, nil
}
