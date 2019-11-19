package stats

import (
	"io"
	"strings"
	"testing"
	"time"
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

type UpdateCase struct {
	Input   Stats
	Output  Stats
	Success bool
}

func UpdateTest(t *testing.T) {
	/** setup **/
	// test case 1
	testCases := []UpdateCase{
		UpdateCase{
			NewStats(),
			NewStats(),
			true,
		},
	}

	// test case 2
	input := NewStats()
	output := NewStats()

	tr := Tr{
		0,
		"test_transaction",
		time.Now(),
		100,
		0,
		"test description",
	}
	input.Transactions = append(input.Transactions, tr)
	output.Transactions = append(output.Transactions, tr)
	output.Cache.Treasury = 100

	testCases = append(testCases, UpdateCase{
		input,
		output,
		true,
	})
	/*********/

	// begin test
	for i, tc := range testCases {
		err := Update(&tc.Input)
		if err != nil {
			if tc.Success {
				t.Fatalf("test %d failed with %s", i, err)
			}
			continue
		}

		// check cache
		if tc.Input.Cache.Treasury != tc.Output.Cache.Treasury {
			if tc.Success {
				t.Fatalf("test %d treasury mismatch", i)
			}
		}
		if tc.Input.Cache.Month != tc.Output.Cache.Month {
			if tc.Success {
				t.Fatalf("test %d month mismatch", i)
			}
		}
		if tc.Input.Cache.Year != tc.Output.Cache.Year {
			if tc.Success {
				t.Fatalf("test %d year mismatch", i)
			}
		}

		// check transaction
		for j, transaction := range tc.Input.Transactions {
			if transaction != tc.Output.Transactions[j] {
				if tc.Success {
					t.Fatalf("test %d transaction mismatch", i)
				}
			}
		}

		// check events
		for j, e := range tc.Input.Events {
			if e != tc.Output.Events[j] {
				if tc.Success {
					t.Fatalf("test %d events mismatch", i)
				}
			}
		}
	}
}
