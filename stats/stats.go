package stats

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var TRINDEX uint
var EVINDEX uint

/* -- output -- */
type Stats struct {
	Treasury Activity
	Income   Activity
	Expenses Activity
	Balance  float64
}

type Activity struct {
	Total   float64
	Entries []Entry
}

type Entry struct {
	Name   string
	Amount float64
	Date   time.Time
}

/* -------------- */

/* -- state -- */
type Transaction struct {
	Id          uint
	Name        string
	Description string
	Date        time.Time
	Amount      float64
}

type Event struct {
	Id          uint
	Name        string
	Description string
	Date        time.Time // date the amount will be added/subtracted
	Times       int       // times this event will repeat (-1 is indefinite)
	Step        [3]int    // time step for next repetition (if times is 0 this is ignored)
	Amount      float64
}

type Timer struct {
	Id   uint
	Date time.Time
}

/* -------------- */

func Parse(in io.Reader) ([]string, error) {
	r := csv.NewReader(in)
	r.Comma = ' '

	return r.Read()
}

func ProcessTransaction(in []string) (Transaction, error) {
	if len(in) < 5 {
		return Transaction{}, fmt.Errorf("process transaction: missing arguments")
	}

	// index
	index := TRINDEX
	TRINDEX++

	// name
	name := in[1]

	// description
	description := in[2]

	// parse date
	date, err := time.Parse("2006-01-02", in[3])
	if err != nil {
		return Transaction{}, fmt.Errorf("process transaction: %s", err)
	}

	// parse amount
	amount, err := strconv.ParseFloat(in[4], 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("process transaction: %s", err)
	}

	return Transaction{
		index,
		name,
		description,
		date,
		amount,
	}, nil
}

func ProcessEvent(in []string) (Event, error) {
	if len(in) < 7 {
		return Event{}, fmt.Errorf("process event: missing arguments")
	}

	// index
	index := EVINDEX
	EVINDEX++

	// name
	name := in[1]

	// description
	description := in[2]

	// parse date
	date, err := time.Parse("2006-01-02", in[3])
	if err != nil {
		return Event{}, fmt.Errorf("process event: %s", err)
	}

	// parse times
	times, err := strconv.ParseInt(in[4], 10, 32)
	if err != nil {
		return Event{}, fmt.Errorf("process event: %s", err)
	}

	// parse step
	var step [3]int

	for i, stepStr := range strings.Split(in[5], ",") {
		s, err := strconv.ParseInt(stepStr, 10, 32)
		if err != nil {
			return Event{}, fmt.Errorf("process event: %s", err)
		}

		step[i] = int(s)
	}

	// parse amount
	amount, err := strconv.ParseFloat(in[6], 64)
	if err != nil {
		return Event{}, fmt.Errorf("process event: %s", err)
	}

	return Event{
		index,
		name,
		description,
		date,
		int(times),
		step,
		amount,
	}, nil
}

func BuildStats(Transactions []Transaction, Events []Event) (stats Stats) {
	for _, tr := range Transactions {
		stats.Treasury.Total += tr.Amount
		stats.Treasury.Entries = append(stats.Treasury.Entries, Entry{
			tr.Name,
			tr.Amount,
			tr.Date,
		})
	}

	for _, ev := range Events {
		if !inMonth(ev.Date, time.Now()) {
			continue
		}

		stats.Balance += ev.Amount

		if ev.Amount >= 0 {
			stats.Income.Total += ev.Amount
			stats.Income.Entries = append(stats.Income.Entries, Entry{
				ev.Name,
				ev.Amount,
				ev.Date,
			})
		} else {
			stats.Expenses.Total += ev.Amount
			stats.Expenses.Entries = append(stats.Expenses.Entries, Entry{
				ev.Name,
				ev.Amount,
				ev.Date,
			})
		}
	}

	return
}

func inMonth(d, month time.Time) bool {
	return d.Year() == month.Year() && d.Month() == month.Month()
}

func UpdateStats(s Stats, f *os.File) error {
	serialized, err := json.Marshal(s)
	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Seek(0, 0)
	if _, err := f.Write(serialized); err != nil {
		return err
	}

	return nil
}

func StartTimer(ev Event, now time.Time, timer chan<- Timer) {
	duration := ev.Date.Sub(now)

	go func() {
		t := <-time.After(duration)
		timer <- Timer{
			ev.Id,
			t,
		}
	}()
}
