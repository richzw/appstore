package appstore

import (
	"math/rand"
	"time"
)

type Backoff interface {
	Pause() time.Duration
}

type OneSecondBackoff struct{}

func (bo *OneSecondBackoff) Pause() time.Duration {
	return time.Duration(1) * time.Second
}

type JitterBackoff struct {
	Initial    time.Duration
	Max        time.Duration
	Multiplier float64

	cur time.Duration
}

func (bo *JitterBackoff) Pause() time.Duration {
	if bo.Initial == 0 {
		bo.Initial = time.Second
	}
	if bo.cur == 0 {
		bo.cur = bo.Initial
	}
	if bo.Max == 0 {
		bo.Max = 30 * time.Second
	}
	if bo.Multiplier < 1 {
		bo.Multiplier = 2
	}

	// https://www.awsarchitectureblog.com/2015/03/backoff.html
	d := time.Duration(1 + rand.Int63n(int64(bo.cur)))
	bo.cur = time.Duration(float64(bo.cur) * bo.Multiplier)

	// stop the backoff
	if bo.cur > bo.Max {
		d = -1
	}
	return d
}
