package util

import (
	"time"
)

func Duration(callback func()) (duration time.Duration) {
	start := time.Now()
	callback()
	end := time.Now()
	duration = end.Sub(start)
	return
}

func BrNow() time.Time {
	z, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.Now()
	}
	t := time.Now()
	local := t.In(z)
	return local
}

func NycNow() time.Time {
	z, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.Now()
	}
	t := time.Now()
	local := t.In(z)
	return local
}
