package main

import (
	"fmt"
	"github.com/argot42/money/config"
	"github.com/argot42/watcher"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	ctl, err := setupFiles(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM)

	err = start(ctl, sigs)
	if err != nil {
		log.Fatalln(err)
	}
}

func setupFiles(cfg *config.Config) (ctl watcher.Sub, err error) {
	// watch control file
	ctl, e := watcher.Watch(cfg.CtlFilePath)
	if e != nil {
		err = fmt.Errorf("error setting up ctl file: %s", e)
		return
	}

	// for writing control file

	return
}

func start(ctl watcher.Sub, sigs chan os.Signal) error {
End:
	for {
		select {
		case input := <-ctl.Out:
			fmt.Println(input)
		case e := <-ctl.Err:
			return fmt.Errorf("error control file: %s", e)
		case <-sigs:
			break End
		}
	}
	return nil
}
