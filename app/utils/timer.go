package utils

import (
	"time"
)

/*
	Why not just do time.Now()? Well, this is to get around the problems time.Now() causes in unit tests. In a test
	scenario, I'll pass a pre-defined datetime to my test suite as a fake timer, so that dynamic values become static
*/
type Timer interface {
	GetTimeNow() time.Time
}

type RealTimer struct{}

func (t *RealTimer) GetTimeNow() time.Time {
	return time.Now()
}

type FakeTimer struct {
	t time.Time
}

func NewFakeTimer(t time.Time) *FakeTimer {
	return &FakeTimer{t}
}

func (t *FakeTimer) GetTimeNow() time.Time {
	return t.t
}
