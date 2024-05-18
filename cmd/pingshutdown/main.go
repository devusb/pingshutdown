package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/devusb/pingshutdown/pkg/countdown"
	"github.com/devusb/pingshutdown/pkg/pushover"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus-community/pro-bing"
)

type Specification struct {
	Target            string `default:"www.google.com"`
	NotificationUser  string
	NotificationToken string
	Notification      bool          `default:"false"`
	Delay             time.Duration `default:"5m"`
	StatusPort        string        `default:"8081"`
	DryRun            bool          `default:"false"`
}

func shutdown() {
	fmt.Println("shutting down system!")
	if !dryRun {
		cmd := exec.Command("poweroff")
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	}
}

var dryRun = false

func main() {
	var s Specification
	err := envconfig.Process("pingshutdown", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	shutdownTimer := countdown.Countdown{Duration: s.Delay}
	notification := pushover.Notification{
		User:         s.NotificationUser,
		Token:        s.NotificationToken,
		RetryWaitMax: s.Delay / 2,
		RetryMax:     20,
	}
	timerLockout := false
	dryRun = s.DryRun
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", HandleStatus(&shutdownTimer, &timerLockout, &s))
	http.HandleFunc("/lockout", HandleLockout(&timerLockout))

	go func() {
		for {
			pinger, err := probing.NewPinger(s.Target)
			if err != nil {
				panic(err)
			}
			pinger.Count = 5
			pinger.Timeout = 5 * time.Second
			err = pinger.Run()
			if err != nil {
				panic(err)
			}
			if (pinger.Statistics().PacketLoss == 100 || pinger.Statistics().PacketsSent == 0) && !timerLockout && !shutdownTimer.Status() {
				fmt.Println("all pings failed, timer started")
				shutdownTimer.StartAfterFunc(shutdown)
				if s.Notification {
					go func() {
						_, err := notification.Send(fmt.Sprintf("%s shutdown timer has started, shutdown at %s", hostname, shutdownTimer.EndTime()))
						if err != nil {
							log.Println(err)
						}
					}()

				}

			} else if (pinger.Statistics().PacketLoss < 100 || timerLockout) && shutdownTimer.Status() {
				fmt.Println("timer stopped")
				shutdownTimer.Stop()
			}
		}
	}()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.StatusPort), nil))

	select {}
}
