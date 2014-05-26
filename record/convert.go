package record

import (
	"strconv"
	"time"
)

func StrToBool(val string, defaultVal bool) bool {

	if val, err := strconv.ParseBool(val); err == nil {
		return val
	}
	return defaultVal
}

func StrToTime(t string) time.Time {
	if val, err := time.Parse(USER_TIME_STR, t); err == nil {
		return val
	}

	return time.Unix(0, 0)

}

func StrToInt(t string) int {
	if val, err := strconv.ParseInt(t, 10, 32); err == nil {
		if val >= 0 {
			return int(val)
		}
	}
	return 0
}
