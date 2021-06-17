package util

import (
	"strconv"
	"time"

	"github.com/mundipagg/boleto-api/log"
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
		lg := log.CreateLog()
		lg.Warn(err.Error(), "Could not get Timezone")
		return time.Now()
	}
	t := time.Now()
	local := t.In(z)
	return local
}

func NycNow() time.Time {
	z, err := time.LoadLocation("America/New_York")
	if err != nil {
		lg := log.CreateLog()
		lg.Warn(err.Error(), "Could not get Timezone")
		return time.Now()
	}
	t := time.Now()
	local := t.In(z)
	return local
}

func GetDurationTimeoutRequest(t string) time.Duration {
	tTime, _ := strconv.Atoi(t)
	tOut := time.Duration(tTime)
	return tOut
}
