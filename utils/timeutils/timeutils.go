package timeutils

import "time"

func ConvertUnixTimestampToDate(timestamp int64, layout string) string {
	return time.Unix(timestamp, 0).Format(layout)
}
