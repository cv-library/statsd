package statsd

import (
	"bytes"
	"net"
	"strconv"
	"testing"
	"time"
	"math/rand"
)

func TestTime(t *testing.T) {
	test( t, func (conn *net.UDPConn) {
		dur, err := time.ParseDuration("1h30m45s")
		if err != nil {
			t.Fatal(err)
		}

		Time("foo", dur)

		exp := []byte("foo:" +
			strconv.FormatUint(uint64(dur.Nanoseconds()/1e6), 10) + "|ms")

		got := make([]byte, len(exp))
		if _, _, err := conn.ReadFromUDP(got); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, exp) {
			t.Errorf("got: %s; want: %s", got, exp)
		}
	})
}

func TestTimeSampled(t *testing.T) {
	test( t, func (conn *net.UDPConn) {
		dur, err := time.ParseDuration("1h30m45s")
		if err != nil {
			t.Fatal(err)
		}

		TimeSampled(0.5, "foo", dur)

		exp := []byte("foo:" +
			strconv.FormatUint(uint64(dur.Nanoseconds()/1e6), 10) + "|ms|@0.5")

		got := make([]byte, len(exp))
		if _, _, err := conn.ReadFromUDP(got); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, exp) {
			t.Errorf("got: %s; want: %s", got, exp)
		}
	})
}

func TestTimerReset(t *testing.T) {
	timer := Timer()
	time.Sleep(time.Millisecond)

	then := timer.start
	timer.Reset()

	if then == timer.start {
		t.Errorf("timer did not reset: was %s; is %s", then, timer.start)
	}
}

func TestTimerSend(t *testing.T) {
	test( t, func (conn *net.UDPConn) {
		timer := Timer()
		time.Sleep(time.Millisecond)

		took := timer.Send("foo")
		if took == 0 {
			t.Error("Send() took no time")
		}

		exp := []byte("foo:" +
			strconv.FormatUint(uint64(took.Nanoseconds()/1e6), 10) + "|ms")

		got := make([]byte, len(exp))
		if _, _, err := conn.ReadFromUDP(got); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, exp) {
			t.Errorf("got: %s; want: %s", got, exp)
		}
	})
}

func TestTimerSendSampled(t *testing.T) {
	test( t, func (conn *net.UDPConn) {
		timer := Timer()
		time.Sleep(time.Millisecond)

		took := timer.SendSampled(0.5, "foo")
		if took == 0 {
			t.Error("Send() took no time")
		}

		exp := []byte("foo:" +
			strconv.FormatUint(uint64(took.Nanoseconds()/1e6), 10) + "|ms|@0.5")

		got := make([]byte, len(exp))
		if _, _, err := conn.ReadFromUDP(got); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, exp) {
			t.Errorf("got: %s; want: %s", got, exp)
		}
	})
}

func TestInc(t *testing.T) {
	test( t, func (conn *net.UDPConn) {
		Inc("foo")

		exp := []byte("foo:1|c")
		got := make([]byte, len(exp))
		if _, _, err := conn.ReadFromUDP(got); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, exp) {
			t.Errorf("got: %s; want: %s", got, exp)
		}
	})
}

func TestIncSampled(t *testing.T) {
	test( t, func (conn *net.UDPConn) {
		IncSampled(0.5, "foo")

		exp := []byte("foo:1|c|@0.5")
		got := make([]byte, len(exp))
		if _, _, err := conn.ReadFromUDP(got); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, exp) {
			t.Errorf("got: %s; want: %s", got, exp)
		}
	})
}

func TestGauge(t *testing.T) {
	test( t, func (conn *net.UDPConn) {
		Gauge("foo", 42)

		exp := []byte("foo:42|g")
		got := make([]byte, len(exp))
		if _, _, err := conn.ReadFromUDP(got); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, exp) {
			t.Errorf("got: %s; want: %s", got, exp)
		}
	})
}

func TestGaugeSampled(t *testing.T) {
	test( t, func (conn *net.UDPConn) {
		GaugeSampled(0.5, "foo", 42)

		exp := []byte("foo:42|g|@0.5")
		got := make([]byte, len(exp))
		if _, _, err := conn.ReadFromUDP(got); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, exp) {
			t.Errorf("got: %s; want: %s", got, exp)
		}
	})
}

func test (t *testing.T, check func(*net.UDPConn) ) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 8125})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second))

	rand.Seed(1)
	check(conn)
}
