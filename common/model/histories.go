package model

import "time"

type Period string

const (
	Period30M Period = "30M"
	Period1H  Period = "1H"
	Period1D  Period = "1D"
)

type HistoriesRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Period    Period    `json:"period"`
	Symbol    string    `json:"symbol"`
}

type HistoriesResponse struct {
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Time   int     `json:"time"`
	Change float64 `json:"change"`
}
