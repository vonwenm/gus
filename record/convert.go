package record

import (
	"strconv"
	"strings"
	"time"
)

func StrToBool(val string, defaultVal bool) bool {

	if val, err := strconv.ParseBool(val); err == nil {
		return val
	}
	val = strings.TrimSpace(strings.ToLower(val))
	if val == `yes` || val == `ok` {
		return true
	}
	if val == `no` {
		return false
	}
	return defaultVal
}

func StrToTime(t string) time.Time {
	if val, err := time.Parse(USER_TIME_STR, t); err == nil {
		return val
	}

	return time.Time{}

}

func StrToInt(t string) int {
	if val, err := strconv.Atoi(t); err == nil {
		return val
	}
	return 0
}
