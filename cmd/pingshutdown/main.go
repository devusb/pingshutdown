package main

import (
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/devusb/pingshutdown/internal/countdown"
	"github.com/devusb/pingshutdown/internal/pushover"
	"github.com/prometheus-community/pro-bing"
	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Target string `default:"www.google.com"`
	NotificationUser string
	NotificationToken string
	Delay time.Duration
}

func shutdown() {
	fmt.Println("shutting down system!")
}

func main() {
	var s Specification
	err := envconfig.Process("pingshutdown", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	shutdownTimer := countdown.Countdown{Duration: s.Delay}
	notification := pushover.Notification{
		User: s.NotificationUser,
		Token: s.NotificationToken,
	}
	timerLockout := false

	http.HandleFunc("/status", HandleStatus(&shutdownTimer, &timerLockout))
	http.HandleFunc("/lockout", HandleLockout(&timerLockout))

	go func() {
		for {
			pinger, _ := probing.NewPinger(s.Target)
			pinger.Count = 5
			pinger.Timeout = 5 * time.Second
			pinger.Run()
			if (pinger.Statistics().PacketLoss == 100 && !timerLockout && !shutdownTimer.Status()) {
				fmt.Println("all pings failed")
				_, err := notification.Send("all pings failed")
				if err != nil {
					log.Println(err)
				}
				shutdownTimer.StartAfterFunc(shutdown)

			} else if ((pinger.Statistics().PacketLoss < 100 || timerLockout) && shutdownTimer.Status()) {
				fmt.Println("some pings succeeded")
				shutdownTimer.Stop()
			} else {
				fmt.Printf("no state change, current timer state is %s\n", shutdownTimer.Status())
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":8081",nil))

	select {}
}
