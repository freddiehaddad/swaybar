package utilization

import (
	"fmt"
	"testing"
)

func TestParseInts(t *testing.T) {
	tests := []struct {
		input       []string
		expected    []int64
		errExpected bool
	}{
		{[]string{}, []int64{}, false},
		{[]string{"0"}, []int64{0}, false},
		{[]string{"0", "1"}, []int64{0, 1}, false},
		{[]string{"a", "1"}, []int64{}, true},
		{[]string{"1", "a"}, []int64{}, true},
	}

	for index, test := range tests {
		results, err := parseInts(test.input)

		if test.errExpected {
			if err == nil {
				t.Errorf("Test %d failed. Expected err, got err=nil", index)
			}
			continue
		}

		if err != nil {
			t.Errorf("Test %d failed. Expected passing test, got err=%s", index, err)
			continue
		}

		if len(results) != len(test.expected) {
			t.Errorf("Test %d failed. Expected length=%d, got length=%d\n", index, len(test.expected), len(results))
			continue
		}

		for i, result := range results {
			if test.expected[i] != result {
				t.Errorf("Test %d failed at i=%d. Expected value=%d, got value=%d\n", index, i, test.expected[i], result)
			}
		}
	}
}

func TestSumArray(t *testing.T) {
	tests := []struct {
		input []int64
		sum   int64
	}{
		{[]int64{}, 0},
		{[]int64{0}, 0},
		{[]int64{0, 0}, 0},
		{[]int64{1}, 1},
		{[]int64{0, 1}, 1},
		{[]int64{1, 0}, 1},
		{[]int64{1, 1}, 2},
	}

	for index, test := range tests {
		sum := sumArray(test.input)
		if sum != test.sum {
			t.Errorf("Test %d failed. Expected sum=%d, got sum=%d", index, test.sum, sum)
		}
	}
}

func TestGetStatValues(t *testing.T) {
	tests := []struct {
		input    []byte
		expected []string
		err      error
	}{
		{
			[]byte("cpu  0 0 0 0 0 0 0 0 0 0\n"),
			[]string{"0", "0", "0", "0", "0", "0", "0", "0", "0", "0"},
			nil,
		},
		{
			[]byte("cpu  1528029 235 1406982 142714174 167167 223736 40373 0 0 0\n"),
			[]string{"1528029", "235", "1406982", "142714174", "167167", "223736", "40373", "0", "0", "0"},
			nil,
		},
		{
			[]byte("cpu"),
			nil,
			fmt.Errorf("error splitting cpu, expected a length 2, but got length 1"),
		},
		{
			[]byte("cpu\n"),
			nil,
			fmt.Errorf("error procesing values, expected length 10, but got length 1"),
		},
		{
			[]byte("cpu  0 0 0 0 0 0 0 0 0\n"),
			nil,
			fmt.Errorf("error procesing values, expected length 10, but got length 9"),
		},
	}

	for index, test := range tests {
		result, err := getStatValues(test.input)
		if test.err != nil && test.err.Error() != err.Error() {
			t.Errorf("Test %d failed. Expected err=%v, got err=%v", index, test.err, err)
			continue
		}

		if len(test.expected) != len(result) {
			t.Errorf("Test %d failed. Expected result len=%d, got len=%d", index, len(test.expected), len(result))
			continue
		}

		for i, s := range test.expected {
			if s != result[i] {
				t.Errorf("Test %d failed. Expected result[%d]=%s, got=%s", index, i, s, result[i])
			}
		}
	}
}
