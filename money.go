package main

import (
	"fmt"
	"github.com/argot42/money/config"
	"github.com/argot42/money/stats"
	"github.com/argot42/watcher"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

	ctl, status, err := setupFiles(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM)

	err = start(ctl, status, cfg.Timeout, sigs)
	if err != nil {
		log.Fatalln(err)
	}

	// cleaning
	fmt.Println("Closing...")
	status.Close()
	close(ctl.Done)
	// wait for children gorouting to end
	fmt.Println("Waiting for goroutines to finish")
	<-ctl.Out
	fmt.Println("bye :)")
}

func setupFiles(cfg *config.Config) (ctl watcher.Sub, status *os.File, err error) {
	// watch control file
	ctl, e := watcher.Watch(cfg.CtlFilePath)
	if e != nil {
		err = fmt.Errorf("error setting up ctl file: %s", e)
		return
	}

	// create status file
	status, e = os.Create(cfg.StatusPath)
	if e != nil {
		err = fmt.Errorf("error setting up ctl file: %s", e)
		return
	}

	return
}

func start(ctl watcher.Sub, status *os.File, timeout time.Duration, sigs chan os.Signal) error {
	var buffer []byte
	s := stats.NewStats()

End:
	for {
		select {
		case input := <-ctl.Out:
			if input != 10 { // 10 is newline
				buffer = append(buffer, input)
				continue
			}

			// parse input
			parsed, err := stats.Parse(strings.NewReader(string(buffer)))
			if err != nil {
				return err
			}

			err := stats.Process(parsed, &s)
			if err != nil {
				return err
			}
			err = stats.Output(status, s)
			if err != nil {
				return err
			}
			buffer = nil

		case e := <-ctl.Err:
			return fmt.Errorf("error control file: %s", e)

		case <-time.After(timeout * time.Second):
			// force an update
			err := stats.Update(&s)
			if err != nil {
				return err
			}
			err = stats.Output(status, s)
			if err != nil {
				return err
			}

		case <-sigs:
			break End
		}
	}

	return nil
}
