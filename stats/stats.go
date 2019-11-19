package stats

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"time"
)

type Stats struct {
	Transactions []Tr
	Events       []Ev
	Cache        Cache
	LastCheck    time.Time
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
	Date        time.Time // date money will be added/subtracted
	Name        string
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
const INPUT string = "in"
const PERIOD string = "pe"

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

func checkEvents(s *Stats) (nTrans int) {
	now := time.Now()
	nTrans = 0

	for _, event := range s.Events {
		if now.Before(event.Date) {
			continue
		}

		newTr := Tr{
			0, // gen id
			event.Name,
			event.Date,
			event.Amount,
			event.Type,
			event.Description,
		}
		s.Transactions = append(s.Transactions, newTr)
		nTrans++
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

func input(args []string, s *Stats) error {
	// TODO
	return nil
}

func period(args []string, s *Stats) error {
	// TODO
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
