package utilization

import (
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
		input        []byte
		expected     []string
		expectsError bool
	}{
		{
			[]byte("cpu  0 0 0 0 0 0 0 0 0 0\n"),
			[]string{"0", "0", "0", "0", "0", "0", "0", "0", "0", "0"},
			false,
		},
		{
			[]byte("cpu  1528029 235 1406982 142714174 167167 223736 40373 0 0 0\n"),
			[]string{"1528029", "235", "1406982", "142714174", "167167", "223736", "40373", "0", "0", "0"},
			false,
		},
		{
			[]byte("cpu"),
			nil,
			true,
		},
		{
			[]byte("cpu\n"),
			nil,
			true,
		},
		{
			[]byte("cpu  0 0 0 0 0 0 0 0 0\n"),
			nil,
			true,
		},
	}

	for index, test := range tests {
		result, err := getStatValues(test.input)
		if test.expectsError {
			if err == nil {
				t.Errorf("Test %d failed. shouldFail=%v, err=%v", index, test.expectsError, err)
			}
			continue
		}

		if err != nil {
			t.Errorf("Test %d bad. Exepcted err=%v, got err=%s", index, nil, err)
			continue
		}

		if len(test.expected) != len(result) {
			t.Errorf("Test %d failed. Result expected len=%d, got len=%d", index, len(test.expected), len(result))
			continue
		}

		for i, s := range test.expected {
			if s != result[i] {
				t.Errorf("Test %d failed. Expected result[%d]=%s, got=%s", index, i, s, result[i])
			}
		}
	}
}
