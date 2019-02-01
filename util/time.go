package util

import (
	"fmt"
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

func ConvertDuration(dur time.Duration) string {
	h := dur / time.Hour
	dur -= h * time.Hour
    m := dur / time.Minute
    dur -= m * time.Minute
    s := dur / time.Second
    dur -= s * time.Second
    ms := dur / time.Millisecond
    return fmt.Sprintf("%02d:%02d:%02d.%03d",h,m,s,ms)
}
