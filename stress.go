package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"
)

func main() {

	cpuCount := runtime.NumCPU()

	app := &cli.App{
		Name:  "stress",
		Usage: "A tool for stress testing commands",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cmd",
				Usage:    "Command to run for stress testing",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "runs",
				Usage:   "Number of times to run the command",
				Value:   100,
				Aliases: []string{"r"},
			},
			&cli.IntFlag{
				Name:    "parallel",
				Usage:   "Number of parallel executions",
				Value:   cpuCount,
				Aliases: []string{"p"},
			},
		},
		Action: func(c *cli.Context) error {

			log.Println("Running stress test with the following parameters:")
			log.Printf("Command: %s", c.String("cmd"))
			log.Printf("Runs: %d", c.Int("runs"))
			log.Printf("Parallel: %d", c.Int("parallel"))
			log.Printf("CPU count: %d", cpuCount)

			if cpuCount < c.Int("parallel") {
				log.Printf("Warning: The number of parallel executions is greater than the number of CPUs. This may cause performance issues.")
			}

			return runStressTest(c.String("cmd"), c.Int("runs"), c.Int("parallel"))
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runStressTest(cmdString string, runs int, parallelLimit int) error {
	statusChan := make(chan error, runs)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, parallelLimit)

	for i := 0; i < runs; i++ {
		wg.Add(1)
		go func(runNumber int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			// Split the command string into command and arguments.
			cmdParts := strings.Fields(cmdString)
			var stdout, stderr bytes.Buffer
			cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()

			if err != nil {
				errorMessage := fmt.Sprintf("Run %d failed: %v\nStderr: %s\n", runNumber, err, stderr.String())
				log.Println(errorMessage)
				statusChan <- fmt.Errorf(errorMessage)
			} else {
				log.Printf("Run %d succeeded\n", runNumber)
				statusChan <- nil
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(statusChan)
	}()

	for err := range statusChan {
		if err != nil {
			return fmt.Errorf("Test failed. Stopping all runs. Error: %v", err)
		}
	}

	return nil
}
