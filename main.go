package main

import (
	"github.com/florhusq/digibank/config"
	"github.com/florhusq/digibank/event"
	"github.com/florhusq/digibank/rest"
)

func main() {
	config, err := config.Load("config.json")
	if err != nil {
		panic(err)
	}

	db, err := event.Open(config.Storage.Connection)
	if err != nil {
		panic(err)
	}

	rest.ServeAPI(config.Rest.Endpoint, config.Prometheus.Endpoint, db)
}
