package main

import "time"

var fakeTime time.Time

func Now() time.Time {
	if !fakeTime.IsZero() {
		return fakeTime
	}
	return time.Now()
}

func Set(t time.Time) {
	fakeTime = t
}
