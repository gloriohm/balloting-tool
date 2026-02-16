package io

import (
	"fmt"
	"os"
	"time"
)

func validateTimestamp(f *os.File) error {
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("Kunne ikke finne timestamp på fil: %w", err)
	}

	lastMod := info.ModTime()
	useBy := time.Now().Add(-2 * time.Hour)
	if lastMod.Before(useBy) {
		return fmt.Errorf("Ballots må være lastet ned seneste for to timer siden.")
	}

	return nil
}
