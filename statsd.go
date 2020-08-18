package statsd

import (
	"math/rand"
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

	rand.Seed(time.Now().UTC().UnixNano())

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

// Send takes a list of remote timer names, and submits the time that
// has ellapsed since the creation of the timer to each in turn.
// It returns a time.Duration representing the amount of time that was sent.
func (t *timer) Send(names ...interface{}) (took time.Duration) {
	return t.SendSampled(1.0, names...)
}

// SendSampled works like Send but sends the timing information
// with the given sample rate.
func (t *timer) SendSampled(rate float64, names ...interface{}) (took time.Duration) {
	took = time.Since(t.start)

	var message string
	if sampled, suffix := sampleData(rate); sampled {
		message = ":" +
			strconv.FormatUint(uint64(took.Nanoseconds()/1e6), 10) +
			"|ms" + suffix
	} else {
		return
	}

	for _, name := range names {
		send(name.(string), message)
	}

	return
}

// Gauge sets arbitrary numeric value for a given metric.
func Gauge(name string, value int64) {
	GaugeSampled(1.0, name, value)
	return
}

// GaugeSampled sets arbitrary numeric value for a given metric
// with the given sample rate.
func GaugeSampled(rate float64, name string, value int64) {
	var message string
	if sampled, suffix := sampleData(rate); sampled {
		message = ":" + strconv.FormatInt(value, 10) + "|g" + suffix
	} else {
		return
	}

	send(name, message)
}

// Inc increments a counter.
func Inc(name string) {
	IncSampled(1.0, name)
}

// IncSampled increments a counter with the given sample rate.
func IncSampled(rate float64, name string) {
	var message string
	if sampled, suffix := sampleData(rate); sampled {
		message = ":1|c" + suffix
	} else {
		return
	}

	send(name, message)
}

// Time sends duration in ms for a given metric.
func Time(name string, took time.Duration) {
	TimeSampled(1.0, name, took)
}

// TimeSampled sends duration in ms for a given metric
// with the given sample rate.
func TimeSampled(rate float64, name string, took time.Duration) {
	var message string
	if sampled, suffix := sampleData(rate); sampled {
		message = ":" +
			strconv.FormatUint(uint64(took.Nanoseconds()/1e6), 10) +
			"|ms" + suffix
	} else {
		return
	}

	send(name, message)
}

func send(name string, message string) {
	if err := getConnection(); err != nil {
		return
	}

	conn.Write([]byte(name + message))

	// Send a host suffixed stat too.
	if AlsoAppendHost {
		conn.Write([]byte(name + "." + host + message))
	}

	return
}

func sampleData(rate float64) (bool, string) {
	if rate == 1.0 {
		return true, ""
	}

	if rand.Float64() >= rate {
		return false, ""
	}

	return true, "|@" + strconv.FormatFloat(rate, 'f', -1, 64)
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
