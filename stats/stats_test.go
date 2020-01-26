package stats

import (
	"encoding/json"
	"fmt"
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

type UpdateCase struct {
	Input   Stats
	Output  Stats
	Success bool
}

type ProcessCase struct {
	Input    []string
	InState  Stats
	OutState Stats
	Success  bool
}

type OutputCase struct {
	Input  Stats
	Output string
}

func cmpStats(in Stats, out Stats) (err []string) {
	// transactions check
	if len(in.Transactions) == len(out.Transactions) {
		for i := 0; i < len(in.Transactions); i++ {
			err = append(err, cmpTr(in.Transactions[i], out.Transactions[i])...)
		}

	} else if len(in.Transactions) > len(out.Transactions) {
		err = append(err, fmt.Sprintf("Transaction slice size mismatch (%v > %v)", in.Transactions, out.Transactions))

	} else {
		err = append(err, fmt.Sprintf("Transactions slice size mismatch (%v < %v)", in.Transactions, out.Transactions))
	}

	// events check
	if len(in.Events) == len(out.Events) {
		for i := 0; i < len(in.Events); i++ {
			err = append(err, cmpEv(in.Events[i], out.Events[i])...)
		}

	} else if len(in.Events) > len(out.Events) {
		err = append(err, fmt.Sprintf("Events slice size mismatch (%v > %v)", in.Events, out.Events))
	} else {
		err = append(err, fmt.Sprintf("Events slice size mismatch (%v < %v)", in.Events, out.Events))
	}

	// cache check
	err = append(err, cmpCache(in.Cache, out.Cache)...)

	// last check check
	if !in.LastCheck.Truncate(24 * time.Hour).Equal(out.LastCheck.Truncate(24 * time.Hour)) {
		err = append(err, fmt.Sprintf("Last check mismatch (%s should be %s)", in.LastCheck, out.LastCheck))
	}

	// index check
	if in.Index != out.Index {
		err = append(err, fmt.Sprintf("Index mismatch (%d should be %d)", in.Index, out.Index))
	}

	return
}

func cmpTr(in Tr, out Tr) (err []string) {
	// Id check
	if in.Id != out.Id {
		err = append(err, fmt.Sprintf("Tr Id mismatch (%d should be %d)", in.Id, out.Id))
	}

	// name check
	if in.Name != out.Name {
		err = append(err, fmt.Sprintf("Tr Name mismatch (%s should be %s)", in.Name, out.Name))
	}

	// date check
	if !in.Date.Truncate(24 * time.Hour).Equal(out.Date.Truncate(24 * time.Hour)) {
		err = append(err, fmt.Sprintf("Tr Date mismatch (%s should be %s)", in.Date, out.Date))
	}

	// amount check
	if in.Amount != out.Amount {
		err = append(err, fmt.Sprintf("Tr Amount mismatch (%f should be %f)", in.Amount, out.Amount))
	}

	// type check
	if in.Type != out.Type {
		err = append(err, fmt.Sprintf("Tr Type mismatch (%d should be %d)", in.Type, out.Type))
	}

	// description check
	if in.Description != out.Description {
		err = append(err, fmt.Sprintf("Tr Description mismatch (%s should be %s)", in.Description, out.Description))
	}

	return
}

func cmpEv(in Ev, out Ev) (err []string) {
	// Id check
	if in.Id != out.Id {
		err = append(err, fmt.Sprintf("Ev Id mismatch (%d should be %d)", in.Id, out.Id))
	}

	// date check
	if !in.Date.Truncate(24 * time.Hour).Equal(out.Date.Truncate(24 * time.Hour)) {
		err = append(err, fmt.Sprintf("Ev Date mismatch (%s should be %s)", in.Date, out.Date))
	}

	// name check
	if in.Name != out.Name {
		err = append(err, fmt.Sprintf("Ev Name mismatch (%s should be %s)", in.Name, out.Name))
	}

	// amount check
	if in.Amount != out.Amount {
		err = append(err, fmt.Sprintf("Ev Amount mismatch (%f should be %f)", in.Amount, out.Amount))
	}

	// type check
	if in.Type != out.Type {
		err = append(err, fmt.Sprintf("Ev Type mismatch (%d should be %d)", in.Type, out.Type))
	}

	// description check
	if in.Description != out.Description {
		err = append(err, fmt.Sprintf("Ev Description mismatch (%s should be %s)", in.Description, out.Description))
	}

	return
}

func cmpCache(in Cache, out Cache) (err []string) {
	if in.Treasury != out.Treasury {
		err = append(err, fmt.Sprintf("Treasury mismatch (%f should be %f)", in.Treasury, out.Treasury))
	}

	err = append(err, cmpPeriod(in.Month, out.Month)...)
	err = append(err, cmpPeriod(in.Year, out.Year)...)

	return
}

func cmpPeriod(in Period, out Period) (err []string) {
	// expenses check
	if in.Expenses != out.Expenses {
		err = append(err, fmt.Sprintf("Expenses mismatch (%f should be %f)", in.Expenses, out.Expenses))
	}

	// income check
	if in.Income != out.Income {
		err = append(err, fmt.Sprintf("Income mismatch (%f should be %f)", in.Income, out.Income))
	}

	// total check
	if in.Total != out.Total {
		err = append(err, fmt.Sprintf("Total mismatch (%f should be %f)", in.Total, out.Total))
	}

	return
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
	var out strings.Builder
	s := NewStats()
	sout, _ := json.Marshal(s.Cache)

	testCases := []OutputCase{
		OutputCase{
			s,
			string(sout),
		},
	}

	for _, tc := range testCases {
		Output(&out, s)

		if out.String() != tc.Output {
			t.Fatalf("expected: %s | got %s", tc.Output, out.String())
		}
	}
}

func TestNewStats(t *testing.T) {
	now := time.Now()

	in := NewStats()
	in.LastCheck = now

	out := Stats{
		[]Tr{},
		[]Ev{},
		Cache{
			0.0,
			Period{
				0.0,
				0.0,
				0.0,
			},
			Period{
				0.0,
				0.0,
				0.0,
			},
		},
		now,
		0,
	}

	errs := cmpStats(in, out)

	for _, err := range errs {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	now := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

	testCases := []UpdateCase{
		/* tc 0 */
		UpdateCase{
			NewStats(),
			NewStats(),
			true,
		},
		/* ********** */

		/* tc 1 */
		UpdateCase{
			Stats{
				[]Tr{
					Tr{
						0,
						"foo",
						time.Now(),
						100.0,
						0,
						"bar",
					},
				},
				[]Ev{},
				Cache{
					100.0,
					Period{},
					Period{},
				},
				time.Now(),
				0,
			},
			Stats{
				[]Tr{
					Tr{
						0,
						"foo",
						time.Now(),
						100.0,
						0,
						"bar",
					},
				},
				[]Ev{},
				Cache{
					100.0,
					Period{},
					Period{},
				},
				time.Now(),
				0,
			},
			true,
		},
		/* ********** */

		/* tc 2 */
		UpdateCase{
			Stats{
				[]Tr{},
				[]Ev{
					Ev{
						0,
						"foo",
						now,
						1,
						[3]int{0, 0, 1},
						100.0,
						0,
						"bar",
					},
				},
				Cache{},
				time.Now(),
				0,
			},
			Stats{
				[]Tr{
					Tr{
						0,
						"foo",
						now,
						100.0,
						0,
						"bar",
					},
				},
				[]Ev{},
				Cache{
					100.0,
					Period{},
					Period{},
				},
				time.Now(),
				1,
			},
			true,
		},
		/* ********** */

		/* tc 3 */
		UpdateCase{
			Stats{
				[]Tr{
					Tr{
						0,
						"test",
						time.Date(2019, 12, 12, 0, 0, 0, 0, time.UTC),
						200.0,
						0,
						"test",
					},
				},
				[]Ev{
					Ev{
						1,
						"foo",
						now,
						2,
						[3]int{0, 0, 1},
						100.0,
						0,
						"bar",
					},
				},
				Cache{},
				time.Now(),
				1,
			},
			Stats{
				[]Tr{
					Tr{
						0,
						"test",
						time.Date(2019, 12, 12, 0, 0, 0, 0, time.UTC),
						200.0,
						0,
						"test",
					},
					Tr{
						1,
						"foo",
						now,
						100.0,
						0,
						"bar",
					},
				},
				[]Ev{
					Ev{
						1,
						"foo",
						now,
						1,
						[3]int{0, 0, 1},
						100.0,
						0,
						"bar",
					},
				},
				Cache{
					300.0,
					Period{},
					Period{},
				},
				time.Now(),
				2,
			},
			true,
		},
		/* ********** */
	}

	for i, tc := range testCases {
		err := Update(&tc.Input)

		if err != nil {
			if tc.Success {
				t.Fatalf("TC %d returned with error %s", i, err)
			}

			return
		}

		for _, errs := range cmpStats(tc.Input, tc.Output) {
			t.Errorf("TC %d -> %s", i, errs)
		}
	}
}

func TestProcess(t *testing.T) {
	base := NewStats()

	testCases := []ProcessCase{
		/* tc 0 */
		ProcessCase{
			[]string{
				INPUT,
				"foo",
				"2020-01-01",
				"100.02",
				"0",
				"bar",
			},
			base,
			Stats{
				[]Tr{
					Tr{
						0,
						"foo",
						time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
						100.02,
						0,
						"bar",
					},
				},
				[]Ev{},
				Cache{
					100.02,
					Period{},
					Period{},
				},
				base.LastCheck,
				1,
			},
			true,
		},
		/* **** */

		/* tc 1 */
		ProcessCase{
			[]string{
				PERIOD,
				"foo",
				"2020-01-02",
				"-1",
				"0,0,0",
				"20",
				"0",
				"bar",
			},
			base,
			Stats{
				[]Tr{},
				[]Ev{
					Ev{
						0,
						"foo",
						time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
						-1,
						[3]int{0, 0, 0},
						20,
						0,
						"bar",
					},
				},
				Cache{},
				base.LastCheck,
				1,
			},
			true,
		},
		/* **** */

		/* tc 2 */
		ProcessCase{
			[]string{
				PERIOD,
				"foo",
				"2019-05-19",
				"-1",
				"0,0,1",
				"20",
				"0",
				"bar",
			},
			Stats{
				[]Tr{
					Tr{
						0,
						"foo0",
						time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
						10.0,
						0,
						"bar0",
					},
					Tr{
						1,
						"foo1",
						time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
						10.1,
						0,
						"bar1",
					},
				},
				[]Ev{
					Ev{
						2,
						"foo2",
						time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
						1,
						[3]int{0, 0, 1},
						30.5,
						0,
						"bar2",
					},
					Ev{
						3,
						"foo3",
						time.Now().AddDate(1, 0, 0),
						1,
						[3]int{0, 0, 1},
						1000,
						0,
						"bar3",
					},
				},
				Cache{},
				time.Now(),
				4,
			},
			Stats{
				[]Tr{
					Tr{
						0,
						"foo0",
						time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
						10.0,
						0,
						"bar0",
					},
					Tr{
						1,
						"foo1",
						time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
						10.1,
						0,
						"bar1",
					},
					Tr{
						4,
						"foo2",
						time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
						30.5,
						0,
						"bar2",
					},
				},
				[]Ev{
					Ev{
						3,
						"foo3",
						time.Now().AddDate(1, 0, 0),
						1,
						[3]int{0, 0, 1},
						1000,
						0,
						"bar3",
					},
					Ev{
						5,
						"foo",
						time.Date(2019, time.May, 19, 0, 0, 0, 0, time.UTC),
						-1,
						[3]int{0, 0, 1},
						20,
						0,
						"bar",
					},
				},
				Cache{
					50.6,
					Period{},
					Period{},
				},
				time.Now(),
				6,
			},
			true,
		},
		/* **** */

		/* tc 3 */
		ProcessCase{
			[]string{
				"foo",
			},
			Stats{},
			Stats{},
			false,
		},
		/* **** */

		/* tc 4 */
		ProcessCase{
			[]string{},
			Stats{},
			Stats{},
			false,
		},
		/* **** */
	}

	for i, tc := range testCases {
		err := Process(tc.Input, &tc.InState)

		if err != nil {
			if tc.Success {
				t.Fatalf("TC %d returned with error %s", i, err)
			}
			continue
		}

		for _, errs := range cmpStats(tc.InState, tc.OutState) {
			t.Errorf("TC %d -> %s", i, errs)
		}
	}
}
