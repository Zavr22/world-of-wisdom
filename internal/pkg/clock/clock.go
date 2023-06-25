package clock

import "time"

type SystemClock struct{}

func (c *SystemClock) Now() time.Time {
	return time.Now()
}

func (c *SystemClock) Sleep(duration time.Duration) {
	time.Sleep(duration)
}
