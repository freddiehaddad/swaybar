package utils

import (
	"io/fs"
	"testing"
)

func TestGetSensorValue(t *testing.T) {
	tests := []struct {
		file       string
		expected   int64
		shouldFail bool
	}{
		{"tests/bad", 0, true},
		{"tests/nofile", 0, true},
		{"tests/0", 0, false},
		{"tests/10000", 10000, false},
	}

	for index, test := range tests {
		result, err := GetSensorValue(test.file)

		if test.shouldFail {
			if err == nil {
				t.Errorf("Test %d failed shouldFail=%v err=%v", index, test.shouldFail, err)
			}
			continue
		}

		if err != nil {
			t.Errorf("Test %d failed with unexpected err=%v", index, err)
		}

		if result != test.expected {
			t.Errorf("Test %d failed result=%d expected=%d", index, result, test.expected)
		}
	}
}

func TestReadSensorValue(t *testing.T) {
	tests := []struct {
		input    int64
		expected float64
	}{
		{0, 0.0},
		{10000, 10.0},
		{10500, 10.5},
	}

	for index, test := range tests {
		result := ReadSensorValue(test.input)
		if result != test.expected {
			t.Error("Test", index, "failed:", "input", test.input, "expected", test.expected, "got", result)
		}
	}
}

func TestSymbolicLink(t *testing.T) {
	tests := []struct {
		input    fs.FileMode
		expected bool
	}{
		{fs.ModeSymlink, true},
		{0xFFFFFFFF ^ fs.ModeSymlink, false},
	}

	for index, test := range tests {
		result := symbolicLink(test.input)
		if result != test.expected {
			t.Errorf("Test %d failed result=%v expected=%v\n", index, result, test.expected)
		}
	}
}
