package timeutils

import "time"

// ConvertUnixTimestampToDate converts Unix Timestamp to date format (layout)
func ConvertUnixTimestampToDate(timestamp int64, layout string) string {
	return time.Unix(timestamp, 0).Format(layout)
}
