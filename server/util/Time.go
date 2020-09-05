package util

import (
	"time"
)

func ParseTimeDuration(text string) (time.Duration, error) {
	return time.ParseDuration(text)
}
