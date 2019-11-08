package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	StatusPath  string
	CtlFilePath string
	Timeout     time.Duration
}

// errors
type ErrCfgFormat int

func (e ErrCfgFormat) Error() string {
	return "Error in config file at line " + strconv.Itoa(int(e))
}

var ErrConfigFormat ErrCfgFormat
var ErrConfigFilePath = errors.New("No valid configuration file path provided")

func GetConfig(args []string) (cfg *Config, err error) {
	statusPath, _ := filepath.Abs("./status.json")
	ctlFilePath, _ := filepath.Abs("./ctl")

	return &Config{
		statusPath,
		ctlFilePath,
		10,
	}, nil
}

func Usage() {
	fmt.Println("usage:", os.Args[0], "config_file")
}
