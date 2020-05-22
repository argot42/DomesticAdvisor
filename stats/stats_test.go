package stats

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

type UpdateStatsCase struct {
	Input  Stats
	Output string
}

type BuildTransactionCase struct {
    Input TrCase
    Output Transaction
}

type TrCase struct {
    Name string
    Description string
    Date time.Time
    Amount float64
}

type BuildEventCase struct {
    Input EvCase
    Output Event
}

type EvCase struct {
    Name string
    Description string
    Date time.Time
    Times int
    Step [3]int
    Amount float64
}

func TestParse(t *testing.T) {
	// setup
	parseCases := []ParseCase{
		{
			"foo b a r",
			[]string{"foo", "b", "a", "r"},
			true,
		},
		{
			"foo \"foo bar\"",
			[]string{"foo", "foo bar"},
			true,
		},
		{
			"foo \"foo bar\" bar",
			[]string{"foo", "foo bar", "bar"},
			true,
		},
		{
			"foo",
			[]string{"foo"},
			true,
		},
		{
			"",
			[]string{""},
			true,
		},
		{
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
		{
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
		{
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
		{
			[]string{"Tr"},
			Transaction{},
			false,
		},
		{
			[]string{"foo", "bar", "2020-01-01", "200"},
			Transaction{},
			false,
		},
		{
			[]string{"Tr", "x", "x", "a", "b"},
			Transaction{},
			false,
		},
		{
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
		{
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
		{
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
		{
			[]string{"", "", "", "", "", "", ""},
			Event{},
			false,
		},
		{
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
		{
			[]Transaction{
				{
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
						{
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
		{
			[]Transaction{
				{
					0,
					"foo",
					"bar",
					time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
					10,
				},
				{
					1,
					"bar",
					"",
					time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC),
					5.5,
				},
			},
			[]Event{
				{
					0,
					"event",
					"",
					now,
					1,
					[3]int{0, 0, 1},
					100.101,
				},
				{
					1,
					"event1",
					"",
					now,
					1,
					[3]int{0, 0, 2},
					10.5,
				},
				{
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
						{
							"foo",
							10,
							time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
						},
						{
							"bar",
							5.5,
							time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC),
						},
					},
				},
				Activity{
					110.601,
					[]Entry{
						{
							"event",
							100.101,
							now,
						},
						{
							"event1",
							10.5,
							now,
						},
					},
				},
				Activity{
					-22.1,
					[]Entry{
						{
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

func TestUpdateStats(t *testing.T) {
	usc := []UpdateStatsCase{
		{
			Stats{
				Activity{
					100.4,
					[]Entry{
						{
							"foo",
							100.4,
							time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
						},
					},
				},
				Activity{},
				Activity{},
				0,
			},
			"{\"Treasury\":{\"Total\":100.4,\"Entries\":[{\"Name\":\"foo\",\"Amount\":100.4,\"Date\":\"2020-01-01T00:00:00Z\"}]},\"Income\":{\"Total\":0,\"Entries\":null},\"Expenses\":{\"Total\":0,\"Entries\":null},\"Balance\":0}",
		},
	}

	for i, c := range usc {
		// create temporal file
		tmp, err := ioutil.TempFile("", fmt.Sprintf("%d_update-stats_", i))
		if err != nil {
			t.Fatalf("Tmp file: %s", err)
		}

		// write to file
		err = UpdateStats(c.Input, tmp)
		if err != nil {
			t.Errorf("%d failed: %s", i, err)
		}

		fileName := tmp.Name()

		// close writing file
		if err = tmp.Close(); err != nil {
			t.Errorf("%d closing %s: %s", i, fileName, err)
		}

		// open same file for reading
		readTmp, err := os.Open(fileName)
		if err != nil {
			t.Errorf("%d opening %s: %s", i, fileName, err)
			continue
		}

		// read all file
		content, err := ioutil.ReadAll(readTmp)
		if err != nil {
			t.Errorf("%d reading %s: %s", i, fileName, err)
			continue
		}

		// compare content with output
		if string(content) != c.Output {
			t.Errorf("%d got %s but should be %s", i, string(content), c.Output)
		}

		// close file for reading
		if err = readTmp.Close(); err != nil {
			t.Errorf("%d reading file close %s: %s", i, fileName, err)
		}

		// remove file
		if err = os.Remove(fileName); err != nil {
			t.Errorf("%d removing file %s: %s", i, fileName, err)
		}
	}
}

func TestStartTimer(t *testing.T) {
	ev := Event{
		0,
		"foo",
		"bar",
		time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
		0,
		[3]int{0, 0, 0},
		230.10,
	}
	now := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	out := make(chan Timer, 5)

	StartTimer(ev, now, out)

	select {
	case timer := <-out:
		if ev.Id != timer.Id {
			t.Fatalf("%d -> %d", timer.Id, ev.Id)
		}

	case <-time.After(5 * time.Second):
		t.Fatal("Timeout!")
	}
}

func TestBuildTransactions(t * testing.T) {
    TRINDEX = 0
    now := time.Now()

    cases := []TrCase{
        {
            "foo",
            "bar",
            now,
            200,
        },
        {
            "bar",
            "foo",
            now,
            130.9,
        },
        {
            "a",
            "b",
            now,
            11,
        },
    }

    // build actual TCs
    actualTCs := make([]BuildTransactionCase, 0, len(cases))

    for i, c := range cases {
        tr := Transaction{
            uint(i),
            c.Name,
            c.Description,
            c.Date,
            c.Amount,
        }

        actualTCs = append(actualTCs, BuildTransactionCase{c, tr})
    }

    for i, tc := range actualTCs {
        tr := BuildTransaction(
            tc.Input.Name,
            tc.Input.Description,
            tc.Input.Date,
            tc.Input.Amount,
        )

        if tr.Id != tc.Output.Id {
            t.Errorf("TC %d: got id %d and should be %d", i, tr.Id, tc.Output.Id)
        }
        if tr.Name != tc.Output.Name {
            t.Errorf("TC %d: got name %s and should be %s", i, tr.Name, tc.Output.Name)
        }
        if tr.Description != tc.Output.Description {
            t.Errorf("TC %d: got description %s and should be %s", i, tr.Description, tc.Output.Description)
        }
        if !tr.Date.Equal(tc.Output.Date) {
            t.Errorf("TC %d: got date %s and should be %s", i, tr.Date, tc.Output.Date)
        }
        if tr.Amount != tc.Output.Amount {
            t.Errorf("TC %d: got amount %f and should be %f", i, tr.Amount, tc.Output.Amount)
        }
    }
}

func TestBuildEvent(t *testing.T) {
    EVINDEX = 0
    now := time.Now()

    cases := []EvCase {
        {
            "foo",
            "bar",
            now,
            2,
            [3]int{1, 2, 3},
            200,
        },
        {
            "bar",
            "foo",
            now,
            -1,
            [3]int{0, 0, 1},
            20,
        },
        {
            "a",
            "x",
            now,
            1,
            [3]int{0, 0, 0},
            2000,
        },
    }

    // build actual TCs
    actualTCs := make([]BuildEventCase, 0, len(cases))

    for i, c := range cases {
        ev := Event{
            uint(i),
            c.Name,
            c.Description,
            c.Date,
            c.Times,
            c.Step,
            c.Amount,
        }

        actualTCs = append(actualTCs, BuildEventCase{c, ev})
    }

    for i, tc := range actualTCs {
        ev := BuildEvent(
            tc.Input.Name,
            tc.Input.Description,
            tc.Input.Date,
            tc.Input.Times,
            tc.Input.Step,
            tc.Input.Amount,
        )

        if ev.Id != tc.Output.Id {
            t.Errorf("TC %d: got id %d and should be %d", i, ev.Id, tc.Output.Id)
        }
        if ev.Name != tc.Output.Name {
            t.Errorf("TC %d: got name %s and should be %s", i, ev.Name, tc.Output.Name)
        }
        if ev.Description != tc.Output.Description {
            t.Errorf("TC %d: got description %s and should be %s", i, ev.Description, tc.Output.Description)
        }
        if !ev.Date.Equal(tc.Output.Date) {
            t.Errorf("TC %d: got date %s and should be %s", i, ev.Date, tc.Output.Date)
        }
        if ev.Times != tc.Output.Times {
            t.Errorf("TC %d: got times %d and should be %d", i, ev.Times, tc.Output.Times)
        }
        for j := 0; j < 3; j++ {
            if ev.Step[j] != tc.Output.Step[j] {
                t.Errorf("TC %d: got step[%d] %d and should be %d", i, j, ev.Step[j], tc.Output.Step[j])
            }
        }
        if ev.Amount != tc.Output.Amount {
            t.Errorf("TC %d: got amount %f and should be %f", i, ev.Amount, tc.Output.Amount)
        }
    }
}
