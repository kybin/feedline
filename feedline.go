package main

import (
    "os"
    "bufio"
    "fmt"
    "time"
    "log"
)

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func readStdin(flow chan<- struct{}, exit chan struct{}) {
    defer close(exit)
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        flow <- struct{}{}
        fmt.Println(scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "cannot read from pipe:", err)
    }
}

func feedLine(times []time.Duration, flow <-chan struct{}) {
    i := 0
    for {
        select {
        case <-flow:
            i = 0
        case <-time.After(times[min(i, len(times)-1)]):
            if i == len(times) {
                continue
            }
            fmt.Println("")
            i++
        }
    }
}

func main() {
    args := os.Args[1:]
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
    flow := make(chan struct{})

    go readStdin(flow, exit)
    go feedLine(feedTimes, flow)

    <-exit
}
