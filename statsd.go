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

// DefaultOptions holds the default options, used by the "...WithOptions"
// functions when `nil` is passed
var DefaultOptions = &Options{
	Rate:       1.0,
	AlwaysSend: false,
}

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

// Options holds key/value pairs to be used when calling the functions
// with the "...WithOptions" suffix.
type Options struct {
	// Rate specifies the sampling rate to use, between 0 and 1. A value of
	// 0.1 means that this value was sampled 1 out of every 10 times. Under
	// normal circumstances, the remaining 9 out of 10 times will send no data
	// to the server, to reduce network use. The `AlwaysSend` option below can
	// be used to change this behaviour.
	Rate float64
	// AlwaysSend specifies whether the packet should be sent to the
	// server regardless of the sampling rate used. Under normal circumstances,
	// every call that would result in sending data to the server checks the
	// sampling rate and uses that to decide whether to send the data or not
	// (if the rate is 0.1, only 1 out of 10 packets would be sent). When this
	// option is set to "true", the packet will be sent regardless. This is useful
	// if the sampling is being done on the calling side.
	AlwaysSend bool
}

// Timer
type timer struct {
	start time.Time
}

// Timer returns a new timer set to `time.Now()`
func Timer() timer {
	return timer{time.Now()}
}

// Reset sets the start time for the timer to `time.Now()`
func (t *timer) Reset() {
	t.start = time.Now()
}

// Send takes a list of remote timer names, and submits the time that
// has ellapsed since the creation of the timer to each in turn.
// It returns a time.Duration representing the amount of time that was sent.
func (t *timer) Send(names ...interface{}) (took time.Duration) {
	return t.SendWithOptions(nil, names...)
}

// SendWithOptions works like Send but sends the timing information
// using the provided options.
func (t *timer) SendWithOptions(
	options *Options,
	names ...interface{},
) (took time.Duration) {
	took = time.Since(t.start)

	var message string
	if sampled, suffix := check(options); sampled {
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
	GaugeWithOptions(nil, name, value)
	return
}

// GaugeWithOptions sets arbitrary numeric value for a given metric
// using the provided options.
func GaugeWithOptions(options *Options, name string, value int64) {
	var message string
	if sampled, suffix := check(options); sampled {
		message = ":" + strconv.FormatInt(value, 10) + "|g" + suffix
	} else {
		return
	}

	send(name, message)
}

// Inc increments a counter.
func Inc(name string) {
	IncWithOptions(nil, name)
}

// IncSampled increments a counter with the given sample rate.
// Note that this function will send the data to the server every time
// it is called. It is the caller's responsibility to implement the
// sampling.
//
// Deprecated: Use IncWithOptions and pass the sampling rate
// using the Options struct.
func IncSampled(name string, rate float64) {
	IncWithOptions(&Options{Rate: rate, AlwaysSend: true}, name)
}

// IncWithOptions increments a counter using the provided options.
func IncWithOptions(options *Options, name string) {
	var message string
	if sampled, suffix := check(options); sampled {
		message = ":1|c" + suffix
	} else {
		return
	}

	send(name, message)
}

// Time sends duration in ms for a given metric.
func Time(name string, took time.Duration) {
	TimeWithOptions(nil, name, took)
}

// TimeWithOptions sends duration in ms for a given metric
// using the provided options.
func TimeWithOptions(options *Options, name string, took time.Duration) {
	var message string
	if sampled, suffix := check(options); sampled {
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

func check(options *Options) (bool, string) {
	if options == nil {
		options = DefaultOptions
	}

	if !options.AlwaysSend && rand.Float64() >= options.Rate {
		return false, ""
	}

	return true, "|@" + strconv.FormatFloat(options.Rate, 'f', -1, 64)
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
