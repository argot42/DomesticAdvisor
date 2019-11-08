package info

import (
	"os"
)

type Stats struct{}

func NewStats() Stats {
	return Stats{}
}

func Update(s *Stats) (out string, err error) {
	return processInput("", false, s)
}

func Process(in string, s *Stats) (out string, err error) {
	return processInput(in, true, s)
}

func processInput(in string, newInfo bool, s *Stats) (out string, err error) {
	return "", nil
}

func Output(in string, out *os.File) error {
	return nil
}
