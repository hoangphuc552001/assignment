package service

import (
	"github.com/caarlos0/env"
	"github.com/go-playground/assert/v2"
	"github.com/spf13/viper"
	"log"
	"sfvn-test/common/model"
	"testing"
	"time"
)

type (
	Config struct {
		Dir string `env:"CONFIG_DIR" envDefault:"../config/config.json"`
	}
)

var config Config
var coinGecko model.CoinGecko
var historiesTest IHistories

func init() {
	if err := env.Parse(&config); err != nil {
		log.Panicf("failed to parse config: %v", err)
		log.Fatal(err)
	}
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	cfg := Config{
		Dir: config.Dir,
	}
	cGecko := model.CoinGecko{
		Url:    viper.GetString(`coin_gecko.url`),
		ApiKey: viper.GetString(`coin_gecko.api_key`),
	}
	config = cfg
	coinGecko = cGecko
	historiesTest = &Histories{
		coinGecko: coinGecko,
	}
}
func TestGetHistories(t *testing.T) {
	tests := []struct {
		name         string
		params       model.HistoriesRequest
		expectedResp []model.HistoriesResponse
		expectedErr  error
		mockCallAPI  func(method, url string, params map[string]string) ([]byte, error)
	}{
		{
			name: "Return list of histories",
			params: model.HistoriesRequest{
				StartDate: time.Now().AddDate(0, 0, -1),
				EndDate:   time.Now(),
				Period:    model.Period1D,
				Symbol:    "bitcoin",
			},
			expectedResp: []model.HistoriesResponse{
				{
					High:   2,
					Low:    3,
					Open:   4,
					Close:  5,
					Time:   1,
					Change: 0,
				},
			},
			expectedErr: nil,
			mockCallAPI: func(method, url string, params map[string]string) ([]byte, error) {
				return []byte(`[[1,4,2,3,5]]`), nil
			},
		},
		{
			name:         "New Histories Service",
			params:       model.HistoriesRequest{},
			expectedResp: []model.HistoriesResponse{},
			expectedErr:  nil,
		},
		{
			name: "Not valid day",
			params: model.HistoriesRequest{
				StartDate: time.Now().AddDate(0, 0, 1),
				EndDate:   time.Now(),
				Period:    model.Period1D,
				Symbol:    "bitcoin",
			},
			expectedResp: []model.HistoriesResponse{},
		},
		{
			name: "Invalid period",
			params: model.HistoriesRequest{
				StartDate: time.Now().AddDate(0, 0, -1),
				EndDate:   time.Now(),
				Period:    "1A",
				Symbol:    "bitcoin",
			},
			expectedResp: []model.HistoriesResponse{},
		},
		{
			name: "Call response api error",
			params: model.HistoriesRequest{
				StartDate: time.Now().AddDate(0, 0, -1),
				EndDate:   time.Now(),
				Period:    model.Period30M,
				Symbol:    "abc",
			},
			expectedResp: []model.HistoriesResponse{},
		},
	}

	// 1
	tt := tests[0]
	t.Run(tt.name, func(t *testing.T) {
		resp, err := historiesTest.GetHistories(tt.params)
		assert.NotEqual(t, len(resp), 0)
		assert.Equal(t, tt.expectedErr, err)
	})

	// 2
	tt = tests[1]
	t.Run(tt.name, func(t *testing.T) {
		NewHistories(coinGecko)
	})

	// 3
	tt = tests[2]
	t.Run(tt.name, func(t *testing.T) {
		resp, err := historiesTest.GetHistories(tt.params)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, err.Error(), "out of range")
		assert.Equal(t, len(resp), 0)
	})

	// 4
	tt = tests[3]
	t.Run(tt.name, func(t *testing.T) {
		resp, err := historiesTest.GetHistories(tt.params)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, err.Error(), "invalid period")
		assert.Equal(t, len(resp), 0)
	})

	// 5
	tt = tests[4]
	t.Run(tt.name, func(t *testing.T) {
		_, err := historiesTest.GetHistories(tt.params)
		assert.NotEqual(t, err, nil)
	})
}
