package elasticsearch

import (
	"gopkg.in/olivere/elastic.v5"
	"context"
	"todo-service-go/events"
	"github.com/gocql/gocql"
	"github.com/Sirupsen/logrus"
)

type ElasticSearch struct {
	ctx context.Context
	client *elastic.Client
}

var (
	indexName = "todos"
	typeName = "todos"
	mapping = `{
		"properties" : {
			"title" : {
				"type" : "string",
				"fields" : {
					"raw" : {
						"type" : "string",
						"index" : "not_analyzed"
					}
				}
			},
			"completed": {
				"type":"boolean"
			}
		}
	}`
)

func NewElasticSearch(url string) (*ElasticSearch, error) {
	var err error

	es := ElasticSearch{}
	es.ctx = context.Background()
	es.client, err = elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		logrus.Printf("Unable to connect to ES url=%s", url)
		return nil, err
	}

	err = es.initIndex()
	if err != nil {
		return nil, err
	}

	return &es, nil
}

func (es *ElasticSearch) initIndex() error {
	exists, err := es.client.IndexExists(indexName).Do(es.ctx)
	if err != nil {
		return err
	}

	if !exists {
		// Create a new index.
		createIndex, err := es.client.CreateIndex(indexName).Do(es.ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	_, err = es.client.PutMapping().Index(indexName).Type(typeName).BodyString(mapping).Do(context.TODO())
	if err != nil {
		// Handle error
		return err
	}

	return nil
}

func (es *ElasticSearch) InsertTodo(event *events.TodoAdded) error {
	_, err := es.client.Index().
		Index(indexName).
		Type(typeName).
		Id(event.GetId()).
		BodyJson(event).
		Do(es.ctx)
	if err != nil {
		// Handle error
		logrus.Println("Unable to insert todo", err)
		return err
	}

	return nil
}

func (es *ElasticSearch) DeleteTodo(id gocql.UUID) error {
	_, err := es.client.Delete().
		Index(indexName).
		Type(typeName).
		Id(id.String()).
		Do(es.ctx)
	if err != nil {
		// Handle error
		logrus.Println("Unable to delete todo", err)
		return err
	}

	return nil
}