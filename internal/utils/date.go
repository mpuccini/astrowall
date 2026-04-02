package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var nonNumeric = regexp.MustCompile(`[^0-9]+`)

// ParseDate parses a date string in various formats:
//   - YYYYMMDD (e.g., "20191224")
//   - YYYY-MM-DD (e.g., "2019-12-24")
//   - YYYY-M-D (e.g., "2019-7-6")
//   - Any separator between numeric groups (first 3 groups used)
func ParseDate(s string) (time.Time, error) {
	parts := nonNumeric.Split(s, -1)

	// Filter out empty strings
	var nums []string
	for _, p := range parts {
		if p != "" {
			nums = append(nums, p)
		}
	}

	if len(nums) >= 3 {
		year, err1 := strconv.Atoi(nums[0])
		month, err2 := strconv.Atoi(nums[1])
		day, err3 := strconv.Atoi(nums[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return time.Time{}, fmt.Errorf("invalid date components in %q", s)
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
	}

	// Single group: expect exactly 8 digits (YYYYMMDD)
	if len(nums) == 1 && len(nums[0]) == 8 {
		year, _ := strconv.Atoi(nums[0][:4])
		month, _ := strconv.Atoi(nums[0][4:6])
		day, _ := strconv.Atoi(nums[0][6:8])
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
	}

	return time.Time{}, fmt.Errorf("could not parse date %q: expected YYYYMMDD or YYYY-MM-DD format", s)
}
