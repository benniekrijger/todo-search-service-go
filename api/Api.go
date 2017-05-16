package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/Sirupsen/logrus"
)

type Api struct {
	*mux.Router
}

func NewApi() *Api {
	router := mux.NewRouter()
	api := Api{router}

	apiRouter:= router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	apiRouter.HandleFunc("/", api.healthCheck)
	logrus.Info("Initialized API")

	return &api
}

func (a *Api) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))
}