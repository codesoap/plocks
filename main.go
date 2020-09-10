package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type block struct {
	command         string
	interval        string // empty string if no regular update is wanted
	updateOnSIGUSR1 bool
	updateOnSIGUSR2 bool
	output          string // the result of the latest execution
}

func main() {
	ticks := make(chan int)
	generateTicks(ticks)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2)

	for i := range blocks {
		updateOutput(&blocks[i])
	}

	for {
		printStatus()
		select {
		case blockIndex := <-ticks:
			updateOutput(&blocks[blockIndex])
		case signal := <-signals:
			updateBlocksForSignal(signal)
		}
	}
}

func generateTicks(ticks chan int) {
	for i, block := range blocks {
		if block.interval != "" {
			go generateBlocksTicks(ticks, i, block.interval)
		}
	}
}

func generateBlocksTicks(ticks chan int, blockIndex int, interval string) {
	dur, err := time.ParseDuration(interval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse interval: %s\n", err.Error())
		os.Exit(1)
	}
	ticker := time.NewTicker(dur)
	for {
		<-ticker.C
		ticks <- blockIndex
	}
}

func updateBlocksForSignal(signal os.Signal) {
	for i := range blocks {
		sigusr1Matches := signal == syscall.SIGUSR1 && blocks[i].updateOnSIGUSR1
		sigusr2Matches := signal == syscall.SIGUSR2 && blocks[i].updateOnSIGUSR2
		if sigusr1Matches || sigusr2Matches {
			updateOutput(&blocks[i])
		}
	}
}

func updateOutput(block *block) {
	out, err := exec.Command("sh", "-c", block.command).Output()
	if err != nil {
		block.output = fmt.Sprintf("Error: %s", err.Error())
	} else if len(out) == 0 {
		block.output = "Error: no output for block"
	} else {
		block.output = strings.TrimSpace(string(out))
	}
}

func printStatus() {
	outputs := make([]string, 0, len(blocks))
	for _, block := range blocks {
		outputs = append(outputs, block.output)
	}
	fmt.Println(strings.Join(outputs, separator))
}
