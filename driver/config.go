package driver

import (
	"strconv"
	"time"
)

func parseDuration(d string) time.Duration {
	duration, err := time.ParseDuration(d)
	if err != nil {
		panic(err)
	}
	return duration
}

func parseFpath(v string, _default string) string {
	_v := ("/var/log" + v)
	if _v != "" {
		return (_v)
	} else {
		return _default
	}
}

func parseInt(v string, _default int) int {
    _v, err := strconv.ParseInt(v, 10, 0);
    if err == nil {
        return int(_v)
    } else {
        return _default
    }
}

func readWithDefault(m map[string]string, key string, def string) string {
	value, ok := m[key]
	if ok {
		return value
	}

	return def
}
