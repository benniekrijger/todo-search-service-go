package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/nats-io/go-nats"
	"os"
	"todo-search-service-go/elasticsearch"
	"gopkg.in/olivere/elastic.v5"
	"todo-search-service-go/handlers"
	"github.com/Sirupsen/logrus"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func main() {
	// Start NATS client
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}
	natsSession, err := nats.Connect(natsUrl)
	if err != nil {
		panic(err)
	}
	logrus.Println("Initialized NATS")
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

	// Start router
	router := mux.NewRouter()
	apiRouter:= router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	apiRouter.HandleFunc("/", healthCheck)
	logrus.Println("Initialized API")

	logrus.Fatal(http.ListenAndServe(":8080", router))
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))
}