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

const cmToFoot = 0.0328084

func CentimeterToFoot(cm float64) (float64, error) {
	if cm < 0 {
		return 0, errors.New("centimeter cannot be negative")
	}

	return cm * cmToFoot, nil
}

const megabitToMegabyte = 0.125

func MegabitToMegabyte(mb float64) (float64, error) {
	if mb < 0 {
		return 0, errors.New("megabit cannot be negative")
	}

	return mb * megabitToMegabyte, nil
}
