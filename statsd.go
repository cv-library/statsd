package statsd

import (
	"net"
	"strconv"
	"time"
)

var Address = "localhost:8125"

// Cache the conn for perf.
var conn net.Conn

func Timer() timer {
	return timer{time.Now()}
}

// Timer
type timer struct {
	start time.Time
}

func (t *timer) Reset() {
	t.start = time.Now()
}

func (t *timer) Send(names ...interface{}) {
	value := ":" + strconv.FormatUint(
		uint64((time.Now().UnixNano()-t.start.UnixNano())/1e6),
		10,
	) + "|ms"

	// If we don't have a conn, make one.
	if conn == nil {
		var err error
		if conn, err = net.Dial("udp", Address); err != nil {
			conn = nil
			return
		}
	}

	for _, name := range names {
		conn.Write([]byte(name.(string) + value))
	}

	return
}
