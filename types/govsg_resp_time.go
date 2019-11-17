package types

import (
	"strings"
	"time"
)

type GovSgResponseTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05"

func (ct *GovSgResponseTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	return
}
