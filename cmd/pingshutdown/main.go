package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/devusb/pingshutdown/internal/countdown"
	"github.com/prometheus-community/pro-bing"
)

func shutdown() {
	fmt.Println("shutting down system!")
}

func main() {
	target := flag.String("t", "www.google.com", "target to begin shutdown when unavailable")

	flag.Parse()

	shutdownTimer := countdown.Countdown{Duration: 15*time.Second}
	timerLockout := false

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		timerStatus := "not in progress"
		secondsRemaining := shutdownTimer.EndTime().Sub(time.Now())
		if shutdownTimer.Status() && secondsRemaining > 0 {
			timerStatus = fmt.Sprintf("in progress, shutdown in %s", secondsRemaining)
		} else if shutdownTimer.Status() && secondsRemaining < 0 {
			timerStatus = "complete, shutdown in progress"
		} else if timerLockout == true {
			timerStatus = "locked out"
		}
		fmt.Fprintf(w, "Shutdown timer is %s", timerStatus)
	})

	http.HandleFunc("/lockout", func(w http.ResponseWriter, r *http.Request) {
		if timerLockout {
			timerLockout = false
		} else {
			timerLockout = true
		}
		http.Redirect(w, r, "/status", http.StatusSeeOther)
	})

	go func() {
		for {
			pinger, _ := probing.NewPinger(*target)
			pinger.Count = 5
			pinger.Timeout = 5 * time.Second
			pinger.Run()
			if (pinger.Statistics().PacketLoss == 100 && !timerLockout) {
				fmt.Println("all pings failed")
				shutdownTimer.StartAfterFunc(shutdown)	
			} else {
				fmt.Println("some pings succeeded")
				shutdownTimer.Stop()
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":8081",nil))

	select {}
}
