package main

import "filepath"

type Config struct {
	StatusPath  string
	CtlFilePath string
}

func GetConfig(args []string) (cfg *Config, err error) {
	return &Config{
		filepath.Abs("./status.json"),
		filepath.Abs("./ctl"),
	}, nil
}
