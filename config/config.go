package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	// options names
	transactionslog = "transactionslog"
	ctlfile         = "ctlfile"
	statusfile      = "statusfile"
	// default filenames
	filename_transactions_log = "transactions_log.csv"
	filename_status           = "status.json"
	filename_ctlfile          = "ctl"
)

// names on .ini file
type Config struct {
	TransactionsLogPath string
	StatusPath          string
	CtlFilePath         string
}

// errors
type ErrCfgFormat int

func (e ErrCfgFormat) Error() string {
	return "Error in config file at line " + strconv.Itoa(int(e))
}

var ErrConfigFormat ErrCfgFormat
var ErrConfigFilePath = errors.New("No valid configuration file path provided")

func GetConfig(args []string) (cfg *Config, err error) {
	if len(args) < 2 {
		cfg, err = newConfig()
		return
	}
	file, err := os.Open(args[1])
	if err != nil {
		return nil, ErrConfigFilePath
	}
	defer file.Close()
	cfg, err = parseConfig(bufio.NewReader(file))
	return
}

func newConfig() (cfg *Config, err error) {
	transactionsLogPath, err := filepath.Abs("./" + filename_transactions_log)
	if err != nil {
		return
	}
	statusPath, err := filepath.Abs("./" + filename_status)
	if err != nil {
		return
	}
	ctlFilePath, err := filepath.Abs("./" + filename_ctlfile)
	if err != nil {
		return
	}
	return &Config{
		transactionsLogPath,
		statusPath,
		ctlFilePath,
	}, nil
}

func parseConfig(r *bufio.Reader) (cfg *Config, err error) {
	cfg, err = newConfig()
	lineCount := 0

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			} else if err == io.EOF && len(line) == 0 {
				break
			}
		}
		lineCount++
		err = parseLine(line, cfg)
		if err != nil {
			ErrConfigFormat = ErrCfgFormat(lineCount)
			return nil, ErrConfigFormat
		}
	}

	return
}

func parseLine(line string, cfg *Config) error {
	r, err := regexp.Compile(`^\s*(#.*|(.+)\s*=\s*(.+))?\s*$`)
	if err != nil {
		return err
	}
	// error in formatting
	matches := r.FindStringSubmatch(line)
	if len(matches) == 0 {
		return errors.New("Wrongly formatted config file")
	}
	// comment or blank line
	if len(matches[1]) == 0 || len(matches[2]) == 0 {
		return nil
	}

	switch matches[2] {
	case transactionslog:
		transactionsLogPath, err := filepath.Abs(matches[3])
		if err != nil {
			return err
		}
		cfg.TransactionsLogPath = transactionsLogPath

	case statusfile:
		statusPath, err := filepath.Abs(matches[3])
		if err != nil {
			return err
		}
		cfg.StatusPath = statusPath

	case ctlfile:
		ctlFilePath, err := filepath.Abs(matches[3])
		if err != nil {
			return err
		}
		cfg.CtlFilePath = ctlFilePath

	default:
		return errors.New(matches[2] + " is not a valid option")
	}

	return nil
}

func Usage() {
	fmt.Println("usage:", os.Args[0], "config_file")
}
