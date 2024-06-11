package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Symbols    []string `yaml:"symbols"`
	MaxWorkers int      `yaml:"max_workers"`
}

type Ticker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type Worker struct {
	symbols       []string
	requestsCount int
	mu            sync.Mutex
}

func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	previousPrices := make(map[string]string)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, symbol := range w.symbols {
				resp, err := client.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol))
				if err != nil {
					log.Println("Error fetching price:", err)
					continue
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println("Error reading response body:", err)
					continue
				}

				var ticker Ticker
				if err := json.Unmarshal(body, &ticker); err != nil {
					log.Println("Error unmarshalling JSON:", err)
					continue
				}

				w.mu.Lock()
				w.requestsCount++
				w.mu.Unlock()

				changed := ""
				if previousPrice, ok := previousPrices[symbol]; ok && previousPrice != ticker.Price {
					changed = " changed"
				}
				previousPrices[symbol] = ticker.Price

				fmt.Printf("%s price:%s%s\n", ticker.Symbol, ticker.Price, changed)
			}
			time.Sleep(5 * time.Second) // Увеличиваем интервал времени для вывода
		}
	}
}

func (w *Worker) GetRequestsCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.requestsCount
}

func main() {
	configPath := "config.yaml"
	fmt.Printf("Reading config file from: %s\n", configPath)

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	numCPU := runtime.NumCPU()
	if config.MaxWorkers > numCPU {
		config.MaxWorkers = numCPU
	}

	workers := make([]*Worker, config.MaxWorkers)
	for i := 0; i < config.MaxWorkers; i++ {
		workers[i] = &Worker{}
	}

	for i, symbol := range config.Symbols {
		workers[i%config.MaxWorkers].symbols = append(workers[i%config.MaxWorkers].symbols, symbol)
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for _, worker := range workers {
		wg.Add(1)
		go worker.Run(ctx, &wg)
	}

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				totalRequests := 0
				for _, worker := range workers {
					totalRequests += worker.GetRequestsCount()
				}
				fmt.Printf("workers requests total: %d\n", totalRequests)
				fmt.Println("Enter 'STOP' to terminate the program.")
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
	fmt.Println("Program terminated gracefully.")
}
