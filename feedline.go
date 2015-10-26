package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func readStdin(flow chan interface{}, exit chan struct{}, lazy bool) {
	defer close(exit)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		flow <- struct{}{} // indicate stdin streams are flowing
		if lazy {
			var nlazy int
			nlazy = (<-flow).(int)
			for i := 0; i < nlazy; i++ {
				fmt.Println("")
			}
		}
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "cannot read from pipe:", err)
	}
}

func feedLine(times []time.Duration, flow chan interface{}, lazy bool) {
	i := 0
	nlazy := 0
	for {
		select {
		case <-flow:
			if lazy {
				flow <- nlazy
				nlazy = 0
			}
			i = 0
		case <-time.After(times[min(i, len(times)-1)]):
			if i == len(times) {
				continue
			}
			if lazy {
				nlazy++
			} else {
				fmt.Println("")
			}
			i++
		}
	}
}

func main() {
	lazy := flag.Bool("lazy", true, "when set, line will keeped, and added before print stdin")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("please specify feed times.")
		os.Exit(1)
	}

	feedTimes := make([]time.Duration, 0)

	for _, a := range args {
		// an arg indicate feeding time.
		// supported time units are "s", "m", "h".
		// ex) 1s, 10s, 5m, 24h
		t, err := time.ParseDuration(a)
		if err != nil {
			// log.Fatal(fmt.Sprintf("invalid argument: %d", a))
			log.Fatal(err)
		}
		feedTimes = append(feedTimes, t)
	}

	// readStdin will close 'exit' if nothing to read.
	exit := make(chan struct{})

	// 'flow' indicated whether text comes in from stdin.
	// feedLine will keep check
	flow := make(chan interface{})

	go readStdin(flow, exit, *lazy)
	go feedLine(feedTimes, flow, *lazy)

	<-exit
}
