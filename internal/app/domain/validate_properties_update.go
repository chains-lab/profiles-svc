package domain

import (
	"fmt"
	"time"
)

func ValidateUpdateProperty(last time.Time, duration time.Duration) error {
	now := time.Now().UTC()

	if now.Sub(last) < duration {
		return fmt.Errorf("update is not allowed to update now, "+
			"you can update it after %s from last update, "+
			"hope you are really change that and dont lie", duration,
		)
	}

	return nil
}
