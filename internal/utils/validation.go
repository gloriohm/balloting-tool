package utils

import (
	"fmt"
	"os"
	"time"
)

func ValidateTimestamp(f *os.File) error {
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("kunne ikke finne timestamp på fil: %w", err)
	}

	lastMod := info.ModTime()
	useBy := time.Now().Add(-2 * time.Hour)
	if lastMod.Before(useBy) {
		return fmt.Errorf("ballots må være nyere enn to timer gamle.")
	}

	return nil
}
