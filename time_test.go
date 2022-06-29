package main

import (
	"testing"
	"time"
)

func fakeNow(t time.Time) {
	now = func() time.Time {
		return t
	}
}

func TestNow(t *testing.T) {
	expected := time.Now()
	fakeNow(expected)
	Reset()

	if got := Now(); expected != got {
		t.Errorf("Now() = %v, want %v", got, expected)
	}
}
