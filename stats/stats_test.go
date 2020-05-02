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

type ProcessTransactionCase struct {
	Input   []string
	Output  Transaction
	Success bool
}

type ProcessEventCase struct {
	Input   []string
	Output  Event
	Success bool
}

type BuildStatsCase struct {
	TrInput []Transaction
	EvInput []Event
	Output  Stats
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

func TestProcessTransaction(t *testing.T) {
	ptc := []ProcessTransactionCase{
		ProcessTransactionCase{
			[]string{"Tr", "foo", "bar", "2020-01-01", "100"},
			Transaction{
				0,
				"foo",
				"bar",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				100,
			},
			true,
		},
		ProcessTransactionCase{
			[]string{"Tr", "", "", "2021-02-03", "30.10"},
			Transaction{
				1,
				"",
				"",
				time.Date(2021, 2, 3, 0, 0, 0, 0, time.UTC),
				30.10,
			},
			true,
		},
		ProcessTransactionCase{
			[]string{"Tr"},
			Transaction{},
			false,
		},
		ProcessTransactionCase{
			[]string{"foo", "bar", "2020-01-01", "200"},
			Transaction{},
			false,
		},
		ProcessTransactionCase{
			[]string{"Tr", "x", "x", "a", "b"},
			Transaction{},
			false,
		},
		ProcessTransactionCase{
			[]string{"", "", "", "", "", "", "", "", ""},
			Transaction{},
			false,
		},
	}

	for i, c := range ptc {
		failed := false
		tr, err := ProcessTransaction(c.Input)

		if err != nil {
			failed = true
			if c.Success {
				t.Errorf("%d: failed %s\n", i, err)
				continue
			}
		}

		if tr.Id != c.Output.Id {
			failed = true
			if c.Success {
				t.Errorf("%d: Id -> %d and should be %d", i, tr.Id, c.Output.Id)
			}
		}

		if tr.Name != c.Output.Name {
			failed = true
			if c.Success {
				t.Errorf("%d: Name -> %s and should be %s", i, tr.Name, c.Output.Name)
			}
		}

		if tr.Description != c.Output.Description {
			failed = true
			if c.Success {
				t.Errorf("%d: Description -> %s and should be %s", i, tr.Description, c.Output.Description)
			}
		}

		if !tr.Date.Equal(c.Output.Date) {
			failed = true
			if c.Success {
				t.Errorf("%d: Date -> %s and should be %s", i, tr.Date, c.Output.Date)
			}
		}

		if tr.Amount != c.Output.Amount {
			failed = true
			if c.Success {
				t.Errorf("%d: Amount -> %f and should be %f", i, tr.Amount, c.Output.Amount)
			}
		}

		if !c.Success && !failed {
			t.Errorf("%d: should have failed but it did not", i)
		}
	}
}

func TestProcessEvent(t *testing.T) {
	tpe := []ProcessEventCase{
		ProcessEventCase{
			[]string{"Ev", "foo", "bar", "2020-10-10", "-1", "0,0,0", "2020"},
			Event{
				0,
				"foo",
				"bar",
				time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC),
				-1,
				[3]int{0, 0, 0},
				2020,
			},
			true,
		},
		ProcessEventCase{
			[]string{"Ev", "", "", "2100-01-02", "10", "1,2,3", "2"},
			Event{
				1,
				"",
				"",
				time.Date(2100, 1, 2, 0, 0, 0, 0, time.UTC),
				10,
				[3]int{1, 2, 3},
				2,
			},
			true,
		},
		ProcessEventCase{
			[]string{"", "", "", "", "", "", ""},
			Event{},
			false,
		},
		ProcessEventCase{
			[]string{"Ev", "", "", "", "10", "1,2,0", "100"},
			Event{},
			false,
		},
	}

	for i, c := range tpe {
		failed := false
		ev, err := ProcessEvent(c.Input)

		if err != nil {
			failed = true
			if c.Success {
				t.Errorf("%d: failed: %s", i, err)
				continue
			}
		}

		if ev.Id != c.Output.Id {
			failed = true
			if c.Success {
				t.Errorf("%d: Id -> %d should be %d", i, ev.Id, c.Output.Id)
			}
		}

		if ev.Name != c.Output.Name {
			failed = true
			if c.Success {
				t.Errorf("%d: Name -> %s should be %s", i, ev.Name, c.Output.Name)
			}
		}

		if ev.Description != c.Output.Description {
			failed = true
			if c.Success {
				t.Errorf("%d: Description -> %s should be %s", i, ev.Description, c.Output.Description)
			}
		}

		if !ev.Date.Equal(c.Output.Date) {
			failed = true
			if c.Success {
				t.Errorf("%d: Date -> %s should be %s", i, ev.Date, c.Output.Date)
			}
		}

		if ev.Times != c.Output.Times {
			failed = true
			if c.Success {
				t.Errorf("%d: Times -> %d should be %d", i, ev.Times, c.Output.Times)
			}
		}

		if ev.Step[0] != c.Output.Step[0] || ev.Step[1] != c.Output.Step[1] || ev.Step[2] != c.Output.Step[2] {
			failed = true
			if c.Success {
				t.Errorf("%d: Step -> %v should be %v", i, ev.Step, c.Output.Step)
			}
		}

		if ev.Amount != c.Output.Amount {
			failed = true
			if c.Success {
				t.Errorf("%d: Amount -> %f should be %f", i, ev.Amount, c.Output.Amount)
			}
		}

		if !c.Success && !failed {
			t.Errorf("%d: should have failed", i)
		}
	}
}

func TestBuildStats(t *testing.T) {
	now := time.Now()

	bsc := []BuildStatsCase{
		// tc0
		BuildStatsCase{
			[]Transaction{
				Transaction{
					0,
					"foo",
					"bar",
					now,
					200.10,
				},
			},
			[]Event{},
			Stats{
				Activity{
					200.10,
					[]Entry{
						Entry{
							"foo",
							200.10,
							now,
						},
					},
				},
				Activity{},
				Activity{},
				0,
			},
			true,
		},

		// tc1
		BuildStatsCase{
			[]Transaction{
				Transaction{
					0,
					"foo",
					"bar",
					time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
					10,
				},
				Transaction{
					1,
					"bar",
					"",
					time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC),
					5.5,
				},
			},
			[]Event{
				Event{
					0,
					"event",
					"",
					now,
					1,
					[3]int{0, 0, 1},
					100.101,
				},
				Event{
					1,
					"event1",
					"",
					now,
					1,
					[3]int{0, 0, 2},
					10.5,
				},
				Event{
					2,
					"event2",
					"",
					now,
					2,
					[3]int{0, 0, 3},
					-22.1,
				},
			},
			Stats{
				Activity{
					15.5,
					[]Entry{
						Entry{
							"foo",
							10,
							time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
						},
						Entry{
							"bar",
							5.5,
							time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC),
						},
					},
				},
				Activity{
					110.601,
					[]Entry{
						Entry{
							"event",
							100.101,
							now,
						},
						Entry{
							"event1",
							10.5,
							now,
						},
					},
				},
				Activity{
					-22.1,
					[]Entry{
						Entry{
							"event2",
							-22.1,
							now,
						},
					},
				},
				88.501,
			},
			true,
		},
	}

	for i, c := range bsc {
		failed := false
		s := BuildStats(c.TrInput, c.EvInput)

		if !checkActivity(s.Treasury, c.Output.Treasury, c.Success, "treasury", i, t) {
			failed = true
		}
		if !checkActivity(s.Income, c.Output.Income, c.Success, "income", i, t) {
			failed = true
		}
		if !checkActivity(s.Expenses, c.Output.Expenses, c.Success, "expenses", i, t) {
			failed = true
		}

		if s.Balance != c.Output.Balance {
			failed = true
			if c.Success {
				t.Errorf("%d: balance is %f but should be %f", i, s.Balance, c.Output.Balance)
			}
		}

		if !c.Success && !failed {
			t.Errorf("%d: didn't fail", i)
		}
	}
}

func checkActivity(a0, a1 Activity, success bool, name string, i int, t *testing.T) bool {
	failed := true

	if a0.Total != a1.Total {
		failed = false
		if success {
			t.Errorf("%d: %s -> total is %f but should be %f", i, name, a0.Total, a1.Total)
		}
	}

	for j, e0 := range a0.Entries {
		e1 := a1.Entries[j]

		if e0.Name != e1.Name {
			failed = false
			if success {
				t.Errorf("%d - %d: %s -> name is %s but should be %s", i, j, name, e0.Name, e1.Name)
			}
		}

		if e0.Amount != e1.Amount {
			failed = false
			if success {
				t.Errorf("%d - %d: %s -> amount is %f but should be %f", i, j, name, e0.Amount, e1.Amount)
			}
		}

		if !e0.Date.Equal(e1.Date) {
			failed = false
			if success {
				t.Errorf("%d - %d: %s -> date is %s but should be %s", i, j, name, e0.Date, e1.Date)
			}
		}
	}

	return failed
}
