package statsd

import (
	"bytes"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestTimerSend(t *testing.T) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 8125})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second))

	timer := NewTimer()

	time.Sleep(time.Millisecond)

	took := timer.Send("foo")

	if took == 0 {
		t.Error("Send() took no time")
	}

	got := make([]byte, 8)
	if _, _, err := conn.ReadFromUDP(got); err != nil {
		t.Fatal(err)
	}

	exp := []byte("foo:" +
		strconv.FormatUint(uint64(took.Nanoseconds()/1e6), 10) + "|ms")

	if !bytes.Equal(got, exp) {
		t.Errorf("got: %s; want: %s", got, exp)
	}
}
