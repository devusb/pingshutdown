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
		w.Header().Set("Content-Type", "text/html")

		html := fmt.Sprintf("<html><body><p>Pinging %s</p>", s.Target)
		html += fmt.Sprintf("<p>Shutdown timer is %s</p>", timerStatus)
		html += "<p><a href=\"/lockout\">Toggle shutdown timer lockout</a></p></body></html>"

		fmt.Fprintf(w, html)
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
