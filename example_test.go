package statsd_test

import (
	"fmt"
	"time"
	"github.com/cv-library/statsd"
)

func ExampleTiming() {
	for i:=0; i<10; i++ {
		timer := statsd.Timer()

		// Do work

		took := timer.Send("work.duration")
		timer.SendWithOptions(
			&statsd.Options{ Rate: 0.25 },
			"work.duration.sampled",
		)

		fmt.Printf("Took %f seconds\n", took.Seconds())
	}
}

func ExampleTimer() {
	// This is slightly more efficient than the Timing example
	// since we avoid creating a new Timing every time
	timer := statsd.Timer()
	for i:=0; i<10; i++ {
		timer.Reset()

		// Do work

		took := timer.Send("work.duration")
		timer.SendWithOptions(
			&statsd.Options{ Rate: 0.25 },
			"work.duration.sampled",
		)

		fmt.Printf("Took %f seconds\n", took.Seconds())
	}
}

func ExampleTime() {
	for i:=0; i<10; i++ {
		start := time.Now()

		// Do work

		statsd.Time("work.duration", time.Since(start))
	}
}

func ExampleTimeWithOptions() {
	for i:=0; i<10; i++ {
		start := time.Now()

		// Do work

		if i % 2 == 0 {
			continue
		}

		// We use AlwaysSend: true here because we're doing the sampling
		// ourselves; we want every packet to reach the server
		statsd.TimeWithOptions(
			&statsd.Options{ Rate: 0.5, AlwaysSend: true },
			"work.duration.sampled",
			time.Since(start),
		)
	}
}

func ExampleInc() {
	// Increment a counter
	statsd.Inc("stats.success")
}

func ExampleIncWithOptions() {
	// Increment a counter every other time
	statsd.IncWithOptions(
		&statsd.Options{ Rate: 0.5 },
		"stats.success",
	)
}

func ExampleGauge() {
	// Set page size to 10
	statsd.Gauge("page.size", 10)
}

func ExampleGaugeWithOptions() {
	// Set page size to 10 every other time
	statsd.GaugeWithOptions(
		&statsd.Options{ Rate: 0.5 },
		"page.size",
		10,
	)
}
