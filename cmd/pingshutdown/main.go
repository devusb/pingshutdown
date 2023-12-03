package main

import (
	"fmt"
	"github.com/prometheus-community/pro-bing"
	"time"
)

func main() {
	pinger, _ := probing.NewPinger("www.google.com")

	go func() {
		pinger.Run()
	}()

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println(pinger.Statistics())
	}()

	select {}
}
