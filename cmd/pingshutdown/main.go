package main

import (
	"fmt"
	"github.com/prometheus-community/pro-bing"
	"time"
	"flag"
	"github.com/devusb/pingshutdown/internal/countdown"
)

func shutdown() {
	fmt.Println("shutting down system!")
}

func main() {
	target := flag.String("t", "www.google.com", "target to begin shutdown when unavailable")

	flag.Parse()

	shutdownTimer := countdown.Countdown{Duration: 15*time.Second}

	go func() {
		for {
			pinger, _ := probing.NewPinger(*target)
			pinger.Count = 5
			pinger.Timeout = 5 * time.Second
			pinger.Run()
			if pinger.Statistics().PacketLoss == 100 {
				fmt.Println("all pings failed")
				shutdownTimer.StartAfterFunc(shutdown)	
			} else {
				fmt.Println("some pings succeeded")
				shutdownTimer.Stop()
			}
		}
	}()

	select {}
}
