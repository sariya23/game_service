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
