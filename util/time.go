package util

import (
	"fmt"
	"strconv"
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
		fmt.Println("Could not get Timezone - ", err.Error())
		return time.Now()
	}
	t := time.Now()
	local := t.In(z)
	return local
}

func NycNow() time.Time {
	z, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println("Could not get Timezone - ", err.Error())
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

func UtcNow() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.0000000Z")
}
