package models

type Config struct {
	Symbols    []string `yaml:"symbols"`
	MaxWorkers int      `yaml:"max_workers"`
}

type Ticker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}
