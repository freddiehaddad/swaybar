package network

import (
	"math"
	"testing"
)

// Return true if a == b within tolerance e, otherwise false
func withinToleranc(a, b, e float64) bool {
	if a == b {
		return true
	}

	d := math.Abs(a - b)
	if b == 0 {
		return d < e
	}

	return (d / math.Abs(b)) < e
}

func TestConvertBytesToBits(t *testing.T) {
	tests := []struct {
		bytes int64
		bits  int64
	}{
		{0, 0},
		{1, 8},
		{2, 16},
	}
	for index, test := range tests {
		bits := convertBytesToBits(test.bytes)
		if test.bits != bits {
			t.Error("Test", index, "failed: input", test.bytes, "expected", test.bits, "got", bits)
		}
	}
}

func TestConvertNanosecondsToSeconds(t *testing.T) {
	tests := []struct {
		nanoseconds int64
		seconds     float64
	}{
		{0, 0},
		{25e7, 0.25},
		{5e8, 0.5},
		{1e9, 1.0},
		{2e9, 2.0},
		{2e9 + 5e8, 2.5},
	}
	for index, test := range tests {
		seconds := convertNanosecondsToSeconds(test.nanoseconds)
		if test.seconds != seconds {
			t.Error("Test", index, "failed: input", test.nanoseconds, "expected", test.seconds, "got", seconds)
		}
	}
}

func TestShortenThroughput(t *testing.T) {
	tests := []struct {
		bitPerSecond float64
		throughput   float64
		rate         string
	}{
		{0, 0, "bps"},
		{1, 1, "bps"},
		{1<<10 - 1, 1023, "bps"},
		{1 << 10, 1, "Kbps"},
		{1<<10 + 1, 1, "Kbps"},
		{1<<20 - 1, 1023.99, "Kbps"},
		{1 << 20, 1, "Mbps"},
		{1<<20 + 1, 1, "Mbps"},
		{1<<30 - 1, 1023.99, "Mbps"},
		{1 << 30, 1, "Gbps"},
		{1<<30 + 1, 1, "Gbps"},
	}
	for index, test := range tests {
		throughput, rate := shortenThroughput(test.bitPerSecond)

		if !withinToleranc(test.throughput, throughput, 1e9) {
			t.Error("Test", index, "failed: input", test.bitPerSecond, "expected", test.throughput, "got", throughput)
		}

		if test.rate != rate {
			t.Error("Test", index, "failed: input", test.bitPerSecond, "expected", test.rate, "got", rate)
		}
	}
}
