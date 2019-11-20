package stats

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

type ParseCase struct {
	Input   string
	Output  []string
	Success bool
}

func TestParse(t *testing.T) {
	// setup
	parseCases := []ParseCase{
		ParseCase{
			"foo b a r",
			[]string{"foo", "b", "a", "r"},
			true,
		},
		ParseCase{
			"foo \"foo bar\"",
			[]string{"foo", "foo bar"},
			true,
		},
		ParseCase{
			"foo \"foo bar\" bar",
			[]string{"foo", "foo bar", "bar"},
			true,
		},
		ParseCase{
			"foo",
			[]string{"foo"},
			true,
		},
		ParseCase{
			"",
			[]string{""},
			true,
		},
		ParseCase{
			"米 こめ",
			[]string{"米", "こめ"},
			true,
		},
	}

	// begin test
	for i, c := range parseCases {
		out, err := Parse(strings.NewReader(c.Input))
		if err != nil {
			if c.Success && err != io.EOF {
				t.Fatalf("Test case %d failed with %s\n", i, err)
			}
			continue
		}

		for j, elem := range out {
			if c.Output[j] != elem {
				if c.Success {
					t.Fatalf("Test case %d erroneous output: element %d should be %s but it is %s\n", i, j, elem, c.Output[j])
				}
				continue
			}
		}
	}
}

func TestOutput(t *testing.T) {
	s := NewStats()
	var out strings.Builder

	Output(&out, s)

	outStr := out.String()
	j, _ := json.Marshal(s.Cache)
	expectedStr := string(j)

	if outStr != expectedStr {
		t.Fatalf("expected: %s | got %s", expectedStr, outStr)
	}
}

func TestNewStats(t *testing.T) {
	s := NewStats()

	// total
	if s.Cache.Treasury != 0.0 {
		t.Errorf("Treasury mismatch (%f --> 0.0)", s.Cache.Treasury)
	}

	// month
	if s.Cache.Month.Total != 0.0 {
		t.Errorf("Month Total (%f --> 0.0)", s.Cache.Month.Total)
	}
	if s.Cache.Month.Income != 0.0 {
		t.Errorf("Month Income (%f --> 0.0)", s.Cache.Month.Income)
	}
	if s.Cache.Month.Expenses != 0.0 {
		t.Errorf("Month Expenses (%f --> 0.0)", s.Cache.Month.Expenses)
	}

	// year
	if s.Cache.Year.Total != 0.0 {
		t.Errorf("Year Total (%f --> 0.0)", s.Cache.Year.Total)
	}
	if s.Cache.Year.Income != 0.0 {
		t.Errorf("Year Total (%f --> 0.0)", s.Cache.Year.Income)
	}
	if s.Cache.Year.Expenses != 0.0 {
		t.Errorf("Year Expenses (%f --> 0.0", s.Cache.Year.Expenses)
	}

	// events
	if len(s.Events) != 0 {
		t.Errorf("Events should be empty (%v)", s.Events)
	}

	// transactions
	if len(s.Transactions) != 0 {
		t.Errorf("Transactions should be empty (%v)", s.Transactions)
	}
}
