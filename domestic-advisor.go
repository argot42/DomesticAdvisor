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
    for {
        _, open := <-ctl.Out
        if !open {
            break
        }
    }
    log.Println("bye :)")
}

func setupFiles(cfg *config.Config) (ctl watcher.R, status *os.File, err error) {
    // watch control file
    ctl = watcher.Read(cfg.CtlFilePath)

    // create status file
    status, e := os.Create(cfg.StatusPath)
    if e != nil {
        err = fmt.Errorf("status file: %s", e)
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

            log.Printf("recv line [%s]\n", string(buffer))

            // parse input
            parsed, err := stats.Parse(bytes.NewReader(buffer))
            if err != nil {
                log.Printf("parsing: %s\n", err)
                buffer = nil
                continue
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

                log.Printf("got transaction: %+v\n", tr)

                transactions = append(transactions, tr)
            case "ev":
                ev, err := stats.ProcessEvent(parsed)
                if err != nil {
                    log.Println(err)
                    continue
                }

                log.Printf("got event: %+v\n", ev)

                events = append(events, ev)

                log.Println("started timer")

                stats.StartTimer(ev, time.Now(), timer)
            default:
                log.Println(parsed[0], "is not a cmd")
                continue
            }

            // update stats
            s := stats.BuildStats(transactions, events)

            if err = stats.UpdateStats(s, status); err != nil {
                return fmt.Errorf("status update: %s", err)
            }

        case t := <-timer:
            log.Printf("timer triggered: %+v\n", t)

            i := findEvent(t.Id, events)
            if i < 0 {
                log.Printf("The event with id %d does not exist\n", t.Id)
                break
            }

            // selected event
            ev := &events[i]

            // build new transaction
            tr := stats.BuildTransaction(ev.Name, ev.Description, t.Date, ev.Amount)
            transactions = append(transactions, tr)

            ev.Times--

            // when times reaches 0 that means the event should not keep repeating
            // hence we only update the date and set up a new timer only if times
            // is greater than zero or negative (that means it will keep reapeating forever)
            if ev.Times != 0 {
                // update date
                newDate := ev.Date.AddDate(ev.Step[0], ev.Step[1], ev.Step[2])
                ev.Date = newDate

                // set new timer
                stats.StartTimer(*ev, time.Now(), timer)
            }

            // update stats
            s := stats.BuildStats(transactions, events)

            if err := stats.UpdateStats(s, status); err != nil {
                return fmt.Errorf("timer status update: %s", err)
            }

        case err := <-ctl.Err:
            return fmt.Errorf("control file erorr: %s", err)

        case <-sigs:
            break End
        }
    }

    return nil
}

func findEvent(id uint, events []stats.Event) int {
    for i, ev := range events {
        if ev.Id == id {
            return i
        }
    }

    return -1
}
