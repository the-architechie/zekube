package worker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ErrorResponse struct {
	HttpStatusCode int
	Message        string
}

type Api struct {
	Port    int
	Address string
	Worker  *Worker
	Router  *chi.Mux
}

func (a *Api) initRouter() {

	a.Router = chi.NewRouter()

	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler)
		r.Get("/", a.GetTasksHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler)
		})
	})
	a.Router.Route("/stats", func(r chi.Router) {
		r.Get("/", a.GetStatsHandler)
	})
}

func (a *Api) Start() {
	a.initRouter()
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", a.Address, a.Port), a.Router)
	if err != nil {
		log.Printf("Failed to start the server")
		return
	}
}
