package es

import (
	"github.com/zero7cola/gin-admin-core/config"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var EsClient *elasticsearch.Client

func InitEs() {
	cfg := elasticsearch.Config{
		Addresses: config.GetStringSlice("es.hosts"),
		Username:  config.GetString("es.username"),
		Password:  config.GetString("es.password"),
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	EsClient = client
}
