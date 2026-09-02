package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"zekube/task"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.Event{}
	err := d.Decode(&te)
	if err != nil {
		msg := fmt.Sprintf("Error unmarshalling body: %v\n", err)
		log.Printf(msg)
		w.WriteHeader(400)

		e := ErrorResponse{
			HttpStatusCode: 400,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(e)
		return

	}
	a.Worker.AddTask(te.Task)
	log.Printf("Added task %v\n", te.Task.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
}

func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the query param on the url
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Printf("No task Id passed in the request \n")
		w.WriteHeader(http.StatusBadRequest)
	}
	// parse the query param to a valid uuid
	tID, _ := uuid.Parse(taskID)
	// Check if the task exists within the db
	_, ok := a.Worker.Db[tID]

	if !ok {
		log.Printf("No task with the ID %v found", tID)
		w.WriteHeader(http.StatusNotFound)
	}
	// retrieve the task from the db
	taskToStop := a.Worker.Db[tID]
	// Get the ptr to the task to be able to mutate it
	taskPtr := *taskToStop
	// change the state to completed
	taskPtr.State = task.Completed
	// Enqueue the task
	a.Worker.AddTask(taskPtr)

	log.Printf("Added task %v to stop container %v\n", taskToStop.ID, taskToStop.ContainerID)
	w.WriteHeader(http.StatusNoContent)
}

func (a *Api) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Worker.Stats)

}
