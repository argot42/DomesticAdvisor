package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/argot42/DomesticAdvisor/config"
	"github.com/argot42/DomesticAdvisor/stats"
	"github.com/argot42/watcher"
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
        log.Fatalln("setup:", err)
    }

    // sigterm setup
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGTERM)

    // timer setup
    timer := make(chan stats.Timer, 10)

    if err = start(ctl, timer, status, cfg.Timeout, sigs); err != nil {
        log.Fatalln("runtime:", err)
    }

    // cleaning
    log.Println("Closing")
    status.Close()
    ctl.Done <- true
    // wait for the goroutines to end
    log.Println("Wating for goroutines to finish")
    <-ctl.Done
    log.Println("bye :)")
}

func setupFiles(cfg *config.Config) (ctl watcher.R, status *os.File, err error) {
    // watch control file
    ctl, e := watcher.Read(cfg.CtlFilePath)
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

func start(ctl watcher.R, timer chan stats.Timer, status *os.File, timeout time.Duration, sigs chan os.Signal) error {
    /* state */
    var transactions []stats.Transaction
    var events []stats.Event

    /********/
    var buffer []byte

    End:
    for {
        select {
        case input := <-ctl.Out:
            // initialize state
            if input.First {
                transactions = make([]stats.Transaction, 0, 5)
                events = make([]stats.Event, 0, 5)
            }

            if input.Data != 10 { // 10 is newline
                buffer = append(buffer, input.Data)
                continue
            }

            // parse input
            parsed, err := stats.Parse(bytes.NewReader(buffer))
            if err != nil {
                return fmt.Errorf("stats: %s", err)
            }
            // clean buffer
            buffer = nil

            // process input
            switch(parsed[0]) {
            case "tr":
                tr, err := stats.ProcessTransaction(parsed)
                if err != nil {
                    log.Println(err)
                    continue
                }

                transactions = append(transactions, tr)
            case "ev":
                ev, err := stats.ProcessEvent(parsed)
                if err != nil {
                    log.Println(err)
                    continue
                }

                events = append(events, ev)
                stats.StartTimer(ev, time.Now(), timer)
            }

            // update stats
            s := stats.BuildStats(transactions, events)

            if err = stats.UpdateStats(s, status); err != nil {
                return fmt.Errorf("update: %s", err)
            }

        case t := <-timer:
            fmt.Println(t)

        case err := <-ctl.Err:
            return fmt.Errorf("control file erorr: %s", err)

        case <-sigs:
            break End
        }
    }

    return nil
}
