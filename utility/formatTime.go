package utility

import (
	"time"
)

func FormatTimeToString(duration time.Duration) string {
	layoutTime := "15:04:05"
	msec := time.UnixMilli(int64(duration.Seconds() * 1000))
	return msec.Format(layoutTime)
}
