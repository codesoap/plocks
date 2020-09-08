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
	interval        int // in seconds; 0 if only initially
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
		if err := updateOutput(&blocks[i]); err != nil {
			fmt.Fprintf(os.Stderr, "Could not update block: %s", err.Error())
			os.Exit(1)
		}
	}

	for {
		printStatus()
		select {
		case blockIndex := <-ticks:
			if err := updateOutput(&blocks[blockIndex]); err != nil {
				fmt.Fprintf(os.Stderr, "Could not update output: %s", err.Error())
				os.Exit(1)
			}
		case signal := <-signals:
			if err := updateBlocksForSignal(signal); err != nil {
				fmt.Fprintf(os.Stderr, "Could not update output: %s", err.Error())
				os.Exit(1)
			}
		}
	}
}

func generateTicks(ticks chan int) {
	for i, block := range blocks {
		if block.interval > 0 {
			go generateBlocksTicks(ticks, i, block.interval)
		}
	}
}

func generateBlocksTicks(ticks chan int, blockIndex, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		<-ticker.C
		ticks <- blockIndex
	}
}

func updateBlocksForSignal(signal os.Signal) error {
	switch signal {
	case syscall.SIGUSR1:
		for i := range blocks {
			if blocks[i].updateOnSIGUSR1 {
				if err := updateOutput(&blocks[i]); err != nil {
					return err
				}
			}
		}
	case syscall.SIGUSR2:
		for i := range blocks {
			if blocks[i].updateOnSIGUSR2 {
				if err := updateOutput(&blocks[i]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func updateOutput(block *block) error {
	out, err := exec.Command("sh", "-c", block.command).Output()
	if err != nil {
		return err
	}
	block.output = strings.TrimSpace(string(out))
	return nil
}

func printStatus() {
	outputs := make([]string, 0, len(blocks))
	for _, block := range blocks {
		outputs = append(outputs, block.output)
	}
	fmt.Println(strings.Join(outputs, separator))
}
