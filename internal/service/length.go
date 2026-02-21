package service

import "errors"

const kgToPound = 2.20462

func KilometerToMile(km float64) (float64, error) {
	if km < 0 {
		return 0, errors.New("kilometer cannot be negative")
	}

	mile := km * 0.621371
	return mile, nil
}
func KilogramToPound(kg float64) (float64, error) {
	if kg < 0 {
		return 0, errors.New("kilogram cannot be negative")
	}

	return kg * kgToPound, nil
}
