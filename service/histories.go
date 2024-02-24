package service

import (
	"encoding/json"
	"fmt"
	"sfvn-test/common/model"
	"sfvn-test/common/util"
	"time"
)

type (
	IHistories interface {
		GetHistories(params model.HistoriesRequest) (resp []model.HistoriesResponse, err error)
	}
	Histories struct {
		coinGecko model.CoinGecko
	}
)

func NewHistories(gecko model.CoinGecko) IHistories {
	return &Histories{
		coinGecko: gecko,
	}
}

func (h *Histories) GetHistories(params model.HistoriesRequest) (resp []model.HistoriesResponse, err error) {
	day, err := h.handleGetDays(params.StartDate)
	if err != nil {
		return resp, err
	}
	timeCache := h.getTimePeriod(*day)
	//check valid period
	isValidPeriod := h.checkValidPeriod(params.Period, timeCache)
	if !isValidPeriod {
		return resp, fmt.Errorf("invalid period")
	}
	// caching
	keyCache := fmt.Sprintf("histories-%s-%d", params.Symbol, *day)
	cache, ok := util.GetCache(keyCache)
	if ok {
		resp, ok = cache.([]model.HistoriesResponse)
		if !ok {
			return resp, fmt.Errorf("error when get cache")
		}
		return resp, nil
	} else {
		response, err := util.CallAPI("GET", h.coinGecko.Url+"/coins/"+params.Symbol+"/ohlc", map[string]string{
			"vs_currency": "usd",
			"days":        fmt.Sprintf("%d", *day),
			"api_key":     h.coinGecko.ApiKey,
		})
		if err != nil {
			return resp, err
		}
		var historiesResp [][]float64
		err = json.Unmarshal(response, &historiesResp)
		if err != nil {
			return resp, err
		}
		var result []model.HistoriesResponse
		for _, history := range historiesResp {
			var change float64
			if len(result) > 0 {
				change = h.handleGetChange(history[4], result[len(result)-1].Close)
			}
			result = append(result, model.HistoriesResponse{
				High:   history[2],
				Low:    history[3],
				Open:   history[1],
				Close:  history[4],
				Time:   int(history[0]),
				Change: change,
			})
		}
		util.SetCache(keyCache, result, timeCache)
		return result, nil
	}
}

func (h *Histories) handleGetDays(startDate time.Time) (*int, error) {
	today := time.Now()
	day := int(today.Sub(startDate).Hours() / 24)
	isValid := checkValidDay(day)
	if !isValid {
		return nil, fmt.Errorf("out of range")
	}
	return &day, nil
}

func (h *Histories) getTimePeriod(day int) time.Duration {
	switch day {
	case 1:
		return 30 * time.Minute
	case 7:
	case 14:
	case 30:
		return 4 * time.Hour
	case 60:
	case 90:
	case 180:
	case 365:
		return 4 * 24 * time.Hour
	}
	return 5 * time.Minute
}

func (h *Histories) checkValidPeriod(period model.Period, cachePeriod time.Duration) bool {
	switch period {
	case model.Period30M:
		return cachePeriod <= 30*time.Minute
	case model.Period1H:
		return cachePeriod <= time.Hour
	case model.Period1D:
		return cachePeriod <= 24*time.Hour
	}
	return false
}

func (h *Histories) handleGetChange(new, origin float64) float64 {
	return (new - origin) / origin * 100
}

func checkValidDay(day int) bool {
	switch day {
	case 1, 7, 14, 30, 60, 90, 180, 365:
		return true
	}
	return false
}
