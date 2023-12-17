package main

import (
	"net/http"
	"fmt"
	"time"
	
	"github.com/devusb/pingshutdown/internal/countdown"
)


func HandleStatus(c *countdown.Countdown, lock *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timerStatus := "not in progress"
		secondsRemaining := c.EndTime().Sub(time.Now())
		if c.Status() && secondsRemaining > 0 {
			timerStatus = fmt.Sprintf("in progress, shutdown in %s", secondsRemaining)
		} else if c.Status() && secondsRemaining < 0 {
			timerStatus = "complete, shutdown in progress"
		} else if *lock == true {
			timerStatus = "locked out"
		}
		fmt.Fprintf(w, "Shutdown timer is %s", timerStatus)
	}
}

func HandleLockout(lock *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if *lock {
			*lock = false
		} else {
			*lock = true
		}
		http.Redirect(w, r, "/status", http.StatusSeeOther)
	}
}
