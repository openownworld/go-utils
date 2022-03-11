package utils

import (
	"fmt"
	"time"
)

type JsonTime time.Time

var UseLocalTimeZone bool = true

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	var now time.Time
	if UseLocalTimeZone {
		now, err = time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	} else {
		now, err = time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.UTC)
	}
	*t = JsonTime(now)
	return
}

func (t *JsonTime) MarshalJSON() ([]byte, error) {
	//b := make([]byte, 0, len(timeFormat)+2)
	//b = append(b, '"')
	//b = time.Time(*t).AppendFormat(b, timeFormat)
	//b = append(b, '"')
	//return b, nil
	if UseLocalTimeZone {
		return []byte(fmt.Sprintf("\"%s\"", time.Time(*t).Local().Format(timeFormat))), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*t).UTC().Format(timeFormat))), nil
}

func (t *JsonTime) String() string {
	if UseLocalTimeZone {
		return time.Time(*t).Local().Format(timeFormat)
	} else {
		return time.Time(*t).UTC().Format(timeFormat)
	}
}
