package models

import "time"

func Timestamp() int64 {
	return FromTime(time.Now())
}

func FromTime(t time.Time) int64 {
	return t.UTC().UnixNano() / int64(time.Millisecond)
}
