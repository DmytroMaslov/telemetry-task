package ratecalculator

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func Test_WaitToNextMessage(t *testing.T) {
	type args struct {
		frequency     int
		timeFromStart time.Duration
		messageNumber uint64
	}
	testCases := []struct {
		name     string
		args     args
		expected time.Duration
	}{
		{
			name: "first message",
			args: args{
				frequency:     100,
				timeFromStart: 0 * time.Second,
				messageNumber: 0,
			},
			expected: 10 * time.Millisecond,
		},
		{
			name: "last message",
			args: args{
				frequency:     100,
				timeFromStart: 990 * time.Millisecond,
				messageNumber: 99,
			},
			expected: 10 * time.Millisecond,
		},
		{
			name: "send immediately",
			args: args{
				frequency:     100,
				timeFromStart: 990 * time.Millisecond,
				messageNumber: 98,
			},
			expected: 0 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc := &RateCalculatorImpl{
				Frequency: tc.args.frequency,
				Period:    time.Second,
			}

			actual := rc.WaitToNextMessage(tc.args.timeFromStart, tc.args.messageNumber)

			assert.Equal(t, tc.expected, actual)
		})
	}

}
