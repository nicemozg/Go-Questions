package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"go-questions/MyTasks/ticker_parser/models"
	"io"
	"log"
	"net/http"
	"sync"
)

type Worker struct {
	Symbols       []string
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
			for _, symbol := range w.Symbols {
				resp, err := client.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol))

				w.mu.Lock()
				w.requestsCount++
				w.mu.Unlock()

				if err != nil {
					log.Println("Error fetching price:", err)
					continue
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println("Error reading response body:", err)
					continue
				}

				var ticker models.Ticker
				if err := json.Unmarshal(body, &ticker); err != nil {
					log.Println("Error unmarshalling JSON:", err)
					continue
				}

				changed := ""
				if previousPrice, ok := previousPrices[symbol]; ok && previousPrice != ticker.Price {
					changed = " changed"
				}
				previousPrices[symbol] = ticker.Price

				fmt.Printf("%s price:%s%s\n", ticker.Symbol, ticker.Price, changed)
			}
		}
	}
}

func (w *Worker) GetRequestsCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.requestsCount
}
