package models

import "time"

func Timestamp() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}
