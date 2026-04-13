package dealer

import (
	"math"
	"math/rand"
	"time"
)

func LinearDelay(base time.Duration, increment time.Duration) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + increment*time.Duration(att-1)
	}
}

// ExponentialDelay / powBase -- seconds
func ExponentialDelay(base time.Duration, powBase int64) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + float64DurSec(math.Pow(float64(powBase), float64(att-1)))*time.Second
	}
}

func Exponential2Delay(base time.Duration) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + float64DurSec(math.Pow(2, float64(att-1)))
	}
}

func ExponentialJitterDelay(base time.Duration, pow int64) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + float64DurSec(math.Pow(float64(pow), float64(att-1))*(0.5+0.5*rand.Float64()))
	}
}

func ExponentialJitter2Delay(base time.Duration) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + float64DurSec(math.Pow(2, float64(att-1))*(0.5+0.5*rand.Float64()))
	}
}

func float64DurSec(f float64) time.Duration {
	return time.Duration(f) * time.Second
}
