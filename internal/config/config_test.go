package config

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

func setEnvs(t *testing.T) {
	t.Setenv("NUMBERS", "12345")
	t.Setenv("TEST", "@ I AM COOL")
	t.Setenv("EMPTY", "")
	t.Setenv("CSV", "COMMA,SEPARATED,VALUES")
	t.Setenv("DURATION", "45m")
	t.Setenv("DURATION_MALFORMED", "hf83h")
	t.Setenv("COMMAS", ",,,")
}

func TestGet(t *testing.T) {
	setEnvs(t)

	var tests = []struct {
		key, def, expected string
	}{
		{"NUMBERS", "", "12345"},
		{"TEST", "", "@ I AM COOL"},
		{"EMPTY", "", ""},
		{"CSV", "", "COMMA,SEPARATED,VALUES"},
		{"CSV", "DEFAULT", "COMMA,SEPARATED,VALUES"},
		{"NOT_EXISTING", "DEFAULT", "DEFAULT"},
		{"NOT_EXISTING", "", ""},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("\"%s\"->\"%s\"", test.key, test.expected)
		t.Run(testname, func(t *testing.T) {
			actual := get(test.key, test.def)
			if actual != test.expected {
				t.Errorf("expected %s, got %s", test.expected, actual)
			}
		})
	}
}

func TestGetDuration(t *testing.T) {
	setEnvs(t)

	var tests = []struct {
		key           string
		def, expected time.Duration
		expectErr     bool
	}{
		{"DURATION", 0, 45 * time.Minute, false},
		{"DURATION_MALFORMED", 0, 0, true},
		{"NOT_EXISTING", 0, 0, false},
		{"NOT_EXISTING", 1 * time.Hour, 1 * time.Hour, false},
	}

	for _, test := range tests {
		errExpected := "no err"
		if test.expectErr {
			errExpected = "err"
		}
		testname := fmt.Sprintf("\"%s\"->\"%s\"(%s)", test.key, test.expected, errExpected)

		t.Run(testname, func(t *testing.T) {
			actual, err := getDuration(test.key, test.def)
			if err != nil && !test.expectErr {
				t.Errorf("unexpected error: %v", err)
			}

			if actual != test.expected {
				t.Errorf("expected %s, got %s", test.expected, actual)
			}
		})
	}
}

func TestGetCSV(t *testing.T) {
	setEnvs(t)

	var tests = []struct {
		key, def string
		expected []string
	}{
		{"NUMBERS", "", []string{"12345"}},
		{"TEST", "", []string{"@ I AM COOL"}},
		{"EMPTY", "", []string{""}},
		{"CSV", "", []string{"COMMA", "SEPARATED", "VALUES"}},
		{"CSV", "DEFAULT", []string{"COMMA", "SEPARATED", "VALUES"}},
		{"NOT_EXISTING", "DEFAULT", []string{"DEFAULT"}},
		{"NOT_EXISTING", "", []string{""}},
		{"COMMAS", "", []string{"", "", "", ""}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("\"%s\"->\"%s\"", test.key, test.expected)
		t.Run(testname, func(t *testing.T) {
			actual := getCSV(test.key, test.def)
			if slices.Compare(actual, test.expected) != 0 {
				t.Errorf("expected %s, got %s", test.expected, actual)
			}
		})
	}
}
func TestMust(t *testing.T) {
	setEnvs(t)

	var tests = []struct {
		key, expected string
		expectErr     bool
	}{
		{"NUMBERS", "12345", false},
		{"TEST", "@ I AM COOL", false},
		{"EMPTY", "", true},
		{"CSV", "COMMA,SEPARATED,VALUES", false},
		{"NOT_EXISTING", "", true},
	}

	for _, test := range tests {
		errExpected := "no err"
		if test.expectErr {
			errExpected = "err"
		}
		testname := fmt.Sprintf("\"%s\"->\"%s\"(%s)", test.key, test.expected, errExpected)

		t.Run(testname, func(t *testing.T) {
			actual, err := must(test.key)
			if err != nil && !test.expectErr {
				t.Errorf("unexpected error: %v", err)
			}

			if actual != test.expected {
				t.Errorf("expected %s, got %s", test.expected, actual)
			}
		})
	}
}
