package utils

import (
	"time"
)

type JTime time.Time

const (
	timeFormart = "2006-01-02 15:04:05.000"
)

func (t *JTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = JTime(now)
	return
}

func (t JTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t JTime) String() string {
	return time.Time(t).Format(timeFormart)
}
