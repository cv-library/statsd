package statsd

import (
	"testing"
	"time"
)

func TestTimerSend(t *testing.T) {
	timer := Timer()

	time.Sleep(time.Millisecond)

	if timer.Send("foo").Nanoseconds() == 0 {
		t.Fail()
	}
}
