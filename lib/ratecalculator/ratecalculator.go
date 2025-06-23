package ratecalculator

import (
	"fmt"
	"time"
)

type RateCalculator interface {
	WaitToNextMessage(time.Duration, uint64) time.Duration
}

type RateCalculatorImpl struct {
	Frequency int           // messages per period
	Period    time.Duration // period for which the rate is calculated
}

func NewRateCalculator(rate int, period time.Duration) (*RateCalculatorImpl, error) {
	if rate <= 0 {
		return nil, fmt.Errorf("rate must be a positive integer, got %d", rate)
	}
	if period <= 0 {
		return nil, fmt.Errorf("period must be a positive duration, got %v", period)
	}
	return &RateCalculatorImpl{Frequency: rate, Period: period}, nil
}

func (rc *RateCalculatorImpl) WaitToNextMessage(timeFromStart time.Duration, messageNumber uint64) time.Duration {
	expectedNextNumber := uint64(rc.Frequency) * uint64(timeFromStart/rc.Period)

	if messageNumber < expectedNextNumber {
		return 0 // send immediately
	}

	interval := uint64(rc.Period.Nanoseconds() / int64(rc.Frequency))

	nextMessageTime := time.Duration((messageNumber + 1) * interval)
	waitTime := nextMessageTime - timeFromStart
	return waitTime
}
