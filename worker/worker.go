package worker

import (
	"cube-orchestrator/task"
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task // in-memory data store TODO: should evolve to a full fledged db
	TaskCount int
}

func (w *Worker) CollectStatus() {
	fmt.Println("Collect stats")
}

func (w *Worker) RunTask() {
	fmt.Println("Start and Stop the tasks")
}

func (w *Worker) StartTask() {
	fmt.Println("I will start the task")
}

func (w *Worker) StopTask() {
	fmt.Println("I will stop the task")
}
