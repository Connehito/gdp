package main

import (
	"time"
)

var fakeTime time.Time
var now = time.Now

//Now returns the current time
func Now() time.Time {
	if !fakeTime.IsZero() {
		return fakeTime
	}
	return now()
}

//Set fake time
func Set(t time.Time) {
	fakeTime = t
}

//Reset fake time
func Reset() {
	fakeTime = time.Time{}
}
