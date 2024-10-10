package validators

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func DateValidation(fl validator.FieldLevel) bool {
	releaseDateStr := fl.Field().String()

	if releaseDateStr == "" {
		return true
	}

	releaseDate, err := time.Parse("2006-01-02", releaseDateStr)
	if err != nil {
		return false
	}

	lowerBound := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	upperBound := time.Now()

	return releaseDate.After(lowerBound) && releaseDate.Before(upperBound.Add(24*time.Hour))
}
