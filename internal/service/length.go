package service

import "errors"

func KilometerToMile(km float64) (float64, error) {
	if km < 0 {
		return 0, errors.New("kilometer cannot be negative")
	}

	mile := km * 0.621371
	return mile, nil
}
