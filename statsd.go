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

func (t *timer) Send(names ...interface{}) (took time.Duration) {
	took = time.Since(t.start)

	value := ":" + strconv.FormatUint(uint64(took.Nanoseconds()/1e6), 10) + "|ms"

	if err := getConnection(); err != nil {
		return
	}

	for _, name := range names {
		conn.Write([]byte(name.(string) + value))
	}

	return
}

// Inc is a simple counter adding one to a given metric.
func Inc(name string) {
	if err := getConnection(); err != nil {
		return
	}

	conn.Write([]byte(name + ":1|c"))

	return
}

func getConnection() (err error) {
	// If we don't have a conn, make one.
	if conn == nil {
		if conn, err = net.Dial("udp", Address); err != nil {
			conn = nil
			return
		}
	}

	return
}
