package dto

import "errors"

func ValidateShortCode(shortCode string) error {
	if len(shortCode) != 7 {
		return errors.New("shortCode length must be 7")
	}
	return nil
}
