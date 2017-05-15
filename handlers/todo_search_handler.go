package handlers

import (
	"todo-search-service-go/elasticsearch"
	"github.com/nats-io/go-nats"
	"todo-service-go/events"
	"github.com/golang/protobuf/proto"
	"github.com/gocql/gocql"
	"github.com/Sirupsen/logrus"
)

type TodoSearchHandler struct {
	natsSession *nats.Conn
	es *elasticsearch.ElasticSearch
}

func NewTodoSearchHandler(natsSession *nats.Conn, es *elasticsearch.ElasticSearch) (*TodoSearchHandler, error) {
	h := TodoSearchHandler{natsSession, es}

	_, err := natsSession.Subscribe("todos.new", func(msg *nats.Msg) {
		h.insertTodo(msg)
	})
	if err != nil {
		return nil, err
	}

	_, err = natsSession.Subscribe("todos.remove", func(msg *nats.Msg) {
		h.deleteTodo(msg)
	})
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *TodoSearchHandler) insertTodo(m *nats.Msg) error {
	event := events.TodoAdded{}
	err := proto.Unmarshal(m.Data, &event)
	if err != nil {
		logrus.Println("Unable to unmarshal todo added event", err)
		return err
	}

	err = h.es.InsertTodo(&event)
	if err != nil {
		logrus.Println("Unable insert todo", err)
		return err
	}

	logrus.Printf("Inserted todo, id=%s", event.GetId())

	return nil
}

func (h *TodoSearchHandler) deleteTodo(m *nats.Msg) error {
	event := events.TodoRemoved{}
	err := proto.Unmarshal(m.Data, &event)
	if err != nil {
		logrus.Println("Unable to unmarshal todo removed event", err)
		return err
	}

	id, err := gocql.ParseUUID(event.GetId())
	if err != nil {
		logrus.Println("Unable to parse todo id", err)
		return err
	}

	err = h.es.DeleteTodo(id)
	if err != nil {
		logrus.Println("Unable to delete todo", err)
		return err
	}

	logrus.Printf("Removed todo, id=%s", event.GetId())

	return nil
}