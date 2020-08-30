package util

import "time"

func TimeNow() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}
