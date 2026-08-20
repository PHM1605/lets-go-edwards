package main

import (
	"lets-go-edwards/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)), // Central European Time is 1 hour behind UTC
			want: "17 Mar 2024 at 09:15",
		},
	}

	for _, tt := range tests {
		// 1st parameter of t.Run(): name of the sub-test
		// 2nd parameter of t.Run(): actual test
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			// Use Test Helper "assert" (our function) here
			assert.Equal(t, hd, tt.want)
		})
	}

}
