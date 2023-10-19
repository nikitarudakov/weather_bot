package timeutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTimeFormat(t *testing.T) {
	testCases := []struct {
		timeStr string
		errExp  bool
	}{
		{"23:11", false},
		{"11:11PM", true},
		{"08:05", false},
		{"29:21", true},
		{"13:23", false},
		{"24:00", true},
		{"00:00", false},
		{"05:05", false},
		{"asd", true},
		{"03:04AM", true},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.timeStr, func(t *testing.T) {
			t.Parallel()

			parsedTime, err := ParseTimeFormat(tc.timeStr)
			if err != nil {
				if tc.errExp {
					assert.Error(t, err)
				} else {
					t.Error(err)
				}
			}

			t.Log("Parsed time", parsedTime)
		})
	}
}
