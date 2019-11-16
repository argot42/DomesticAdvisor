package stats

import (
	"errors"
	"os"
	"strings"
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
	Amount      float32
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
	Amount      float32
	Type        uint
	Description string
}

type Cache struct {
	Treasury float32
	Month    Period
	Year     Period
}

type Period struct {
	Total    float32
	Income   float32
	Expenses float32
}

// input commands names
const INPUT string = "in"
const PERIOD string = "pe"

// errors
const parsingCommandErr error = errors.New("Error parsing command")

func NewStats() Stats {
	return Stats{
		[]Transactions{},
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

func Update(s *Stats) (out string, err error) {
	return processInput("", false, s)
}

func Process(in string, s *Stats) (out string, err error) {
	return processInput(in, true, s)
}

func processInput(in string, newInfo bool, s *Stats) (err error) {
	nTrans := checkEvents(s)

	if newInfo {
		err = runCmd(in, s)
		if err != nil {
			return
		}
	}

	err = updateCache(s)
	if err != nil {
		return
	}

	s.LastCheck = time.Now()

	return
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

func runCmd(in string, s *Stats) (int changes, err error) {
	return 0, nil
}

func updateCache(s *Stats) (err error) {
	return nil
}

func Output(out *os.File) error {
	return nil
}
