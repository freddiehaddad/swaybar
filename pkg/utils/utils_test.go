package utils

import "testing"

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
