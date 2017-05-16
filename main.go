package main

import (
	"net/http"
	"os"
	"todo-search-service-go/elasticsearch"
	"gopkg.in/olivere/elastic.v5"
	"todo-search-service-go/handlers"
	"github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats-streaming"
	"todo-search-service-go/api"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

const natsClientName = "service_todo-search"
const natsClusterName = "test-cluster"

func main() {
	// Start NATS client
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		natsUrl = stan.DefaultNatsURL
	}
	natsSession, err := stan.Connect(natsClusterName, natsClientName, stan.NatsURL(natsUrl))
	if err != nil {
		panic(err)
	}
	logrus.Info("Initialized NATS streaming")
	defer natsSession.Close()

	// Start ElasticSearch client
	esUrl := os.Getenv("ES_URL")
	if esUrl == "" {
		esUrl = elastic.DefaultURL
	}
	es, err := elasticsearch.NewElasticSearch(esUrl)
	if err != nil {
		panic(err)
	}
	logrus.Println("Initialized ElasticSearch")

	// Start Handlers
	_, err = handlers.NewTodoSearchHandler(natsSession, es)
	if err != nil {
		panic(err)
	}
	logrus.Info("Initialized Handlers")

	// Start api
	a := api.NewApi()

	logrus.Fatal(http.ListenAndServe(":8080", a.Router))
}