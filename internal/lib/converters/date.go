package converters

import (
	"time"

	"google.golang.org/genproto/googleapis/type/date"
)

func ToProtoDate(t time.Time) *date.Date {
	return &date.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}
}

func FromProtoDate(t *date.Date) time.Time {
	return time.Date(int(t.Year), time.Month(t.Month), int(t.Day), 0, 0, 0, 0, time.UTC)
}
