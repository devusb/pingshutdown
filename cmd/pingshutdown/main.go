package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/devusb/pingshutdown/internal/countdown"
	"github.com/devusb/pingshutdown/internal/pushover"
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
		User:  s.NotificationUser,
		Token: s.NotificationToken,
	}
	timerLockout := false
	dryRun = s.DryRun

	http.HandleFunc("/", HandleStatus(&shutdownTimer, &timerLockout))
	http.HandleFunc("/lockout", HandleLockout(&timerLockout))

	go func() {
		for {
			pinger, _ := probing.NewPinger(s.Target)
			pinger.Count = 5
			pinger.Timeout = 5 * time.Second
			pinger.Run()
			if pinger.Statistics().PacketLoss == 100 && !timerLockout && !shutdownTimer.Status() {
				fmt.Println("all pings failed")
				if s.Notification {
					_, err := notification.Send("all pings failed")
					if err != nil {
						log.Println(err)
					}
				}
				shutdownTimer.StartAfterFunc(shutdown)

			} else if (pinger.Statistics().PacketLoss < 100 || timerLockout) && shutdownTimer.Status() {
				fmt.Println("some pings succeeded")
				shutdownTimer.Stop()
			} else {
				fmt.Printf("no state change, current timer state is %s\n", shutdownTimer.Status())
			}
		}
	}()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.StatusPort), nil))

	select {}
}
