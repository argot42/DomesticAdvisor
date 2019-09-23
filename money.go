package main

import (
	"fmt"
	"github.com/argot42/money/config"
	"github.com/argot42/watcher"
	"log"
	"os"
)

func main() {
	cfg, err := config.GetConfig(os.Args)
	if err != nil {
		if err == config.ErrConfigFilePath {
			config.Usage()
			return
		}
		log.Fatalln("config:", err)
	}

	ctl, transaction, status, err := setup(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	err = start(ctl, transaction, status)
	if err != nil {
		log.Fatalln(err)
	}
}

func setup(cfg *config.Config) (ctl watcher.Sub, transaction watcher.Sub, status watcher.Sub, err error) {
	ctl, e := watcher.Watch(cfg.CtlFilePath)
	if e != nil {
		err = fmt.Errorf("error setting up ctl file: %s", e)
		return
	}
	transaction, e = watcher.Watch(cfg.TransactionsLogPath)
	if e != nil {
		err = fmt.Errorf("error setting up transaction file: %s", e)
		return
	}
	status, e = watcher.Watch(cfg.StatusPath)
	if e != nil {
		err = fmt.Errorf("error setting up status file: %s", e)
		return
	}

	return
}

func start(ctl watcher.Sub, transaction watcher.Sub, status watcher.Sub) error {
	return fmt.Errorf("todo: everything")
}
