package network

import (
	"math"
	"testing"
)

// Return true if a == b within tolerance t, otherwise false
func withinTolerance(a, b, t float64) bool {
	c := math.Abs(a - b)
	return c <= t
}

func TestWithinTolerance(t *testing.T) {
	const tolerance = 0.001

	tests := []struct {
		a        float64
		b        float64
		expected bool
	}{
		{0.0, 0.0, true},
		{0.001, 0.0, true},
		{0.0, 0.001, true},
		{0.0, 0.0001, true},
		{0.0001, 0.0, true},
		{0.0, 0.0011, false},
		{0.0011, 0.0, false},
		{0.002, 0.0, false},
		{0.002, 0.0, false},
		{1.0, 1.0, true},
		{1.0, 0.0, false},
		{0.0, 1.0, false},
		{1.0, 1.0001, true},
		{1.0001, 1.0, true},
	}

	for index, test := range tests {
		result := withinTolerance(test.a, test.b, tolerance)
		if result != test.expected {
			t.Error("Test", index, "failed: a", test.a, "b", test.b, "expected", test.expected, "got", result)
		}
	}

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

		if !withinTolerance(test.throughput, throughput, 0.01) {
			t.Error("Test withinTolerance", index, "failed: input", test.bitPerSecond, "expected", test.throughput, "got", throughput)
		}

		if test.rate != rate {
			t.Error("Test", index, "failed: input", test.bitPerSecond, "expected", test.rate, "got", rate)
		}
	}
}
