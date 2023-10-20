package timeutils

import "time"

const Layout24H = "15:04"

// ConvertUnixTimestampToDate converts Unix Timestamp to date format (layout)
func ConvertUnixTimestampToDate(timestamp int64, layout string) string {
	return time.Unix(timestamp, 0).Format(layout)
}

func ParseTimeFormat(ts string) (*time.Time, error) {
	t, err := time.Parse(Layout24H, ts)
	if err != nil {
		return nil, err
	}

	tUTC := t.UTC()

	return &tUTC, nil
}
