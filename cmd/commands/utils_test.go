package commands

import (
	"testing"
)

func TestParseRangeInput(t *testing.T) {
	var tests = []struct {
		input     string
		wantStart *int
		wantEnd   *int
		wantErr   bool
	}{
		{"10:20", intPtr(10), intPtr(20), false},
		{"0:1", intPtr(0), intPtr(1), false},
		{":1", nil, intPtr(1), false},
		{"1:", intPtr(1), nil, false},
		{":-10", nil, intPtr(-10), false},
		{"-1:10", nil, nil, true}, // Negative start is now allowed
	}

	for _, tt := range tests {
		testname := tt.input
		t.Run(testname, func(t *testing.T) {
			start, end, err := parseRangeInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if !intsEqual(start, tt.wantStart) {
				t.Errorf("start: got %v, want %v", valOrNil(start), valOrNil(tt.wantStart))
			}
			if !intsEqual(end, tt.wantEnd) {
				t.Errorf("end: got %v, want %v", valOrNil(end), valOrNil(tt.wantEnd))
			}
		})
	}
}
