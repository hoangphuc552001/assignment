package main

import (
	"github.com/caarlos0/env"
	"github.com/spf13/viper"
	"log"
	"sfvn-test/api"
	v1 "sfvn-test/api/v1"
	"sfvn-test/common/model"
	"sfvn-test/service"
)

type (
	Config struct {
		Dir  string `env:"CONFIG_DIR" envDefault:"config/config.json"`
		Port string
	}
)

var config Config
var coinGecko model.CoinGecko

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
		Dir:  config.Dir,
		Port: viper.GetString(`main.port`),
	}
	cGecko := model.CoinGecko{
		Url:    viper.GetString(`coin_gecko.url`),
		ApiKey: viper.GetString(`coin_gecko.api_key`),
	}
	config = cfg
	coinGecko = cGecko
}

func main() {
	server := api.NewServer()
	v1.NewAPIHistories(server.Engine, service.NewHistories(coinGecko))
	server.Start(config.Port)
}
