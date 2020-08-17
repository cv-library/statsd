package statsd

import (
	"net"
	"os"
	"strconv"
	"time"
)

// Address of the StatsD server.
var Address = "localhost:8125"

// AlsoAppendHost appends hostname along with any metric sent.
var AlsoAppendHost = true

// Cache the conn for perf.
var conn net.Conn
var host string

func init() {
	var err error

	if host, err = os.Hostname(); err != nil {
		panic(err)
	}
}

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

	if err := getConnection(); err != nil {
		return
	}

	value := ":" + strconv.FormatUint(uint64(took.Nanoseconds()/1e6), 10) + "|ms"

	for _, name := range names {
		conn.Write([]byte(name.(string) + value))

		// Send a host suffixed stat too.
		if AlsoAppendHost {
			conn.Write([]byte(name.(string) + "." + host + value))
		}
	}

	return
}

// Gauge sets arbitrary numeric value for a given metric.
func Gauge(name string, value int64) {
	if err := getConnection(); err != nil {
		return
	}

	suffix := ":" + strconv.FormatInt(value, 10) + "|g"

	conn.Write([]byte(name + suffix))

	// Send a host suffixed stat too.
	if AlsoAppendHost {
		conn.Write([]byte(name + "." + host + suffix))
	}

	return
}

// Inc is a simple counter adding one to a given metric.
func Inc(name string) {
	return IncSampled(1.0, name)
}

// IncSampled increments a counter with the given sample rate (between 0.0 - 1.0).
// Example: Passing 0.1 tells statsd that this metric has been sample 1/10 times.
// IE Reported metric will be 10x the number of times called.
func IncSampled(rate float64, name string){
	if err := getConnection(); err != nil {
		return
	}

	message := ":1|c";
	if rate < 1.0 {
		message = message + "|@" + strconv.FormatFloat(rate, 'f', -1, 64)
	}

	conn.Write([]byte(name + message))

	// Send a host suffixed stat too.
	if AlsoAppendHost {
		conn.Write([]byte(name + "." + host + message))
	}

	return
}

// Time sends duration in ms for a given metric.
func Time(name string, took time.Duration) {
	if err := getConnection(); err != nil {
		return
	}

	value := ":" + strconv.FormatUint(uint64(took.Nanoseconds()/1e6), 10) + "|ms"

	conn.Write([]byte(name + value))

	// Send a host suffixed stat too.
	if AlsoAppendHost {
		conn.Write([]byte(name + "." + host + value))
	}

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
