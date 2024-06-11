package main

import (
	"bufio"
	"context"
	"fmt"
	"go-questions/MyTasks/ticker_parser/models"
	"go-questions/MyTasks/ticker_parser/worker"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

func main() {
	configPath := "config.yaml"
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	var config models.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	numCPU := runtime.NumCPU()
	if config.MaxWorkers > numCPU {
		config.MaxWorkers = numCPU
	}

	workers := make([]*worker.Worker, config.MaxWorkers)
	for i := 0; i < config.MaxWorkers; i++ {
		workers[i] = &worker.Worker{}
	}

	for i, symbol := range config.Symbols {
		workers[i%config.MaxWorkers].Symbols = append(workers[i%config.MaxWorkers].Symbols, symbol)
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for _, wrk := range workers {
		wg.Add(1)
		go wrk.Run(ctx, &wg)
	}

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				totalRequests := 0
				for _, wrk := range workers {
					totalRequests += wrk.GetRequestsCount()
				}
				fmt.Printf("workers requests total: %d\n", totalRequests)
			}
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stopChan
		fmt.Println("Received stop signal, shutting down...")
		cancel()
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "STOP" {
			cancel()
			break
		}
	}

	wg.Wait()
}
