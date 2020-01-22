package stats

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

type Stats struct {
	Transactions []Tr
	Events       []Ev
	Cache        Cache
	LastCheck    time.Time
	Index        uint
}

type Tr struct { // transactions
	Id          uint
	Name        string
	Date        time.Time // date the money was added/subtracted
	Amount      float64
	Type        uint
	Description string
}

type Ev struct { // events
	/* Events are money movements that will happen in the future
	* after the date has passed this events gets deleted from the list and
	* the information goes to a transaction */
	Id          uint
	Name        string
	Date        time.Time // date money will be added/subtracted
	Times       int       // times this will repeat (-1 is indefinite)
	Step        [3]int    // time step for the next repetition (if times is -1 this is ignored)
	Amount      float64
	Type        uint
	Description string
}

type Cache struct {
	Treasury float64
	Month    Period
	Year     Period
}

type Period struct {
	Total    float64
	Income   float64
	Expenses float64
}

// input commands names
const (
	INPUT  = "in"
	PERIOD = "pe"
)

func NewStats() Stats {
	return Stats{
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
		time.Now(),
		0,
	}
}

func Parse(in io.Reader) (out []string, err error) {
	r := csv.NewReader(in)
	r.Comma = ' '

	return r.Read()
}

func Update(s *Stats) error {
	return processInput(nil, s)
}

func Process(in []string, s *Stats) error {
	return processInput(in, s)
}

func processInput(in []string, s *Stats) error {
	nTrans := checkEvents(s)

	if in != nil {
		err := runCmd(in, s)
		if err != nil {
			return err
		}
	}

	if nTrans > 0 || in != nil {
		err := updateCache(s)
		if err != nil {
			return err
		}
	}

	s.LastCheck = time.Now()

	return nil
}

func checkEvents(s *Stats) (n int) {
	now := time.Now()
	remaining := make([]Ev, 0)

	for _, event := range s.Events {
		if now.Before(event.Date) {
			remaining = append(remaining, event)
			continue
		}

		newTr := Tr{
			s.Index,
			event.Name,
			event.Date,
			event.Amount,
			event.Type,
			event.Description,
		}
		s.Transactions = append(s.Transactions, newTr)
		s.Index++
		n++

		event.Times--

		if event.Times > 0 {
			remaining = append(remaining, Ev{
				event.Id,
				event.Name,
				event.Date.AddDate(event.Step[0], event.Step[1], event.Step[2]),
				event.Times,
				event.Step,
				event.Amount,
				event.Type,
				event.Description,
			})
		}
	}

	if len(remaining) < len(s.Events) {
		s.Events = remaining
	}

	return
}

func runCmd(in []string, s *Stats) (err error) {
	if len(in) < 1 {
		return errors.New("no command")
	}

	switch in[0] {
	case INPUT:
		err = input(in[1:], s)
	case PERIOD:
		err = period(in[1:], s)
	default:
		return errors.New("invalid command")
	}

	return
}

func updateCache(s *Stats) (err error) {
	total := 0.0

	for _, tr := range s.Transactions {
		total += tr.Amount
	}

	s.Cache.Treasury = total

	return nil
}

// input name date amount type description
func input(args []string, s *Stats) error {
	if len(args) < 5 {
		return errors.New("inputcmd: missing arguments")
	}

	// index
	index := s.Index
	s.Index++

	// name
	name := args[0]

	// parse date
	date, err := time.Parse("2006-01-02", args[1])
	if err != nil {
		return err
	}

	// parse amount
	amount, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return err
	}

	// parse type
	t, err := strconv.ParseUint(args[3], 10, 32)
	if err != nil {
		return err
	}

	// description
	description := args[4]

	s.Transactions = append(s.Transactions, Tr{
		index,
		name,
		date,
		amount,
		uint(t),
		description,
	})

	return nil
}

// period name date times year,month,day amount type description
func period(args []string, s *Stats) error {
	if len(args) < 7 {
		return errors.New("periodcmd: missing arguments")
	}

	// index
	index := s.Index
	s.Index++

	// name
	name := args[0]

	// parse date
	date, err := time.Parse("2006-01-02", args[1])
	if err != nil {
		return err
	}

	// parse times
	times, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		return err
	}

	// parse step
	var step [3]int

	for i, stepStr := range strings.Split(args[3], ",") {
		s, err := strconv.ParseInt(stepStr, 10, 32)
		if err != nil {
			return err
		}

		step[i] = int(s)
	}

	// parse amount
	amount, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return err
	}

	// parse type
	t, err := strconv.ParseUint(args[5], 10, 32)
	if err != nil {
		return err
	}

	// description
	description := args[6]

	s.Events = append(s.Events, Ev{
		index,
		name,
		date,
		int(times),
		step,
		amount,
		uint(t),
		description,
	})
	return nil
}

func Output(out io.Writer, s Stats) error {
	encoded, err := json.Marshal(s.Cache)
	if err != nil {
		return err
	}

	_, err = out.Write(encoded)
	return err
}
