package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/devusb/pingshutdown/internal/countdown"
)

func HandleStatus(c *countdown.Countdown, lock *bool, s *Specification) http.HandlerFunc {
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
		fmt.Fprintf(w, "Pinging %s\nShutdown timer is %s", s.Target, timerStatus)
	}
}

func HandleLockout(lock *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if *lock {
			*lock = false
		} else {
			*lock = true
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
