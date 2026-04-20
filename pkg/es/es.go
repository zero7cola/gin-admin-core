package es

import (
	"github.com/elastic/go-elasticsearch/v8"
)

var EsClient *elasticsearch.Client

func InitEs() {
	//cfg := elasticsearch.Config{
	//	Addresses: setting.GetStringSlice("es.hosts"),
	//	Username:  setting.GetString("es.username"),
	//	Password:  setting.GetString("es.password"),
	//}
	//client, err := elasticsearch.NewClient(cfg)
	//if err != nil {
	//	log.Fatalf("Error creating the client: %s", err)
	//}
	//EsClient = client
}
