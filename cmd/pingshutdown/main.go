package main

import (
	"fmt"
	"github.com/prometheus-community/pro-bing"
	"time"
	"flag"
)

func main() {
	target := flag.String("t", "www.google.com", "target to begin shutdown when unavailable")

	flag.Parse()

	c := make(chan bool)

	go func(c chan bool) {
		for {
			pinger, _ := probing.NewPinger(*target)
			pinger.Count = 10
			pinger.Timeout = 10 * time.Second
			pinger.Run()
			if pinger.Statistics().PacketLoss == 100 {
				fmt.Println("all pings failed") 
			} else {
				fmt.Println("some pings succeeded")
			}
		}
	}(c)

	go func(c chan bool) {
		shutdownTimer := time.NewTimer(3 * time.Minute)

	}(c)

	// go func() {
	// 	for {
	// 		time.Sleep(2 * time.Second)
	// 		fmt.Println(pinger.Statistics())

	// 	}

	// }()

	select {}
}
